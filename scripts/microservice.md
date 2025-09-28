### Microservice Integration Workflow

```bash
# 1. Update microservice branch with latest dev
git microservice-sync-dev '<microservice>'

# 2. Prepare microservice branch for PR with dev
git prepare-microservice-pr '<microservice-name>'

# 2. Open PR on GitHub: 
'<microservice-name>' → dev

# 3. After PR is merged, pull latest dev
git switch dev
git pull origin dev
```

