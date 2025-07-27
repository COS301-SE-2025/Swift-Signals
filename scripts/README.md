# Git Workflow Setup for PR-Based Microservice Development

This script sets up Git aliases to streamline a PR-based microservice development workflow with proper branching strategy and automated synchronization.

## Branching Strategy

### Branch Hierarchy
```
main (production)
├── dev (development/staging)
│   ├── microservice (feature integration)
│   │   └── feature/microservice/<feature_name>
│   ├── frontend (feature integration)
│   │   └── feature/frontend/<feature_name>
│   ├── deployment (feature integration)
│   │   └── feature/deployment/<feature_name>
│   └── integration (feature integration)
│       └── feature/integration/<feature_name>
└── hotfix/<critical-security-patch_name> (emergency fixes)
```

### Workflow Overview

1. **Feature Development**: Features are developed in `feature/microservice-name/feature-name` branches
2. **Microservice Integration**: Features are merged into their respective microservice branches via PR
3. **Development Integration**: Microservice branches are merged into `dev` via PR
4. **Production Deployment**: `dev` is merged into `main` via PR
5. **Hotfixes**: Critical fixes branch from `main` and merge back to `main`, then to `dev`

## Installation

1. Navigate to your Git repository
2. Make the script executable:
   ```bash
      chmod +x scripts/setup-git-aliases.sh 
   ```

3. Run the setup script:
   ```bash
   ./scripts/setup-git-aliases.sh
   ```

## Available Commands

### Feature Development

#### `git new-feature <microservice> <feature-name>`
Creates a new feature branch from the specified microservice branch.

**Example:**
```bash
git new-feature user-service login-system
# Creates: feature/user-service/login-system
```

**What it does:**
- Fetches latest changes from origin
- Switches to the user-service branch
- Pulls latest changes
- Creates new feature branch
- Pushes branch to origin with upstream tracking

#### `git feature-sync <microservice> <feature-name>`
Syncs feature branch with its parent microservice branch.

**Example:**
```bash
git feature-sync user-service login-system
```

**What it does:**
- Switches to feature branch
- Fetches latest changes
- Rebases on top of origin/user-service
- Force pushes with lease protection

#### `git prepare-feature-pr <microservice> <feature-name>`
Prepares feature branch for PR creation.

**Example:**
```bash
git prepare-feature-pr user-service login-system
```

**What it does:**
- Runs feature-sync
- Provides PR creation instructions

#### `git delete-feature <microservice> <feature-name>`
Safely deletes feature branch locally and remotely.

**Example:**
```bash
git delete-feature user-service login-system
```

**What it does:**
- Checks if branch exists
- Ensures you're not on the branch being deleted
- Deletes local branch
- Deletes remote branch if it exists

### Microservice Integration

#### `git microservice-sync-dev <microservice>`
Syncs microservice branch with dev branch.

**Example:**
```bash
git microservice-sync-dev user-service
```

**What it does:**
- Fetches latest changes
- Switches to user-service branch
- Merges origin/dev
- Pushes changes

#### `git prepare-microservice-pr <microservice>`
Prepares microservice branch for PR to dev.

**Example:**
```bash
git prepare-microservice-pr user-service
```

### Hotfix Management

#### `git new-hotfix <hotfix-name>`
Creates a new hotfix branch from main.

**Example:**
```bash
git new-hotfix security-patch-001
```

**What it does:**
- Fetches latest changes
- Switches to main
- Pulls latest main
- Creates hotfix branch
- Pushes with upstream tracking

#### `git hotfix-sync <hotfix-name>`
Syncs hotfix branch with main.

**Example:**
```bash
git hotfix-sync security-patch-001
```

#### `git prepare-hotfix-pr <hotfix-name>`
Prepares hotfix for PR to main.

**Example:**
```bash
git prepare-hotfix-pr security-patch-001
```

#### `git delete-hotfix <hotfix-name>`
Safely deletes hotfix branch.

**Example:**
```bash
git delete-hotfix security-patch-001
```

### Utility Commands

#### `git graph`
Shows a visual representation of the git history.

**Example:**
```bash
git graph
```

## Complete Workflow Examples

### Feature Development Workflow

```bash
# 1. Start new feature
git new-feature intersection-service GetAllIntersections
# 2. Make your changes and commits
git add .
git commit -m "Add GetAllIntersections endpoint"

# 3. Sync with parent branch before PR
git prepare-feature-pr intersection-service GetAllIntersections

# 4. Open PR on GitHub: feature/intersection-service/GetAllIntersections → intersection-service

# 5. After PR is merged, cleanup
git delete-feature intersection-service GetAllIntersections
```

### Microservice Integration Workflow

```bash
# 1. Sync microservice with dev
git prepare-microservice-pr intersection-service

# 2. Open PR on GitHub: intersection-service → dev

# 3. After PR is merged, pull latest dev
git checkout dev
git pull origin dev
```

### Hotfix Workflow

```bash
# 1. Create hotfix
git new-hotfix critical-auth-bug

# 2. Make fix and commit
git add .
git commit -m "Fix authentication bypass vulnerability"

# 3. Sync and prepare PR
git prepare-hotfix-pr critical-auth-bug

# 4. Open PR on GitHub: hotfix/critical-auth-bug → main

# 5. After merge, cleanup and sync dev
git delete-hotfix critical-auth-bug
git checkout dev
git pull origin dev
git merge origin/main  # Sync hotfix to dev
git push origin dev
```

## Best Practices

### Branch Naming
- Features: `feature/microservice-name/descriptive-feature-name`
- Hotfixes: `hotfix/descriptive-issue-name`
- Microservices: Use either `user-service`, `intersection-service`,
                 `simulation-service`, `optimisation-service`, `metric-service`,
                 `frontend`, `deployment` or `integration`

### Commit Messages
- Use conventional commits format
- Be descriptive and specific
- Reference issue numbers when applicable

### PR Management
- Always sync branches before creating PRs
- Use descriptive PR titles and descriptions
- Review code thoroughly before merging
- Delete feature branches after merging

## Error Handling

The script includes comprehensive error handling:

- **Branch existence checks**: Prevents creating duplicate branches
- **Current branch protection**: Prevents deleting the branch you're on
- **Remote branch verification**: Handles cases where remote branches don't exist
- **Validation**: Ensures required parameters are provided

## Troubleshooting

### Common Issues

1. **Branch already exists**: Use `git branch -a` to see all branches
2. **Cannot delete current branch**: Switch to a different branch first
3. **Merge conflicts**: Resolve conflicts manually during sync operations
4. **Permission denied**: Ensure you have push access to the repository

### Cleanup Commands

```bash
# List all local branches
git branch

# List all remote branches
git branch -r

# Delete local branch (force)
git branch -D branch-name

# Delete remote branch manually
git push origin --delete branch-name
```

## Customization

To modify the aliases, edit the script and re-run it. The aliases are stored in `.git/config` and can be viewed with:

```bash
git config --local --get-regexp alias
```

## Security Considerations

- Uses `--force-with-lease` instead of `--force` for safer force pushes
- Includes branch existence checks to prevent accidental operations
- Validates required parameters before execution

## Support

For issues or questions:
1. Check the troubleshooting section
2. Review Git logs for specific error messages
3. Ensure you have proper repository permissions
4. Verify branch names and microservice names are correct
