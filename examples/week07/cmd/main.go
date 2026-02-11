// 详细注释: package main
package main

// 详细注释: import (
import (
	// 详细注释: "fmt"
	"fmt"
	// 详细注释: "net/http"
	"net/http"
	// 详细注释: "net/http/httptest"
	"net/http/httptest"
	// 详细注释: "time"
	"time"

	// 详细注释: "github.com/yang/go-learning-backend/examples/week07"
	"github.com/yang/go-learning-backend/examples/week07"
	// 详细注释: )
)

// 详细注释: func main() {
func main() {
	// 详细注释: auth := week07.NewAuthService("week07-secret", time.Hour)
	auth := week07.NewAuthService("week07-secret", time.Hour)
	// 详细注释: _ = auth.Register("alice", "pass123", "user")
	_ = auth.Register("alice", "pass123", "user")

	// 详细注释: token, err := auth.Login("alice", "pass123")
	token, err := auth.Login("alice", "pass123")
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: fmt.Println("login error:", err)
		fmt.Println("login error:", err)
		// 详细注释: return
		return
		// 详细注释: }
	}
	// 详细注释: fmt.Println("issued token:", token)
	fmt.Println("issued token:", token)

	// 详细注释: protected := week07.AuthMiddleware(auth, "user", func(w http.ResponseWriter, _ *http.Request) {
	protected := week07.AuthMiddleware(auth, "user", func(w http.ResponseWriter, _ *http.Request) {
		// 详细注释: w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusOK)
		// 详细注释: _, _ = w.Write([]byte("access granted"))
		_, _ = w.Write([]byte("access granted"))
		// 详细注释: })
	})

	// 详细注释: okReq := httptest.NewRequest(http.MethodGet, "/secure", nil)
	okReq := httptest.NewRequest(http.MethodGet, "/secure", nil)
	// 详细注释: okReq.Header.Set("Authorization", "Bearer "+token)
	okReq.Header.Set("Authorization", "Bearer "+token)
	// 详细注释: okRec := httptest.NewRecorder()
	okRec := httptest.NewRecorder()
	// 详细注释: protected(okRec, okReq)
	protected(okRec, okReq)
	// 详细注释: fmt.Printf("authorized => status=%d body=%s\n", okRec.Code, okRec.Body.String())
	fmt.Printf("authorized => status=%d body=%s\n", okRec.Code, okRec.Body.String())

	// 详细注释: badReq := httptest.NewRequest(http.MethodGet, "/secure", nil)
	badReq := httptest.NewRequest(http.MethodGet, "/secure", nil)
	// 详细注释: badRec := httptest.NewRecorder()
	badRec := httptest.NewRecorder()
	// 详细注释: protected(badRec, badReq)
	protected(badRec, badReq)
	// 详细注释: fmt.Printf("missing token => status=%d body=%s", badRec.Code, badRec.Body.String())
	fmt.Printf("missing token => status=%d body=%s", badRec.Code, badRec.Body.String())
	// 详细注释: }
}
