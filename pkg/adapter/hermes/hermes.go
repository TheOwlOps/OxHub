package hermes

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/TheOwlOps/oxhud/pkg/model"
)

type HermesSessionConfig struct {
	Model         string `json:"model"`
	ContextLength int64  `json:"context_length"`
}

// Parse reads JSON from stdin or detects active Hermes session
func Parse(r io.Reader) (*model.AgentState, error) {
	state := &model.AgentState{
		Agent: "hermes",
		Model: "hermes",
		Context: model.ContextUsage{
			UsedTokens:  0,
			TotalTokens: 200000,
			Percentage:  0,
		},
	}

	// Try reading stdin if piped
	var input struct {
		Model       string `json:"model"`
		UsedTokens  int64  `json:"used_tokens"`
		TotalTokens int64  `json:"total_tokens"`
		ActiveTool  string `json:"active_tool"`
		Task        string `json:"task"`
	}

	if r != nil {
		if err := json.NewDecoder(r).Decode(&input); err == nil {
			if input.Model != "" {
				state.Model = input.Model
			}
			state.Context.UsedTokens = input.UsedTokens
			if input.TotalTokens > 0 {
				state.Context.TotalTokens = input.TotalTokens
			}
			if state.Context.TotalTokens > 0 {
				state.Context.Percentage = float64(state.Context.UsedTokens) / float64(state.Context.TotalTokens) * 100.0
			}
			state.ActiveTool = input.ActiveTool
			state.RunningTask = input.Task
			return state, nil
		}
	}

	// Fallback: Check local config if available
	home, err := os.UserHomeDir()
	if err == nil {
		cfgPath := filepath.Join(home, ".hermes", "config.yaml")
		if _, err := os.Stat(cfgPath); err == nil {
			state.Model = "active"
		}
	}

	return state, nil
}
