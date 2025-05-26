#!/bin/bash

set -e

echo "Building the Docker image..."
docker build --no-cache -t swift-signals-frontend:latest -f ../build/package/Dockerfile.frontend ../

echo "Loading the image into Minikube..."
minikube image load swift-signals-frontend:latest

echo "Applying Kubernetes manifests..."
kubectl apply -f ../deployments/k8s/development/frontend.yaml

echo "Connecting port 30080 to localhost..."
kubectl port-forward svc/frontend 30080:80 

echo "Frontend available at http://localhost:30080"

