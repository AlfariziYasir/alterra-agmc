package auth

import (
	"api-mvc/database/postgres"
	"api-mvc/database/seeder"
	"api-mvc/internal/dto"
	"api-mvc/internal/factory"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceLogin(t *testing.T) {
	type args struct {
		ctx     context.Context
		payload *dto.LoginRequest
	}

	cases := []struct {
		name    string
		args    args
		want    *dto.UserJwtResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.TODO(),
				payload: &dto.LoginRequest{
					Email:    "test@test.com",
					Password: "12345678",
				},
			},
			want:    &dto.UserJwtResponse{},
			wantErr: false,
		},
		{
			name: "record not foudn",
			args: args{
				ctx: context.TODO(),
				payload: &dto.LoginRequest{
					Email:    "test2@test.com",
					Password: "12345678",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "password unmatch",
			args: args{
				ctx: context.TODO(),
				payload: &dto.LoginRequest{
					Email:    "test@test.com",
					Password: "1234567",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	_, err := postgres.NewClient()
	if err != nil {
		panic(err)
	}
	seeder.NewSeeder().DeleteAll()
	seeder.NewSeeder().SeedAll()
	asserts := assert.New(t)

	var (
		authService = NewService(factory.NewFactory())
	)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := authService.Login(c.args.ctx, c.args.payload)

			asserts.Equal((err != nil), c.wantErr)
			asserts.Equal(c.want, res)
		})
	}
}
