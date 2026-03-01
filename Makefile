GOC := go

GOFLAGS_PROD := -ldflags "-s -w -X main.IS_PROD=TRUE  -X main.PORT=:8000" -trimpath -buildvcs=false
GOFLAGS_DEV  := -ldflags "-s -w -X main.IS_PROD=FALSE -X main.PORT=:3000" -trimpath -buildvcs=false

BIN_DIR := ./bin
BUILD_DIR := ./build
APPLICATION := mandana
APP_ENTRY := ./cmd/

.PHONY: dev prod build clean

prod:
	mkdir -p $(BIN_DIR)
	$(GOC) build $(GOFLAGS_PROD) -o $(BIN_DIR)/$(APPLICATION) $(APP_ENTRY)

dev:
	mkdir -p $(BIN_DIR)
	$(GOC) build $(GOFLAGS_DEV) -o $(BIN_DIR)/$(APPLICATION)-dev $(APP_ENTRY)

build: prod
	mkdir -p $(BUILD_DIR)
	tar -czvf $(BUILD_DIR)/$(APPLICATION).tar.gz $(BIN_DIR)

clean:
	rm -rf $(BIN_DIR)
