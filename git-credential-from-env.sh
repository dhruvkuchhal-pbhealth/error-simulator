#!/bin/sh
# Git credential helper - reads GITHUB_TOKEN and GITHUB_OWNER from .env
# Used when: git config credential.helper "$(pwd)/git-credential-from-env.sh"
REPO_ROOT="$(cd "$(dirname "$0")" && pwd)"
if [ -f "$REPO_ROOT/.env" ]; then
  GITHUB_OWNER=$(grep -E "^GITHUB_OWNER=" "$REPO_ROOT/.env" | head -1 | cut -d= -f2- | sed 's/#.*//' | tr -d ' ')
  GITHUB_TOKEN=$(grep -E "^GITHUB_TOKEN=" "$REPO_ROOT/.env" | head -1 | cut -d= -f2- | sed 's/#.*//' | tr -d ' ')
  echo "username=${GITHUB_OWNER}"
  echo "password=${GITHUB_TOKEN}"
fi
