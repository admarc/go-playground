package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/admarc/users/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUsers_Create(t *testing.T) {
	type fields struct {
		user UsersService
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCode int
		wantBody []byte
	}{
		{
			name: "failure when payload can't be decoded",
			fields: fields{
				user: &UsersServiceMock{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", strings.NewReader(`bad payload`)),
			},
			wantCode: http.StatusInternalServerError,
			wantBody: []byte("\n"),
		},
		{
			name: "failure when service fails with invalid name",
			fields: fields{
				user: &UsersServiceMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						assert.Equal(t, "mike", name)
						return models.User{}, models.UserCreateParamInvalidNameErr
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", strings.NewReader(`{"name": "mike"}`)),
			},
			wantCode: http.StatusBadRequest,
			wantBody: []byte("\n"),
		},
		{
			name: "failure when service fails with unknown error",
			fields: fields{
				user: &UsersServiceMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						assert.Equal(t, "mike", name)
						return models.User{}, errors.New("unknown error")
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", strings.NewReader(`{"name": "mike"}`)),
			},
			wantCode: http.StatusInternalServerError,
			wantBody: []byte("\n"),
		},
		{
			name: "success",
			fields: fields{
				user: &UsersServiceMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						assert.Equal(t, "mike", name)
						return models.User{Name: name, ID: "1"}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", strings.NewReader(`{"name": "mike"}`)),
			},
			wantCode: http.StatusOK,
			wantBody: []byte(`{"id":"1","name":"mike"}` + "\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{
				user: tt.fields.user,
			}
			u.Create(tt.args.w, tt.args.r)
			assert.Equal(t, tt.wantCode, tt.args.w.Code)
			assert.Equal(t, tt.args.w.Body.Bytes(), tt.wantBody)
		})
	}
}
