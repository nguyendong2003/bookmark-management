package main

import "github.com/nguyendong2003/bookmark-management/internal/api"

func main() {
	app := api.New()
	if err := app.Start(); err != nil {
		panic(err)
	}
}
