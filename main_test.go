package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MyHandler struct {
	body []byte
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(h.body)
}

func Test_computeHash(t *testing.T) {
	testBody := []byte("test")
	srv := httptest.NewServer(&MyHandler{body: testBody})
	defer srv.Close()

	hashRes := fmt.Sprintf("%x", md5.Sum(testBody))

	tests := []struct {
		name    string
		link    string
		want    string
		wantErr bool
	}{
		{
			name:    "first",
			link:    "not exist",
			want:    "",
			wantErr: true,
		},
		{
			name:    "two",
			link:    srv.URL,
			want:    hashRes,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			got, err := computeHash(ctx, md5.New(), tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("computeHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}
