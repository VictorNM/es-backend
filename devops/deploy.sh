#!/bin/bash
# This script deploy to kubernetes

NAMESPACE="es-dev" # or "es-test" or "es-prod"

echo "Deploying to $NAMESPACE namespace"

echo "Install KUBECTL"
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt` $TRAVIS_BUILD_DIR/kubectl
chmod +x $TRAVIS_BUILD_DIR/kubectl
export KUBECTL = $TRAVIS_BUILD_DIR/kubectl

echo "Applying yaml files"
$KUBECTL --namespace $NAMESPACE apply -f $TRAVIS_BUILD_DIR/devops/yaml_temp/back-end-deployment.yaml
$KUBECTL --namespace $NAMESPACE apply -f $TRAVIS_BUILD_DIR/devops/yaml_temp/back-end-service.yaml
