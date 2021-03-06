#!/bin/bash
# This script push image to docker registry

set -o errexit

echo "---Install IBM Cloud CLI and registry plugins"
curl -fsSL https://clis.cloud.ibm.com/install/linux | sh
curl -sL https://ibm.biz/idt-installer | bash

echo "---Login into IBM Cloud CLI"
ibmcloud api cloud.ibm.com
ibmcloud login --apikey $IBMCLOUD_API_KEY -c $IBMCLOUD_ACC_ID --no-region
docker login -u iamapikey -p $IBMCLOUD_API_KEY us.icr.io

echo "---Build Docker Image"
docker build . -f $TRAVIS_BUILD_DIR/Dockerfile -t esbackend:$TRAVIS_BUILD_NUMBER

echo "---Tag docker image with IBM Cloud"
docker tag esbackend:$TRAVIS_BUILD_NUMBER us.icr.io/esregistry/esbackend:$TRAVIS_BUILD_NUMBER

echo "---Push docker image with IBM Cloud Registry"
docker push us.icr.io/esregistry/esbackend:$TRAVIS_BUILD_NUMBER
