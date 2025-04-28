package anthropic

import "github.com/anthropics/anthropic-sdk-go"

func getMaxTokens(max *int) int64 {
	if max != nil {
		return int64(*max)
	}
	return 1024
}

func mergeContent(blocks []anthropic.ContentBlockUnion) string {
	var combined string
	for _, block := range blocks {
		if block.Text != "" {
			combined += block.Text
		}
	}
	return combined
}
