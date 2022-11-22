BINARY_NAME=RacoonApp

build:
	@go mod vendor
	@echo "Building GoRacoon..."
	@go build -o tmp/${BINARY_NAME} .
	@echo "GoRacoon has been built!"

run: build
	@echo "Starting GoRacoon..."
	@./tmp/${BINARY_NAME} &
	@echo "GoRacoon has been started!"

clean:
	@echo "Cleaning..."
	@go clean
	@rm tmp/${BINARY_NAME}
	@echo "Cleaning finished!"

test:
	@echo "Testing..."
	@go test ./...
	@echo "Testing done!"

start: run

stop:
	@echo "Stopping GoRacoon..."
	@-pkill -SIGTERM -f "./tmp/${BINARY_NAME}"
	@echo "Stopped GoRacoon!"

restart: stop start