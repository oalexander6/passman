package main

import (
	"log"

	gin_binding "github.com/oalexander6/passman/pkg/bindings/gin"
	"github.com/oalexander6/passman/pkg/config"
	"github.com/oalexander6/passman/pkg/services"
	memory_store "github.com/oalexander6/passman/pkg/stores/memory"
)

func main() {
	conf := config.GetConfig()
	memoryStore := memory_store.New()
	services := services.New(conf, memoryStore.NotesStore)
	app := gin_binding.New(services, conf)

	log.Fatal(app.Run())
}
