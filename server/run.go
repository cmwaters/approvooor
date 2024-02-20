package server

import (
	"net/http"

	"github.com/cmwaters/blobusign/node"
)

func Run(node node.Node) error {
	mux := http.NewServeMux()
	// Handle PDF submission
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
			return
		}

		

		// Assuming PDF processing happens here
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PDF submitted successfully."))
	})
	return http.ListenAndServe(":8080", mux)
}
