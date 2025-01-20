package agent

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func HandleProxyRequest(path string, w http.ResponseWriter, r *http.Request) error {
	uri := fmt.Sprintf("https://ollama.wantia.app/%s", path)

	req, err := http.NewRequestWithContext(r.Context(), r.Method, uri, r.Body)
	if err != nil {
		return fmt.Errorf("failed to create proxy request: %w", err)
	}

	copyRequestHeaders(r, req)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		if r.Context().Err() != nil {
			return fmt.Errorf("request cancelled by client: %w", r.Context().Err())
		}

		http.Error(w, "Failed to forward request to Ollama", http.StatusBadGateway)
		return fmt.Errorf("error forwarding request: %w", err)
	}
	defer resp.Body.Close()

	copyResponseHeaders(resp, w)

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		if r.Context().Err() != nil {
			return fmt.Errorf("response streaming cancelled by client: %w", r.Context().Err())
		}

		return fmt.Errorf("error copying response: %w", err)
	}

	return nil
}

func copyRequestHeaders(oldReq *http.Request, newReq *http.Request) {
	for name, values := range oldReq.Header {
		for _, value := range values {
			newReq.Header.Add(name, value)
		}
	}
}

func copyResponseHeaders(res *http.Response, w http.ResponseWriter) {
	for name, values := range res.Header {
		for _, value := range values {
			w.Header().Set(name, value)
		}
	}
}
