#!/usr/bin/env bash
set -e

echo "🏗️ Building Docker images..."

# Build Backend
echo "   📦 Building backend..."
docker build -t attends-moi-backend:latest -f apps/backend/Dockerfile apps/backend

# Build Frontend
echo "   📦 Building frontend..."
docker build -t attends-moi-frontend:latest -f apps/frontend/Dockerfile apps/frontend

echo "📥 Loading images into Minikube..."
minikube image load attends-moi-backend:latest
minikube image load attends-moi-frontend:latest

echo "📝 Applying Kubernetes manifests..."
kubectl apply -k k8s/overlays/minikube

echo "⏳ Waiting for deployments to be ready..."
kubectl rollout status deployment/backend -n attends-moi --timeout=120s
kubectl rollout status deployment/frontend -n attends-moi --timeout=120s
kubectl rollout status deployment/postgres -n attends-moi --timeout=120s

echo ""
echo "✅ Deployment successful!"
echo "🌐 Access the app:"
minikube service frontend -n attends-moi --url
