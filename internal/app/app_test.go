package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	originalStore "github.com/rawen554/go-loyal/internal/adapters/store"
	"github.com/rawen554/go-loyal/internal/adapters/store/mocks"
	"github.com/rawen554/go-loyal/internal/config"
	"github.com/rawen554/go-loyal/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config, err := config.ParseFlags()
	if err != nil {
		t.Error(err)
	}

	store := mocks.NewMockStore(ctrl)

	gomock.InOrder(
		store.EXPECT().CreateUser(gomock.Any()).Return(int64(1), nil),
		store.EXPECT().CreateUser(gomock.Any()).Return(int64(0), originalStore.ErrLoginNotFound),
	)

	app := NewApp(config, store, zap.L().Sugar())
	r, err := app.SetupRouter()
	if err != nil {
		t.Error(err)
	}

	srv := httptest.NewServer(r)
	defer srv.Close()

	tests := []struct {
		userCreds models.UserCredentialsSchema
		name      string
		url       string
		method    string
		status    int
	}{
		{
			name:      "Register user",
			userCreds: models.UserCredentialsSchema{Login: "a", Password: "b"},
			url:       "/api/user/register",
			status:    200,
			method:    http.MethodPost,
		},
		{
			name:      "Register user with conflict",
			userCreds: models.UserCredentialsSchema{Login: "a", Password: "b"},
			url:       "/api/user/register",
			status:    409,
			method:    http.MethodPost,
		},
	}

	for _, tt := range tests {
		tt := tt

		b, err := json.Marshal(tt.userCreds)
		if err != nil {
			t.Error(err)
		}

		url, err := url.JoinPath(srv.URL, tt.url)
		if err != nil {
			t.Error(err)
		}

		req, err := http.NewRequest(tt.method, url, bytes.NewBuffer(b))
		if err != nil {
			t.Error(err)
		}
		if err := req.Body.Close(); err != nil {
			t.Error(err)
		}

		res, err := srv.Client().Do(req)
		if err != nil {
			t.Error(err)
		}
		if err := res.Body.Close(); err != nil {
			t.Error(err)
		}
		require.Equal(t, tt.status, res.StatusCode)
		if tt.status == http.StatusOK {
			require.Contains(t, res.Header.Get("Set-Cookie"), "jwt")
		}
	}
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config, err := config.ParseFlags()
	if err != nil {
		t.Error(err)
	}

	store := mocks.NewMockStore(ctrl)

	gomock.InOrder(
		store.EXPECT().CreateUser(gomock.Any()).Return(int64(1), nil),
		store.EXPECT().CreateUser(gomock.Any()).Return(int64(0), originalStore.ErrDuplicateLogin),
	)

	app := NewApp(config, store, zap.L().Sugar())
	r, err := app.SetupRouter()
	if err != nil {
		t.Error(err)
	}

	srv := httptest.NewServer(r)
	defer srv.Close()

	tests := []struct {
		userCreds models.UserCredentialsSchema
		name      string
		url       string
		method    string
		status    int
	}{
		{
			name:      "Register user",
			userCreds: models.UserCredentialsSchema{Login: "a", Password: "b"},
			url:       "/api/user/register",
			status:    200,
			method:    http.MethodPost,
		},
		{
			name:      "Register user with conflict",
			userCreds: models.UserCredentialsSchema{Login: "a", Password: "b"},
			url:       "/api/user/register",
			status:    409,
			method:    http.MethodPost,
		},
	}

	for _, tt := range tests {
		tt := tt

		b, err := json.Marshal(tt.userCreds)
		if err != nil {
			t.Error(err)
		}

		url, err := url.JoinPath(srv.URL, tt.url)
		if err != nil {
			t.Error(err)
		}

		req, err := http.NewRequest(tt.method, url, bytes.NewBuffer(b))
		if err != nil {
			t.Error(err)
		}
		if err := req.Body.Close(); err != nil {
			t.Error(err)
		}

		res, err := srv.Client().Do(req)
		if err != nil {
			t.Error(err)
		}
		if err := res.Body.Close(); err != nil {
			t.Error(err)
		}
		require.Equal(t, tt.status, res.StatusCode)
		if tt.status == http.StatusOK {
			require.Contains(t, res.Header.Get("Set-Cookie"), "jwt")
		}
	}
}
