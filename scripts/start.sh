#!/bin/bash
. scripts/scriptUtils.sh

COMPOSE_FILES=docker/docker-compose-template.yaml
COMPOSE_FILE_CHAINCODE=docker/chaincode-template.yaml

if [[ $# -eq 1 ]]; then
  if [ "$1" == "-c" ]; then
    if [ -z "${CC_SRC_PATH}" ]; then
      source .env
      if [ -z "${CC_SRC_PATH}" ]; then
        fatalln 'The environment variable CC_SRC_PATH is not exist, please edit .env or add environment variable to specific it'
      fi
    fi
    COMPOSE_FILES="${COMPOSE_FILES} -f ${COMPOSE_FILE_CHAINCODE}"
  fi
fi

infoln "Service start..."
docker-compose -f ${COMPOSE_FILES} up -d
res=$?
if [ $res -ne 0 ]; then
  fatalln "Failed to start network..."
fi

