### Feature Development Workflow

```bash
# 1. Start new feature
git new-feature '<microservice-name>' '<feature-name>'

# 2. Make your changes and commits
git add .
git commit -m "Meaningful message"

# 3. Sync with parent branch to get latest changes
git feature-sync '<microservice>' '<feature-name>'

# 4. Make more changes and commits
git add .
git commit -m "Meaningful message"

# 5. Sync with parent branch before PR
git prepare-feature-pr '<microservice-name>' '<feature-name>'

# 6. Open PR on GitHub: 
feature/'<microservice-name>'/'<feature-name>' → '<microservice-name>'

# 7. After PR is merged, cleanup
git delete-feature '<microservice-name>' '<feature-name>'
```

