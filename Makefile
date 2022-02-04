TARGET = jaken

BUILD_DIR = ./build
CMD_DIR = ./cmd

all: $(BUILD_DIR)/$(TARGET)

$(BUILD_DIR):
	mkdir -p "$(BUILD_DIR)"

$(BUILD_DIR)/$(TARGET): $(BUILD_DIR)
	go build -v -o "$(BUILD_DIR)/$(TARGET)" "$(CMD_DIR)/$(TARGET)/main.go"

clean:
	rm -rf "$(BUILD_DIR)"