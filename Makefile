clean:
	rm -rf ./tmp ./dist

dev:
	air

build:
	go build -o ./dist/passman ./cmd/main.go

run:
	./dist/passman