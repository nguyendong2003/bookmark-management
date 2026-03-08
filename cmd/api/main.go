package main

import "github.com/nguyendong2003/bookmark-management/internal/api"

func main() {
	cfg, err := api.NewConfig("BOOKMARK_SERVICE")
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
