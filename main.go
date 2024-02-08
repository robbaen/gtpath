package main

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
    app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        e.Router.GET("/", func(c echo.Context) error {
			return c.String(200, "Hello, World!")
	})
		return nil
	})


    // serves static files from the provided public dir (if exists)
    app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        e.Router.GET("/*", func(c echo.Context) error {
			return c.String(404, "Sorry, this page does not exist!")
	})
		return nil
	})

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}