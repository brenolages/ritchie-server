package repository

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"ritchie-server/server"
	"ritchie-server/server/config"
	"ritchie-server/server/mock"
	"testing"
)

func TestHandler_Handler(t *testing.T) {
	type fields struct {
		Config server.Config
		org    string
		method string
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		{
			name: "success",
			fields: fields{
				Config: mock.DummyConfig(),
				org:    "zup",
				method: http.MethodGet,
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-type", "application/json")
					json.NewEncoder(w).Encode(repositoryConfigWant())
				}
			}(),
		},
		{
			name: "not found",
			fields: fields{
				Config: mock.DummyConfig(),
				org:    "zup",
				method: http.MethodPost,
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusNotFound)
				}
			}(),
		},
		{
			name: "not found org",
			fields: fields{
				Config: mock.DummyConfig(),
				org:    "notfound",
				method: http.MethodGet,
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusNotFound)
				}
			}(),
		},
		{
			name: "nil config",
			fields: fields{
				Config: config.Configuration{
					Configs: map[string]*server.ConfigFile{
						"empty": {
							RepositoryConfig: nil,
						}},
				},
				org:    "empty",
				method: http.MethodGet,
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusNotFound)
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := NewConfigHandler(tt.fields.Config)

			r, _ := http.NewRequest(tt.fields.method, "/usage", bytes.NewReader([]byte{}))

			r.Header.Add(server.OrganizationHeader, tt.fields.org)
			r.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			tt.want.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			mu.Handler().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v want %v", g.Code, w.Code)
			}

			if g.Code == http.StatusOK {
				if !reflect.DeepEqual(g.Body, w.Body) {
					t.Errorf("Handler returned wrong body: got %v \n want %v", g.Body, w.Body)
				}
			}
		})
	}
}

func repositoryConfigWant() []server.Repository {
	conf, _ := mock.DummyConfig().ReadRepositoryConfig("zup")
	return conf
}
