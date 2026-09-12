param(
    [string]$Target = "all"
)

$Repo = "https://github.com/TheOwlOps/OxHub.git"
Write-Host "⚡ Installing OxHub Skills & Plugins for: $Target" -ForegroundColor Cyan

function Install-Hermes {
    Write-Host "🦅 Configuring Hermes Agent..." -ForegroundColor Green
    if (Get-Command hermes -ErrorAction SilentlyContinue) {
        hermes skills tap add TheOwlOps/OxHub
    } else {
        $hermesDir = "$HOME\.hermes\skills"
        New-Item -ItemType Directory -Force -Path $hermesDir | Out-Null
        git clone --depth 1 $Repo "$env:TEMP\oxhub-tmp"
        Copy-Item -Recurse -Force "$env:TEMP\oxhub-tmp\hub\skills\*" $hermesDir
        Remove-Item -Recurse -Force "$env:TEMP\oxhub-tmp"
    }
}

function Install-Claude {
    Write-Host "🤖 Configuring Claude Code..." -ForegroundColor Green
    $claudeDir = "$HOME\.claude\skills"
    New-Item -ItemType Directory -Force -Path $claudeDir | Out-Null
    git clone --depth 1 $Repo "$env:TEMP\oxhub-tmp"
    Copy-Item -Recurse -Force "$env:TEMP\oxhub-tmp\hub\skills\*" $claudeDir
    Remove-Item -Recurse -Force "$env:TEMP\oxhub-tmp"
}

function Install-Codex {
    Write-Host "💻 Configuring OpenAI Codex..." -ForegroundColor Green
    $codexDir = "$HOME\.codex\skills"
    New-Item -ItemType Directory -Force -Path $codexDir | Out-Null
    git clone --depth 1 $Repo "$env:TEMP\oxhub-tmp"
    Copy-Item -Recurse -Force "$env:TEMP\oxhub-tmp\hub\skills\*" $codexDir
    Remove-Item -Recurse -Force "$env:TEMP\oxhub-tmp"
}

switch ($Target) {
    "hermes" { Install-Hermes }
    "claude" { Install-Claude }
    "codex"  { Install-Codex }
    "all"    { Install-Hermes; Install-Claude; Install-Codex }
    default  { Write-Host "Unknown target: $Target" -ForegroundColor Red; exit 1 }
}

Write-Host "✅ OxHub setup complete!" -ForegroundColor Green
