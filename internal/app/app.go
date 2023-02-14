package app

import (
	"HFLabs-test/internal/parser"
	"HFLabs-test/internal/sheetsApi"
	"fmt"
	"log"
)

type App struct {
	parser parser.Parser
}

func New() App {
	return App{
		parser: parser.New(sheetsApi.New()),
	}
}

func (a *App) Run() {
	err := a.parser.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Успешно")
}
