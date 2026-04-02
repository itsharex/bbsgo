.PHONY: all build dev clean

all: build

build:
	@echo "Building frontend..."
	cd site && npm install && npm run build
	@echo "Building backend..."
	cd server && go mod tidy && go build -o bbsgo.exe .
	@echo "Build complete!"

dev:
	@echo "Starting development servers..."
	@echo "Please run the following commands in separate terminals:"
	@echo "  Terminal 1: cd server && go run main.go"
	@echo "  Terminal 2: cd site && npm run dev"

clean:
	rm -f server/bbsgo.exe
	rm -f server/bbsgo.db
	rm -rf site/node_modules
	rm -rf site/dist

install:
	cd server && go mod download
	cd site && npm install

run-server:
	cd server && go run main.go

run-site:
	cd site && npm run dev
