#!/bin/bash
. scripts/scriptUtils.sh

function stop() {
    docker-compose -f docker/docker-compose-template.yaml -f docker/chaincode-template.yaml down --volumes --remove-orphans
}
infoln "Service stop..."

stop