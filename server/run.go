package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cmwaters/blobusign/node"
)

type Node interface {
	Publish(ctx context.Context, data []byte) (node.ID, error)
	Get(ctx context.Context, id node.ID) (node.SignedDocument, error)
	Sign(ctx context.Context, id node.ID) error
}

// Run runs a http server over the node
func Run(n Node) error {
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

		enableCORS(w)
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

		docBytes, err := json.Marshal(doc)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal document: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		enableCORS(w)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(docBytes); err != nil {
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

		enableCORS(w)
		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(":8080", mux)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
