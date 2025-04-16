package validation

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorMapper struct {
	userMessages map[error]string
}

func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{
		userMessages: map[error]string{
			ErrInvalidArgument:   "Некорректное значение",
			ErrNotFound:          "Не найдено",
			ErrUserNotFound:      "Пользователь не найден",
			ErrInvalidToken:      "Неверный токен",
			ErrTokenNotFound:     "Токен не найден или истёк",
			ErrTokenSaveFailed:   "Ошибка сохранения токена",
			ErrTokenRemoveFailed: "Ошибка удаления токена",
		},
	}
}

var (
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrNotFound          = errors.New("not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidToken      = errors.New("invalid refresh token")
	ErrAlreadyExists     = errors.New("user already exists")
	ErrTokenNotFound     = errors.New("refresh token not found or expired")
	ErrTokenSaveFailed   = errors.New("failed to save refresh token")
	ErrTokenRemoveFailed = errors.New("failed to remove refresh token")
)

type FieldError struct {
	Field string
	Err   error
}

func (e *FieldError) Error() string {
	msg := e.Err.Error()
	if idx := strings.Index(msg, ":"); idx != -1 {
		msg = strings.TrimSpace(msg[idx+1:])
	}
	return fmt.Sprintf("%s: %s", e.Field, msg)
}
func (e *FieldError) Unwrap() error {
	return e.Err
}

type FieldErrors []*FieldError

func (fe FieldErrors) Error() string {
	var b strings.Builder
	for i, f := range fe {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(f.Error())
	}
	return b.String()
}

func AsFieldErrors(err error) (FieldErrors, bool) {
	var fe FieldErrors
	if errors.As(err, &fe) {
		return fe, true
	}
	return nil, false
}

func (em *ErrorMapper) HandleGRPC(err error) error {
	if fe, ok := AsFieldErrors(err); ok {
		return status.Errorf(codes.InvalidArgument, fe.Error())
	}

	switch {
	case errors.Is(err, ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, em.messageFor(err))
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrUserNotFound):
		return status.Error(codes.NotFound, em.messageFor(err))
	case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrTokenNotFound):
		return status.Error(codes.Unauthenticated, em.messageFor(err))
	case errors.Is(err, ErrTokenSaveFailed), errors.Is(err, ErrTokenRemoveFailed):
		return status.Error(codes.Internal, em.messageFor(err))
	default:
		return status.Errorf(codes.Internal, "Внутренняя ошибка: %v", err)
	}

}

func (em *ErrorMapper) messageFor(err error) string {
	if msg, ok := em.userMessages[err]; ok {
		return msg
	}
	return err.Error()
}
