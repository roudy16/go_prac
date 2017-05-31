GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

CUR_DIR:=${CURDIR}

BIN=./bin/
SPECS=./specs/

SRC_MAIN=./src/main/
SRC_PWLB=./src/pwlb/
EXEC_STEVE=bin/stt-linux-steve

all: build-steve

build-steve: build-pwlb
	env GOOS=linux GOARCH=amd64 GOPATH=$(CUR_DIR) $(GOBUILD) -o $(EXEC_STEVE) $(SRC_MAIN)

build-pwlb:
	env GOOS=linux GOARCH=amd64 GOPATH=$(CUR_DIR) $(GOBUILD) -i $(SRC_MAIN) $(SRC_PWLB)

test:
	$(GOTEST) $(MAIN)

test-pwlb: build-pwlb
	$(GOTEST) -v $(SRC_PWLB)

runall: build-steve run1 run2 run3 run4 run5

run1: build-steve
	$(EXEC_STEVE) $(CUR_DIR)/$(SPECS)spec1.txt > spec1.out

run2: build-steve
	$(EXEC_STEVE) $(CUR_DIR)/$(SPECS)spec2.txt > spec2.out

run3: build-steve
	$(EXEC_STEVE) $(CUR_DIR)/$(SPECS)spec3.txt > spec3.out

run4: build-steve
	$(EXEC_STEVE) $(CUR_DIR)/$(SPECS)spec4.txt > spec4.out

run5: build-steve
	$(EXEC_STEVE) $(CUR_DIR)/$(SPECS)spec5.txt > spec5.out

clean:
	rm -r $(BIN)
