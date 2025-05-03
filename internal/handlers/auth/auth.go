package authHandler

import (
	"context"
	"encoding/json"
	"errors"
	"gateway-api/internal/domain/models"
	"gateway-api/internal/lib/utils/formatTimestamp"
	"gateway-api/internal/lib/validation"
	"gateway-api/internal/services/auth"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	service *auth.Auth
	timeout time.Duration
}

func NewAuthHandler(service *auth.Auth, timeout time.Duration) *AuthHandler {
	return &AuthHandler{service: service, timeout: timeout}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		validation.WriteError(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	resp, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		validation.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	expiresAt, _ := extractExpiresAt(resp.GetToken())

	dto := models.User{
		ID:           resp.GetId(),
		Email:        resp.GetEmail(),
		Name:         resp.GetName(),
		CreatedAt:    formatTimestamp.FormatTimestamp(resp.GetCreatedAt()),
		AccessToken:  resp.GetToken(),
		ExpiresAt:    expiresAt,
		RefreshToken: resp.GetRefreshToken(),
	}

	_ = json.NewEncoder(w).Encode(dto)
}
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		validation.WriteError(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	resp, err := h.service.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		validation.WriteError(w, err, http.StatusBadRequest)
		return
	}

	expiresAt, _ := extractExpiresAt(resp.GetToken())

	dto := models.User{
		ID:           resp.GetId(),
		Email:        resp.GetEmail(),
		Name:         resp.GetName(),
		CreatedAt:    formatTimestamp.FormatTimestamp(resp.GetCreatedAt()),
		AccessToken:  resp.GetToken(),
		ExpiresAt:    expiresAt,
		RefreshToken: resp.GetRefreshToken(),
	}

	_ = json.NewEncoder(w).Encode(dto)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		validation.WriteError(w, err, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	resp, err := h.service.Refresh(ctx, req.RefreshToken)
	if err != nil {
		validation.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	expiresAt, _ := extractExpiresAt(resp.GetToken())

	dto := models.Tokens{
		AccessToken:  resp.GetToken(),
		ExpiresAt:    expiresAt,
		RefreshToken: resp.GetRefreshToken(),
	}

	_ = json.NewEncoder(w).Encode(dto)
}

func (h *AuthHandler) User(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		validation.WriteError(w, validation.ErrInvalidToken, http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	resp, err := h.service.GetUser(ctx, token)
	if err != nil {
		validation.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	dto := models.UserInfo{
		ID:        resp.GetId(),
		Email:     resp.GetEmail(),
		Name:      resp.GetName(),
		CreatedAt: formatTimestamp.FormatTimestamp(resp.GetCreatedAt()),
	}

	_ = json.NewEncoder(w).Encode(dto)
}

func extractExpiresAt(token string) (int64, error) {
	parser := jwt.NewParser()
	claims := jwt.MapClaims{}

	_, _, err := parser.ParseUnverified(token, claims)
	if err != nil {
		return 0, err
	}

	if expRaw, ok := claims["exp"]; ok {
		switch exp := expRaw.(type) {
		case float64:
			return int64(exp), nil
		case json.Number:
			return exp.Int64()
		}
	}

	return 0, errors.New("exp not found in token")
}
