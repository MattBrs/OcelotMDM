package api

import (
	"net/http"

	"github.com/MattBrs/OcelotMDM/internal/api/dto/token_dto"
	"github.com/MattBrs/OcelotMDM/internal/domain/token"
	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	service *token.Service
}

func NewTokenHandler(service *token.Service) *TokenHandler {
	return &TokenHandler{service}
}

func (h *TokenHandler) RequestToken(ctx *gin.Context) {
	token, err := h.service.GenerateNewToken(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Could not process request"},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		token_dto.NewTokenResponse{
			Token:     token.Token,
			ExpiresAt: token.ExpiresAt,
		},
	)
}

func (h *TokenHandler) VerifyToken(ctx *gin.Context) {
	otp := ctx.Query("token")
	if otp == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not validate token"})
		return
	}

	status, err := h.service.Verify(ctx, otp)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "An error occurred while parsing the token"},
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"valid": status})
}
