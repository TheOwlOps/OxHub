#!/bin/bash
# OxHub Universal Installer for Claude Code, Hermes, and OpenAI Codex

set -e

TARGET="${1:-all}"
REPO="https://github.com/TheOwlOps/OxHub.git"

echo "⚡ Installing OxHub Skills & Plugins for: $TARGET"

install_hermes() {
    echo "🦅 Configuring Hermes Agent..."
    if command -v hermes &> /dev/null; then
        hermes skills tap add TheOwlOps/OxHub || true
        echo "✓ Tapped TheOwlOps/OxHub into Hermes!"
    else
        HERMES_DIR="$HOME/.hermes/skills"
        mkdir -p "$HERMES_DIR"
        git clone --depth 1 "$REPO" /tmp/oxhub-tmp 2>/dev/null || (cd /tmp/oxhub-tmp && git pull)
        cp -r /tmp/oxhub-tmp/hub/skills/* "$HERMES_DIR/"
        rm -rf /tmp/oxhub-tmp
        echo "✓ Copied skills to $HERMES_DIR"
    fi
}

install_claude() {
    echo "🤖 Configuring Claude Code..."
    CLAUDE_DIR="$HOME/.claude/skills"
    mkdir -p "$CLAUDE_DIR"
    git clone --depth 1 "$REPO" /tmp/oxhub-tmp 2>/dev/null || (cd /tmp/oxhub-tmp && git pull)
    cp -r /tmp/oxhub-tmp/hub/skills/* "$CLAUDE_DIR/"
    rm -rf /tmp/oxhub-tmp
    echo "✓ Copied skills to $CLAUDE_DIR (Claude Code auto-loads skills in ~/.claude/skills)"
}

install_codex() {
    echo "💻 Configuring OpenAI Codex..."
    CODEX_DIR="$HOME/.codex/skills"
    mkdir -p "$CODEX_DIR"
    git clone --depth 1 "$REPO" /tmp/oxhub-tmp 2>/dev/null || (cd /tmp/oxhub-tmp && git pull)
    cp -r /tmp/oxhub-tmp/hub/skills/* "$CODEX_DIR/"
    rm -rf /tmp/oxhub-tmp
    echo "✓ Copied skills to $CODEX_DIR"
}

case "$TARGET" in
    hermes)
        install_hermes
        ;;
    claude)
        install_claude
        ;;
    codex)
        install_codex
        ;;
    all)
        install_hermes
        install_claude
        install_codex
        ;;
    *)
        echo "Unknown target: $TARGET. Choose: hermes, claude, codex, or all."
        exit 1
        ;;
esac

echo "✅ OxHub setup complete!"
