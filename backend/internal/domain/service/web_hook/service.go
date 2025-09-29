package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Event struct {
	CommandID        string    `json:"command_id"`
	DeviceName       string    `json:"device_name"`
	Status           string    `json:"status"`
	PreviousStatus   string    `json:"previous_status"`
	Data             string    `json:"data,omitempty"`
	ErrorDescription string    `json:"error_desc,omitempty"`
	UpdatedAt        time.Time `json:"updated_at"`
	CallbackURL      *string   `json:"-"`
	CallbackSecret   *string   `json:"-"`
}

type EventBody struct {
	CommandID        string    `json:"command_id"`
	DeviceName       string    `json:"device_name"`
	Status           string    `json:"status"`
	PreviousStatus   string    `json:"previous_status"`
	Data             string    `json:"data,omitempty"`
	ErrorDescription string    `json:"error_desc,omitempty"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type EventWrapper struct {
	event   Event
	attempt int
}

type WebHookService struct {
	eventQueue chan EventWrapper
	httpClient *http.Client
	maxRetry   int
	ctx        context.Context
}

func NewService(ctx context.Context, bufferSize int) *WebHookService {
	service := WebHookService{
		eventQueue: make(chan EventWrapper, bufferSize),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		maxRetry: 3,
		ctx:      ctx,
	}

	go service.loop()

	return &service
}

func (s *WebHookService) loop() {
	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("closing webook loop")
			return
		case ew := <-s.eventQueue:
			s.send(ew)
			// handle event
		}
	}
}

func (s *WebHookService) Publish(event Event) {
	select {
	case s.eventQueue <- EventWrapper{event: event, attempt: 1}:
		// ok
	default:
		// drop
	}
}

func (s *WebHookService) requeue(ew EventWrapper) {
	select {
	case s.eventQueue <- ew:
		// ok
	default:
		// drop
	}
}

func (s *WebHookService) send(ew EventWrapper) {
	eventBody := EventBody{
		CommandID:        ew.event.CommandID,
		DeviceName:       ew.event.DeviceName,
		Status:           ew.event.Status,
		PreviousStatus:   ew.event.PreviousStatus,
		Data:             ew.event.Data,
		ErrorDescription: ew.event.ErrorDescription,
		UpdatedAt:        ew.event.UpdatedAt,
	}

	body, err := json.Marshal(eventBody)
	if err != nil {
		fmt.Println("could not marhal event: ", err.Error())
		return
	}

	if ew.event.CallbackURL == nil || *ew.event.CallbackURL == "" {
		fmt.Println("callback URL is not valid")
		return
	}

	req, err := http.NewRequestWithContext(
		s.ctx,
		"POST",
		*ew.event.CallbackURL,
		bytes.NewReader(body),
	)

	if err != nil {
		fmt.Println("could not create request with context: ", err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "OcelotMDM-EphemeralWebhook/0.1")
	req.Header.Set("X-Ocelot-Command-ID", ew.event.CommandID)
	req.Header.Set("X-Ocelot-Command-Status", ew.event.Status)
	req.Header.Set("X-Ocelot-Delivery-Attempt", strconv.Itoa(ew.attempt))
	req.Header.Set("X-Ocelot-Sent-At", time.Now().UTC().Format(time.RFC3339))
	if ew.event.CallbackSecret != nil && *ew.event.CallbackSecret != "" {
		req.Header.Set(
			"X-Ocelot-Signature",
			sign(*ew.event.CallbackSecret, body),
		)
	}

	res, err := s.httpClient.Do(req)
	if res != nil && res.Body != nil {
		_ = res.Body.Close()
	}

	if err != nil || res.StatusCode < 200 || res.StatusCode >= 300 {
		fmt.Println("could not complete request: ", err.Error())
		if ew.attempt < s.maxRetry {
			ew.attempt++
			s.requeue(ew)
		}
		return
	}
}

func sign(secret string, body []byte) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)

	return "sha256=" + hex.EncodeToString(h.Sum(nil))
}
