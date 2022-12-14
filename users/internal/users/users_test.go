package users

import (
	"context"
	"errors"
	"testing"

	"github.com/admarc/users/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestService_Create(t *testing.T) {
	type fields struct {
		repo Repository
	}
	type args struct {
		name string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       models.User
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "failure - empty name",
			fields: fields{
				repo: &RepositoryMock{},
			},
			args: args{
				name: "",
			},
			want:       models.User{},
			wantErr:    true,
			wantErrMsg: "invalid name argument",
		},
		{
			name: "failure - repository error",
			fields: fields{
				repo: &RepositoryMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						assert.Equal(t, "Tod", name)
						return models.User{}, errors.New("Failed to execute insert")
					},
				},
			},
			args: args{
				name: "Tod",
			},
			want:       models.User{},
			wantErr:    true,
			wantErrMsg: "Failed to execute insert",
		},
		{
			name: "success",
			fields: fields{
				repo: &RepositoryMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						assert.Equal(t, "Tod", name)
						return models.User{ID: "964e531c-7aba-49d1-87c6-7d37b0291d77", Name: "Tod"}, nil
					},
				},
			},
			args: args{
				name: "Tod",
			},
			want:       models.User{ID: "964e531c-7aba-49d1-87c6-7d37b0291d77", Name: "Tod"},
			wantErr:    false,
			wantErrMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				repo: tt.fields.repo,
			}
			ctx := context.TODO()
			got, err := s.Create(ctx, tt.args.name)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.Containsf(t, err.Error(), tt.wantErrMsg, "expected error containing %q, got %s", tt.wantErrMsg, err)
			}

		})
	}
}

func TestService_Get(t *testing.T) {
	type fields struct {
		repo Repository
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       models.User
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "failure - repository error",
			fields: fields{
				repo: &RepositoryMock{
					GetFunc: func(ctx context.Context, id string) (models.User, error) {
						assert.Equal(t, "383673b8-bd9a-41b4-adba-79bc1abc889e", id)
						return models.User{}, errors.New("Failed to fetch user")
					},
				},
			},
			args: args{
				id: "383673b8-bd9a-41b4-adba-79bc1abc889e",
			},
			want:       models.User{},
			wantErr:    true,
			wantErrMsg: "Failed to fetch user",
		},
		{
			name: "success",
			fields: fields{
				repo: &RepositoryMock{
					GetFunc: func(ctx context.Context, id string) (models.User, error) {
						assert.Equal(t, "383673b8-bd9a-41b4-adba-79bc1abc889e", id)
						return models.User{ID: "383673b8-bd9a-41b4-adba-79bc1abc889e", Name: "tod"}, nil
					},
				},
			},
			args: args{
				id: "383673b8-bd9a-41b4-adba-79bc1abc889e",
			},
			want:       models.User{ID: "383673b8-bd9a-41b4-adba-79bc1abc889e", Name: "tod"},
			wantErr:    false,
			wantErrMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				repo: tt.fields.repo,
			}
			got, err := s.Get(tt.args.ctx, tt.args.id)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.Containsf(t, err.Error(), tt.wantErrMsg, "expected error containing %q, got %s", tt.wantErrMsg, err)
			}
		})
	}
}
