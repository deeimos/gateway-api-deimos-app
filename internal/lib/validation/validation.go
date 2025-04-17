package validation

import (
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorMapper struct {
	userMessages map[error]string
}

func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{
		userMessages: map[error]string{
			ErrNoMetadata:   "Metadata отсутствует",
			ErrInvalidToken: "Неверный или истекший токен",
			ErrGateway:      "Ошибка при обращении к микросервису",
		},
	}
}

var (
	ErrNoMetadata   = errors.New("missing metadata")
	ErrInvalidToken = errors.New("invalid token")
	ErrGateway      = errors.New("internal gateway error")
)

func (em *ErrorMapper) HandleGRPC(err error) error {
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.Unavailable:
			return status.Error(codes.Unavailable, "Сервис временно недоступен")
		default:
			return s.Err()
		}
	}

	switch {
	case errors.Is(err, ErrNoMetadata):
		return status.Error(codes.Unauthenticated, em.userMessageFor(err))
	case errors.Is(err, ErrInvalidToken):
		return status.Error(codes.Unauthenticated, em.userMessageFor(err))
	case errors.Is(err, ErrGateway):
		return status.Error(codes.Internal, em.userMessageFor(err))
	default:
		return status.Errorf(codes.Internal, "Внутренняя ошибка: %v", err)
	}
}

func (em *ErrorMapper) userMessageFor(err error) string {
	if msg, ok := em.userMessages[err]; ok {
		return msg
	}
	return err.Error()
}

func WrapGatewayError(op string, err error) error {
	return fmt.Errorf("%s: %w", op, ErrGateway)
}

func HandleGRPCServiceError(log *slog.Logger, op string, err error) error {
	if _, ok := status.FromError(err); ok {
		return err
	}

	log.Error("service unreachable", slog.String("op", op), slog.String("err", err.Error()))
	return fmt.Errorf("%s: %w", op, ErrGateway)
}
