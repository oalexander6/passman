clean:
	rm -rf ./tmp ./dist

build: templates styles scripts
	go build -o ./dist/passman ./cmd/main.go

templates:
	templ generate -path=./pkg

templates-fmt:
	templ fmt ./pkg

templates-watch:
	templ generate -path=./pkg --watch

styles:
	tailwindcss -i ./assets/styles.css -o ./dist/public/assets/styles.css

styles-watch:
	tailwindcss -i ./assets/styles.css -o ./dist/public/assets/styles.css --watch