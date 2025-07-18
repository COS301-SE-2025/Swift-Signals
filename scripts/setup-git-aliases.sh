#!/bin/bash

# Git Workflow Setup Script
# Sets up local Git aliases for PR-based microservice workflow

set -e  # Exit on any error

echo "🔧 Setting up local Git aliases for your PR-based microservice workflow..."

# Function to check if a branch exists locally
check_local_branch_exists() {
    git show-ref --verify --quiet "refs/heads/$1"
}

# Function to check if a branch exists on remote
check_remote_branch_exists() {
    git ls-remote --heads origin "$1" | grep -q "$1"
}

# Create new feature branch
git config --local alias.new-feature '!f() { \
  ms=${1:?Microservice name required}; \
  fn=${2:?Feature name required}; \
  branch_name="feature/$ms/$fn"; \
  if git show-ref --verify --quiet "refs/heads/$branch_name"; then \
    echo "❌ Branch $branch_name already exists locally"; \
    exit 1; \
  fi; \
  if git ls-remote --heads origin "$branch_name" | grep -q "$branch_name"; then \
    echo "❌ Branch $branch_name already exists on remote"; \
    exit 1; \
  fi; \
  git fetch origin && \
  git switch "$ms" && \
  git pull origin "$ms" && \
  git switch -c "$branch_name" && \
  git push -u origin "$branch_name" && \
  echo "✅ Created and pushed feature branch: $branch_name"; \
}; f'

# Sync feature branch with microservice branch
git config --local alias.feature-sync '!f() { \
  ms=${1:?Microservice name required}; \
  fn=${2:?Feature name required}; \
  branch_name="feature/$ms/$fn"; \
  if ! git show-ref --verify --quiet "refs/heads/$branch_name"; then \
    echo "❌ Branch $branch_name does not exist locally"; \
    exit 1; \
  fi; \
  git switch "$branch_name" && \
  git fetch origin && \
  git rebase "origin/$ms" && \
  git push --force-with-lease && \
  echo "✅ Synced feature branch with $ms"; \
}; f'

# Prepare feature for PR
git config --local alias.prepare-feature-pr '!f() { \
  ms=${1:?Microservice name required}; \
  fn=${2:?Feature name required}; \
  git feature-sync "$ms" "$fn" && \
  echo "✅ Now open a PR: feature/$ms/$fn → $ms on Github"; \
}; f'

# Delete feature branch (with existence check)
git config --local alias.delete-feature '!f() { \
  ms=${1:?Microservice name required}; \
  fn=${2:?Feature name required}; \
  branch_name="feature/$ms/$fn"; \
  if ! git show-ref --verify --quiet "refs/heads/$branch_name"; then \
    echo "❌ Branch $branch_name does not exist locally"; \
    exit 1; \
  fi; \
  current_branch=$(git rev-parse --abbrev-ref HEAD); \
  if [ "$current_branch" = "$branch_name" ]; then \
    echo "❌ Cannot delete branch $branch_name - you are currently on it"; \
    exit 1; \
  fi; \
  git branch -d "$branch_name" && \
  if git ls-remote --heads origin "$branch_name" | grep -q "$branch_name"; then \
    git push origin --delete "$branch_name" && \
    echo "✅ Deleted feature branch: $branch_name (local and remote)"; \
  else \
    echo "✅ Deleted feature branch: $branch_name (local only - remote did not exist)"; \
  fi; \
}; f'

# Sync microservice branch with dev
git config --local alias.microservice-sync-dev '!f() { \
  ms=${1:?Microservice name required}; \
  git fetch origin && \
  git switch "$ms" && \
  git pull origin "$ms" && \
  git merge "origin/dev" && \
  git push origin "$ms" && \
  echo "✅ Synced $ms with dev branch"; \
}; f'

# Prepare microservice for PR to dev
git config --local alias.prepare-microservice-pr '!f() { \
  ms=${1:?Microservice name required}; \
  git microservice-sync-dev "$ms" && \
  echo "✅ Now open a PR: $ms → dev on Github"; \
}; f'

# Create new hotfix branch
git config --local alias.new-hotfix '!f() { \
  hn=${1:?Hotfix name required}; \
  branch_name="hotfix/$hn"; \
  if git show-ref --verify --quiet "refs/heads/$branch_name"; then \
    echo "❌ Branch $branch_name already exists locally"; \
    exit 1; \
  fi; \
  if git ls-remote --heads origin "$branch_name" | grep -q "$branch_name"; then \
    echo "❌ Branch $branch_name already exists on remote"; \
    exit 1; \
  fi; \
  git fetch origin && \
  git switch main && \
  git pull origin main && \
  git switch -c "$branch_name" && \
  git push -u origin "$branch_name" && \
  echo "✅ Created and pushed hotfix branch: $branch_name"; \
}; f'

# Sync hotfix branch with main
git config --local alias.hotfix-sync '!f() { \
  hn=${1:?Hotfix name required}; \
  branch_name="hotfix/$hn"; \
  if ! git show-ref --verify --quiet "refs/heads/$branch_name"; then \
    echo "❌ Branch $branch_name does not exist locally"; \
    exit 1; \
  fi; \
  git switch "$branch_name" && \
  git fetch origin && \
  git pull origin "$branch_name" && \
  git merge "origin/main" && \
  git push origin "$branch_name" && \
  echo "✅ Synced hotfix branch with main"; \
}; f'

# Prepare hotfix for PR
git config --local alias.prepare-hotfix-pr '!f() { \
  hn=${1:?Hotfix name required}; \
  git hotfix-sync "$hn" && \
  echo "✅ Now open a PR: hotfix/$hn → main on Github"; \
}; f'

# Delete hotfix branch (with existence check and fixed variable name)
git config --local alias.delete-hotfix '!f() { \
  hn=${1:?Hotfix name required}; \
  branch_name="hotfix/$hn"; \
  if ! git show-ref --verify --quiet "refs/heads/$branch_name"; then \
    echo "❌ Branch $branch_name does not exist locally"; \
    exit 1; \
  fi; \
  current_branch=$(git rev-parse --abbrev-ref HEAD); \
  if [ "$current_branch" = "$branch_name" ]; then \
    echo "❌ Cannot delete branch $branch_name - you are currently on it"; \
    exit 1; \
  fi; \
  git branch -d "$branch_name" && \
  if git ls-remote --heads origin "$branch_name" | grep -q "$branch_name"; then \
    git push origin --delete "$branch_name" && \
    echo "✅ Deleted hotfix branch: $branch_name (local and remote)"; \
  else \
    echo "✅ Deleted hotfix branch: $branch_name (local only - remote did not exist)"; \
  fi; \
}; f'

# Enhanced git graph visualization
git config --local alias.graph \
"log --oneline --graph --decorate -n 10"

echo "✅ Git aliases configured locally:"
git config --local --get-regexp alias

echo ""
echo "🎉 Setup complete! Run 'git <alias-name> --help' or check the README for usage instructions."
