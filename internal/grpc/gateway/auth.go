package gateway

import (
	"context"
	"fmt"
	"gateway-api/internal/lib/validation"

	auth_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/auth-api"
	gateway_apiv1 "github.com/deeimos/proto-deimos-app/gen/go/public-api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string) (*auth_apiv1.LoginResponse, error)
	Register(ctx context.Context, name string, email string, password string) (*auth_apiv1.RegisterResponse, error)
	Refresh(ctx context.Context, refersh string) (*auth_apiv1.RefreshResponse, error)
	GetUser(ctx context.Context, token string) (*auth_apiv1.GetUserResponse, error)
}

func (s *serverApi) Login(ctx context.Context, req *gateway_apiv1.LoginRequest) (*gateway_apiv1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, s.errMapper.HandleGRPC(err)
	}
	panic("implemet me")
}

func (s *serverApi) Register(ctx context.Context, req *gateway_apiv1.RegisterRequest) (*gateway_apiv1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, s.errMapper.HandleGRPC(err)
	}
	panic("implemet me")
}

func (s *serverApi) Refresh(ctx context.Context, req *gateway_apiv1.RefreshRequest) (*gateway_apiv1.RefreshResponse, error) {
	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.Unauthenticated, "Отсутствует токен")
	}

	panic("implemet me")
}

func (s *serverApi) GetUser(ctx context.Context, req *gateway_apiv1.GetUserRequest) (*gateway_apiv1.GetUserResponse, error) {
	panic("implemet me")
}

func validateLogin(req *gateway_apiv1.LoginRequest) error {
	var errs validation.FieldErrors

	if req.GetEmail() == "" {
		errs = append(errs, &validation.FieldError{
			Field: "email",
			Err:   fmt.Errorf("%w: Поле не должно быть пустым", validation.ErrInvalidArgument),
		})
	}
	if req.GetPassword() == "" {
		errs = append(errs, &validation.FieldError{
			Field: "password",
			Err:   fmt.Errorf("%w: Поле не должно быть пустым", validation.ErrInvalidArgument),
		})
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateRegister(req *gateway_apiv1.RegisterRequest) error {
	var errs validation.FieldErrors

	if req.GetName() == "" {
		errs = append(errs, &validation.FieldError{
			Field: "name",
			Err:   fmt.Errorf("%w: Поле не должно быть пустым", validation.ErrInvalidArgument),
		})
	}
	if req.GetEmail() == "" {
		errs = append(errs, &validation.FieldError{
			Field: "email",
			Err:   fmt.Errorf("%w: Поле не должно быть пустым", validation.ErrInvalidArgument),
		})
	}
	if req.GetPassword() == "" {
		errs = append(errs, &validation.FieldError{
			Field: "password",
			Err:   fmt.Errorf("%w: Поле не должно быть пустым", validation.ErrInvalidArgument),
		})
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
