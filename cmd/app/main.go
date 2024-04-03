package main

import (
	"floodcontrol/internal/app"
	"fmt"
)

func main() {
	cfg := app.NewConfig()
	fmt.Println(cfg)
	app.Run(cfg)
}
