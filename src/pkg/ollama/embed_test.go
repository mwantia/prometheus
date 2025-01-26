package ollama

import (
	"context"
	"net/http"
	"testing"

	"github.com/mwantia/prometheus/internal/config"
)

func TestEmbed(tst *testing.T) {
	cfg, err := config.ParseConfig("../../tests/config.hcl")
	if err != nil {
		tst.Fatalf("Unable to parse config: %v", err)
	}

	model := "mxbai-embed-large:latest"
	ctx := context.TODO()

	client := CreateClient(cfg.Ollama.Endpoint, model, http.DefaultClient)

	tst.Run("Ollama.Embed", func(t *testing.T) {
		resp, err := client.Embed(ctx, EmbedRequest{
			Input: "The apply is orange",
		})
		if err != nil {
			tst.Errorf("Unable to perform chat request: %v", err)
		}

		tst.Log(resp.Embeddings)
	})
}
