#!/bin/bash
# Copyright London Stock Exchange Group All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
set -e
#echo "11223"
#sleep 60000000
# This script expedites the chaincode development process by automating the
# requisite channel create/join commands

# We use a pre-generated orderer.block and channel transaction artifact (myc.tx),
# both of which are created using the configtxgen tool

# first we create the channel against the specified configuration in myc.tx
# this call returns a channel configuration block - myc.block - to the CLI container
sleep 3
peer channel create -c ch1 -f ./channel-artifacts/ch1.tx -o orderer:7050

# now we will join the channel and start the chain with myc.block serving as the
# channel's first block (i.e. the genesis block)
peer channel join -b ch1.block

# Now the user can proceed to build and start chaincode in one terminal
# And leverage the CLI container to issue install instantiate invoke query commands in another
peer lifecycle chaincode approveformyorg  -o orderer:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1  --signature-policy "OR ('SampleOrg.member')" --package-id mycc:1.0
peer lifecycle chaincode checkcommitreadiness -o orderer:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1  --signature-policy "OR ('SampleOrg.member')"
peer lifecycle chaincode commit -o orderer:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1  --signature-policy "OR ('SampleOrg.member')" --peerAddresses peer:7051

# sleep 600000000
# exit 0