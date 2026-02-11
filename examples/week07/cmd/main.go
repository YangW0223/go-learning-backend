package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/yang/go-learning-backend/examples/week07"
)

func main() {
	auth := week07.NewAuthService("week07-secret", time.Hour)
	_ = auth.Register("alice", "pass123", "user")

	token, err := auth.Login("alice", "pass123")
	if err != nil {
		fmt.Println("login error:", err)
		return
	}
	fmt.Println("issued token:", token)

	protected := week07.AuthMiddleware(auth, "user", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("access granted"))
	})

	okReq := httptest.NewRequest(http.MethodGet, "/secure", nil)
	okReq.Header.Set("Authorization", "Bearer "+token)
	okRec := httptest.NewRecorder()
	protected(okRec, okReq)
	fmt.Printf("authorized => status=%d body=%s\n", okRec.Code, okRec.Body.String())

	badReq := httptest.NewRequest(http.MethodGet, "/secure", nil)
	badRec := httptest.NewRecorder()
	protected(badRec, badReq)
	fmt.Printf("missing token => status=%d body=%s", badRec.Code, badRec.Body.String())
}
