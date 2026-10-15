// drill-app is hardened-app with two deliberately vulnerable dependencies,
// the live-drill definition of the dhc catalogue (task 16.2). It serves the
// same version banner on / and liveness probe on /healthz, behind a CORS
// middleware, plus /whoami, which reads the bearer token's subject.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
)

// version is stamped at build time via -ldflags "-X main.version=...".
var version = "dev"

// signingKey verifies the drill's bearer tokens; never a production secret.
var signingKey = []byte(os.Getenv("DRILL_JWT_KEY"))

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "drill-app %s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("/whoami", func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		tok, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) { return signingKey, nil })
		if err != nil || !tok.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		sub, _ := tok.Claims.GetSubject()
		fmt.Fprintf(w, "hello %s\n", sub)
	})
	return mux
}

// newHandler wraps the mux in the CORS middleware the drill pins.
func newHandler() http.Handler {
	return handlers.CORS(handlers.AllowedOrigins([]string{"https://drill.example"}))(newMux())
}

func main() {
	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("drill-app %s listening on %s", version, addr)
	log.Fatal(http.ListenAndServe(addr, newHandler()))
}
