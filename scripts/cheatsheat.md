# Git Aliases Workflow CheatSheet

### Feature Development Workflow

```bash
# 1. Start new feature
git new-feature '<microservice-name>' '<feature-name>'
# 2. Make your changes and commits
git add .
git commit -m "Meaningful message"

# 3. Sync with parent branch before PR
git prepare-feature-pr '<microservice-name>' '<feature-name>'

# 4. Open PR on GitHub: feature/'<microservice-name>'/'<feature-name>' → '<microservice-name>'

# 5. After PR is merged, cleanup
git delete-feature '<microservice-name>' '<feature-name>'
```

### Microservice Integration Workflow

```bash
# 1. Sync microservice with dev
git prepare-microservice-pr '<microservice-name>'

# 2. Open PR on GitHub: '<microservice-name>' → dev

# 3. After PR is merged, pull latest dev
git checkout dev
git pull origin dev
```

### Hotfix Workflow

```bash
# 1. Create hotfix
git new-hotfix '<hotfix-name>'

# 2. Make fix and commit
git add .
git commit -m "Meaningful message"

# 3. Sync and prepare PR
git prepare-hotfix-pr '<hotfix-name>'

# 4. Open PR on GitHub: hotfix/'<hotfix-name>' → main

# 5. After merge, cleanup and sync dev
git delete-hotfix '<hotfix-name>'
git checkout dev
git pull origin dev
git merge origin/main  # Sync hotfix to dev
git push origin dev
```


