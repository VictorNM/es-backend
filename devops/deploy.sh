#!/bin/bash
# This script deploy to kubernetes

NAMESPACE="es-dev" # or "es-test" or "es-prod"

echo "Deploying to $NAMESPACE namespace"

echo "Install KUBECTL"
curl -sSL "https://storage.googleapis.com/kubernetes-release/release/v1.14.9/bin/linux/amd64/kubectl" > $TRAVIS_BUILD_DIR/kubectl
chmod +x $TRAVIS_BUILD_DIR/kubectl
export KUBECTL=$TRAVIS_BUILD_DIR/kubectl

echo "Target the cluster"
ibmcloud ks cluster config --cluster $CLUSTER_ID

echo "Applying yaml files"
$KUBECTL --namespace $NAMESPACE apply -f $TRAVIS_BUILD_DIR/devops/yaml_temp/back-end-deployment.yaml
$KUBECTL --namespace $NAMESPACE apply -f $TRAVIS_BUILD_DIR/devops/yaml_temp/back-end-service.yaml
