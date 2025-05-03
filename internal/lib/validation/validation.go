package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"google.golang.org/grpc/status"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrGateway      = errors.New("internal gateway error")
)

func WrapGatewayError(op string, err error) error {
	return fmt.Errorf("%s: %w", op, ErrGateway)
}

func HandleGRPCServiceError(log *slog.Logger, op string, err error) error {
	if s, ok := status.FromError(err); ok {
		return s.Err()
	}

	log.Error("service unreachable", slog.String("op", op), slog.String("err", err.Error()))
	return WrapGatewayError(op, err)
}

func BuildErrorResponse(err error) map[string]interface{} {
	s, ok := status.FromError(err)
	if !ok {
		return map[string]interface{}{"message": "Внутренняя ошибка сервера"}
	}

	pairs := strings.Split(s.Message(), ";")
	fieldErrors := make(map[string]interface{})

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			message := strings.TrimSpace(parts[1])
			fieldErrors[field] = message
		}
	}

	if len(fieldErrors) > 0 {
		return fieldErrors
	}

	return map[string]interface{}{"message": s.Message()}
}

func WriteError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	payload := BuildErrorResponse(err)
	_ = json.NewEncoder(w).Encode(payload)
}
