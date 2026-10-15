// hardened-app is the supply-chain lab's tiny HTTP service: a fully static
// (CGO_ENABLED=0) binary meant to run from a shell-less, nonroot runtime
// image. It serves a version banner on / and a liveness probe on /healthz.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
)

// version is stamped at build time via -ldflags "-X main.version=...".
var version = "dev"

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hardened-app %s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	return mux
}

func main() {
	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("hardened-app %s listening on %s", version, addr)
	log.Fatal(http.ListenAndServe(addr, newMux()))
}
