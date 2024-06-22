clean:
	rm -rf ./tmp ./dist

dev: templates styles
	echo "To enable hot reloading of templates, run 'make templates-watch' in another terminal"
	air

build: templates styles
	go build -o ./dist/passman ./cmd/main.go

run:
	./dist/passman

templates: templates-format
	templ generate -path=./pkg

templates-format:
	templ fmt ./pkg

templates-watch:
	templ generate -path=./pkg -watch

styles:
	tailwindcss -i ./pkg/pages/styles/styles.css -o ./dist/public/assets/styles.css

styles-watch:
	tailwindcss -i ./pkg/pages/styles/styles.css -o ./dist/public/assets/styles.css --watch