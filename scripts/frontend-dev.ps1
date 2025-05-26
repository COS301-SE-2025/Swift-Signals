Write-Host "Building the Docker image..."
docker build --no-cache -t swift-signals-frontend:latest -f ../build/package/Dockerfile.frontend ../

Write-Host "Loading the image into Minikube..."
minikube image loadswift-signals-frontend:latest 

Write-Host "Applying Kubernetes manifests..."
kubectl apply -f ../deployments/k8s/development/frontend.yaml

Write-Host "Connecting port 30080 to localhost..."
kubectl port-forward svc/frontend 30080:80

Write-Host "Frontend available at http://localhost:30080"

