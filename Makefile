install:
	go install -mod=readonly $(BUILD_FLAGS) ./blanket

build:
	go build -mod=readonly $(BUILD_FLAGS) -o build/blanket ./blanket