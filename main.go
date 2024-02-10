package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/robbaen/gtpath/templates"
	"github.com/robbaen/gtpath/templates/components"
)

type Product struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

type ProductsResponse struct {
    Products []Product `json:"items"` // or the actual key name in your JSON
}

func getProducts() ([]Product, error) {
    url := "http://localhost:8090/api/collections/products/records"

    // Make an HTTP GET request
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Error reading response: ", err)
        return nil, err
    }
    defer resp.Body.Close()

    // Read the body of the response using io.ReadAll instead of ioutil.ReadAll
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading body: ", err)
        return nil, err
    }

    // Unmarshal the JSON data into the products slice
    var respObj ProductsResponse
    if err := json.Unmarshal(body, &respObj); err != nil {
        fmt.Println("Error parsing JSON: ", err)
        return nil, err
    }

    // Return the slice of Product objects and no error
    return respObj.Products, nil
}

func main() {
    app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        e.Router.GET("/", func(c echo.Context) error {
			component := templates.Index("TechNytt")
			return component.Render(context.Background(), c.Response().Writer)
	})
		return nil
	})


	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/produkt", func(c echo.Context) error {
			products, err := getProducts()
			if err != nil {
				return err
			}
			for _, product := range products {
				// Render each product
				err := components.Products(product.Name, product.Description).Render(context.Background(), c.Response().Writer)
				if err != nil {
					// Handle error if rendering fails
					return err
				}
			}
			// Return nil if all products are rendered successfully
			return nil
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