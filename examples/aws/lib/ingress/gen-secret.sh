#!/bin/bash

NAME="secret-name"
SERVER="docker.io"
NAMESPACE="default"

echo -n Registry username: 
read USERNAME
echo

echo -n Registry password: 
read -s PASSWORD
echo

kubectl create secret docker-registry $NAME \
    --dry-run \
    --docker-server=$SERVER \
    --docker-username=$USERNAME \
    --docker-password=$PASSWORD \
    --namespace=$NAMESPACE \
    -o yaml