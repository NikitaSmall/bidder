BIN_DIR=bin
BYNARY=bidder

BYNARY_OUTPUT= ${BIN_DIR}/${BYNARY}

all:
	go build -o ${BYNARY_OUTPUT}
