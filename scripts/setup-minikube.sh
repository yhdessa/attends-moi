#!/usr/bin/env bash
set -e

echo "🚀 Setting up Minikube..."

# Check if minikube is installed
if ! command -v minikube &> /dev/null; then
    echo "❌ minikube is not installed. Please install it first."
    exit 1
fi

# Start cluster
if ! minikube status &> /dev/null; then
    echo "📦 Starting Minikube cluster..."
    minikube start --driver=docker --memory=1096 --cpus=2
else
    echo "✅ Minikube is already running."
fi

# Enable Ingress addon
echo "🔌 Enabling Ingress addon..."
minikube addons enable ingress

echo "✅ Minikube setup complete!"
