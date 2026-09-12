package codex

import (
	"encoding/json"
	"io"

	"github.com/TheOwlOps/oxhud/pkg/model"
)

type CodexStatusInput struct {
	Model       string `json:"model"`
	UsedTokens  int64  `json:"used_tokens"`
	TotalTokens int64  `json:"total_tokens"`
	ActiveTool  string `json:"active_tool"`
	Status      string `json:"status"`
}

func Parse(r io.Reader) (*model.AgentState, error) {
	state := &model.AgentState{
		Agent: "codex",
		Model: "gpt-4o",
		Context: model.ContextUsage{
			UsedTokens:  0,
			TotalTokens: 128000,
			Percentage:  0,
		},
	}

	if r == nil {
		return state, nil
	}

	var input CodexStatusInput
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
		state.RunningTask = input.Status
	}

	return state, nil
}
