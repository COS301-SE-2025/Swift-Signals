#!/bin/bash

set -e

echo "Building the Docker image..."
docker build -t frontend:dev -f ../build/package/Dockerfile.frontend ../

echo "Loading the image into Minikube..."
minikube image load frontend:dev

echo "Applying Kubernetes manifests..."
kubectl apply -f development/frontend.yaml

echo "Connecting port 30080 to localhost..."
kubectl port-forward svc/frontend 30080:80 &

echo "Frontend available at http://localhost:30080

