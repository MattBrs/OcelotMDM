package vpn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Service struct {
	baseUrl    string
	httpClient *http.Client
}

func NewService(baseUrl string) *Service {
	return &Service{
		baseUrl,
		&http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Service) RequestCertCreation(deviceName string) ([]byte, error) {
	url := fmt.Sprintf("%s/vpn/client/new", s.baseUrl)

	reqBody := CreateClientRequest{
		DeviceName: deviceName,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, ErrReqParsing
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, ErrReqCreation
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, ErrReq
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, ErrReadResponse
	}

	if res.StatusCode != 200 {
		return nil, ErrCreatingCertificate
	}

	var parsedRes CreateClientResponse
	if err := json.Unmarshal(body, &parsedRes); err != nil {
		return nil, ErrParsingResponse
	}

	return parsedRes.OvpnFile, nil
}
