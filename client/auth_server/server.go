package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os/exec"
)

var port = flag.String("port", "8080", "server port")

func handler(w http.ResponseWriter, r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	fmt.Printf("Generating identity token from gcloud for request\n%s\n", string(dump))
	cmd := exec.Command("gcloud", "auth", "print-identity-token")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch ID token from gcloud, error: %v", err), http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", out.String()))
	fmt.Println("ok")
	io.WriteString(w, "ok")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Printf("Authz server is listening on port %s\n", *port)
	http.ListenAndServe(":"+*port, nil)
}
