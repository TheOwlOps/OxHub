package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/TheOwlOps/oxhud/pkg/adapter/claude"
	"github.com/TheOwlOps/oxhud/pkg/adapter/codex"
	"github.com/TheOwlOps/oxhud/pkg/adapter/hermes"
	"github.com/TheOwlOps/oxhud/pkg/model"
	"github.com/TheOwlOps/oxhud/pkg/render"
)

func main() {
	agentType := flag.String("agent", "claude", "Agent type: claude, hermes, codex")
	flag.Parse()

	var state *model.AgentState
	var err error

	switch *agentType {
	case "claude":
		state, err = claude.Parse(os.Stdin)
	case "hermes":
		state, err = hermes.Parse(os.Stdin)
	case "codex":
		state, err = codex.Parse(os.Stdin)
	default:
		state, err = claude.Parse(os.Stdin)
	}

	if err != nil {
		// Output minimal fallback on empty or invalid stdin
		state = &model.AgentState{
			Agent: *agentType,
			Model: "unknown",
		}
	}

	output := render.Render(state)
	fmt.Print(output)
}
