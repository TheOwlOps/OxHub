package claude

import (
	"encoding/json"
	"io"

	"github.com/TheOwlOps/oxhud/pkg/model"
)

// ClaudeStatusLineInput represents JSON passed via stdin by Claude Code statusline
type ClaudeStatusLineInput struct {
	Model struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"model"`
	ContextWindow struct {
		CurrentUsage struct {
			InputTokens  int64 `json:"input_tokens"`
			OutputTokens int64 `json:"output_tokens"`
		} `json:"current_usage"`
		ContextWindowSize int64 `json:"context_window_size"`
	} `json:"context_window"`
	TranscriptPath string `json:"transcript_path"`
}

// Parse reads stdin JSON from Claude Code
func Parse(r io.Reader) (*model.AgentState, error) {
	var input ClaudeStatusLineInput
	dec := json.NewDecoder(r)
	if err := dec.Decode(&input); err != nil {
		return nil, err
	}

	modelName := input.Model.DisplayName
	if modelName == "" {
		modelName = input.Model.ID
	}
	if modelName == "" {
		modelName = "claude"
	}

	used := input.ContextWindow.CurrentUsage.InputTokens + input.ContextWindow.CurrentUsage.OutputTokens
	total := input.ContextWindow.ContextWindowSize
	var pct float64
	if total > 0 {
		pct = float64(used) / float64(total) * 100.0
	}

	state := &model.AgentState{
		Agent: "claude",
		Model: modelName,
		Context: model.ContextUsage{
			UsedTokens:  used,
			TotalTokens: total,
			Percentage:  pct,
		},
	}
	return state, nil
}
