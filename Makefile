.PHONY: start stop reload
stop:
	@scripts/stop.sh
start:
	@scripts/start.sh
chaincode-reload:
	@scripts/reload.sh
start-chaincode:
	@scripts/start.sh -c