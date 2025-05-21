Write-Host "Building the Docker image..."
docker build --no-cache -t frontend:dev -f ../build/package/Dockerfile.frontend ../

Write-Host "Loading the image into Minikube..."
minikube image load frontend:dev

Write-Host "Applying Kubernetes manifests..."
kubectl apply -f ../deployments/k8s/development/frontend.yaml

Write-Host "Connecting port 30080 to localhost..."
Start-Process -NoNewWindow kubectl port-forward svc/frontend 30080:80

Write-Host "Frontend available at http://localhost:30080"

