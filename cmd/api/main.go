package main

import "github.com/nguyendong2003/bookmark-management/internal/api"

func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	app := api.New(cfg)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
