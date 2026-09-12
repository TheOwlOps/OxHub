package render

import (
	"fmt"
	"strings"

	"github.com/TheOwlOps/oxhud/pkg/model"
)

const (
	ColorReset  = "\033[0m"
	ColorCyan   = "\033[36m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
	ColorGray   = "\033[90m"
	ColorBold   = "\033[1m"
)

func Render(state *model.AgentState) string {
	var b strings.Builder

	// Line 1: Agent & Context Bar
	// e.g.: [CLAUDE: Opus] [████░░░░░░] 45k/200k (22.5%)
	barColor := ColorGreen
	if state.Context.Percentage > 75.0 {
		barColor = ColorYellow
	}
	if state.Context.Percentage > 90.0 {
		barColor = ColorRed
	}

	bar := renderProgressBar(state.Context.Percentage, 10)
	agentLabel := strings.ToUpper(state.Agent)

	b.WriteString(fmt.Sprintf("%s%s[%s: %s]%s %s%s%s %s%s/%s (%.1f%%)%s\n",
		ColorBold, ColorCyan, agentLabel, state.Model, ColorReset,
		barColor, bar, ColorReset,
		ColorGray, formatTokens(state.Context.UsedTokens), formatTokens(state.Context.TotalTokens), state.Context.Percentage, ColorReset,
	))

	// Line 2: Active Tool / Task
	if state.ActiveTool != "" || state.RunningTask != "" {
		toolStr := state.ActiveTool
		if toolStr == "" {
			toolStr = "idle"
		}
		b.WriteString(fmt.Sprintf("%s◐ Tool:%s %s", ColorYellow, ColorReset, toolStr))
		if state.RunningTask != "" {
			b.WriteString(fmt.Sprintf(" %s| Task:%s %s", ColorGray, ColorReset, state.RunningTask))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func renderProgressBar(pct float64, width int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	filled := int((pct / 100.0) * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled
	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}

func formatTokens(n int64) string {
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000.0)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%dk", n/1_000)
	}
	return fmt.Sprintf("%d", n)
}
