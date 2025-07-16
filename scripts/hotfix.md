### Hotfix Workflow

```bash
# 1. Create hotfix
git new-hotfix '<hotfix-name>'

# 2. Make fix and commit
git add .
git commit -m "Meaningful message"

# 3. Sync with main to get latest changes
git hotfix-sync '<hotfix-name>'

# 4. Make more changes and commits
git add .
git commit -m "Meaningful message"

# 5. Sync with main before PR
git prepare-hotfix-pr '<hotfix-name>'

# 4. Open PR on GitHub: 
hotfix/'<hotfix-name>' → main

# 5. After merge, cleanup and sync dev
git delete-hotfix '<hotfix-name>'
git checkout dev
git pull origin dev
git merge origin/main  # Sync hotfix to dev
git push origin dev
```

