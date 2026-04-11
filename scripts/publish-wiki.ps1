Param(
    [string]$Repo = "Variel42k/flowscope",
    [string]$SourceDir = "docs/wiki",
    [string]$TempDir = ".tmp/wiki-publish",
    [string]$Token = ""
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $SourceDir)) {
    throw "Source directory not found: $SourceDir"
}

if (Test-Path $TempDir) {
    Remove-Item -LiteralPath $TempDir -Recurse -Force
}

$effectiveToken = $Token
if ([string]::IsNullOrWhiteSpace($effectiveToken) -and $env:GITHUB_TOKEN) {
    $effectiveToken = $env:GITHUB_TOKEN
}

$wikiUrl = "https://github.com/$Repo.wiki.git"
if (-not [string]::IsNullOrWhiteSpace($effectiveToken)) {
    $wikiUrl = "https://x-access-token:$effectiveToken@github.com/$Repo.wiki.git"
}

Write-Host "Cloning wiki repository for: $Repo"
git clone $wikiUrl $TempDir
if ($LASTEXITCODE -ne 0 -or -not (Test-Path (Join-Path $TempDir ".git"))) {
    throw "Failed to clone wiki repository. Ensure GitHub Wiki is enabled and authentication is valid."
}

Write-Host "Copying wiki pages from $SourceDir"
Copy-Item -Path "$SourceDir/*" -Destination $TempDir -Recurse -Force

Push-Location $TempDir
try {
    git add .
    $changes = git status --porcelain
    if (-not $changes) {
        Write-Host "No wiki changes to publish."
        exit 0
    }

    git config user.name "flowscope-wiki-bot"
    git config user.email "wiki-bot@users.noreply.github.com"
    $branch = (git rev-parse --abbrev-ref HEAD).Trim()
    if ([string]::IsNullOrWhiteSpace($branch) -or $branch -eq "HEAD") {
        $branch = "master"
    }
    git commit -m "docs(wiki): sync from docs/wiki"
    git push origin $branch
    Write-Host "Wiki published successfully."
}
finally {
    Pop-Location
}
