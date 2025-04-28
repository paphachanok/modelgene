package ollama

import (
	"github.com/ollama/ollama/api"
	"github.com/paphachanok/modelgene/pkg/types"
)

func convertMessages(msgs []types.Message) []api.Message {
	var out []api.Message
	for _, m := range msgs {
		out = append(out, api.Message{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return out
}
