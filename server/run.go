package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/cmwaters/blobusign/node"
)

// Run runs a http server over the node
func Run(n node.Node) error {
	mux := http.NewServeMux()
	// Handle document submission
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		id, err := n.Publish(r.Context(), body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to publish document: %s", err.Error()), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to write response: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	})

	// Handle document retrieval
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
			return
		}

		rawID := r.URL.Query().Get("id")
		id, err := node.Parse([]byte(rawID))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse ID: %s", err.Error()), http.StatusBadRequest)
			return
		}

		doc, err := n.Get(r.Context(), node.ID(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve document: %s", err.Error()), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(doc); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write document: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	})

	// Handle document signing
	mux.HandleFunc("/sign", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
			return
		}

		rawID := r.URL.Query().Get("id")
		id, err := node.Parse([]byte(rawID))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse ID: %s", err.Error()), http.StatusBadRequest)
			return
		}

		err = n.Sign(r.Context(), node.ID(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to sign document: %s", err.Error()), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(":8080", mux)
}
