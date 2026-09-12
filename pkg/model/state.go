package model

// AgentState represents universal normalized snapshot for CLI HUD
type AgentState struct {
	Agent       string       `json:"agent"`        // claude, hermes, codex, generic
	Model       string       `json:"model"`        // e.g. claude-3-7-sonnet, anti, gpt-4o
	Context     ContextUsage `json:"context"`      // token context usage
	ActiveTool  string       `json:"active_tool"`  // current or last executed tool
	RunningTask string       `json:"running_task"` // subagent / running step description
	Todos       TodoProgress `json:"todos"`        // task / todo tracker
}

type ContextUsage struct {
	UsedTokens  int64   `json:"used_tokens"`
	TotalTokens int64   `json:"total_tokens"`
	Percentage  float64 `json:"percentage"`
}

type TodoProgress struct {
	Completed int `json:"completed"`
	Total     int `json:"total"`
}
