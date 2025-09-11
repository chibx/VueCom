package main

import (
	"fmt"
	"os"

	// "sync"
	"vuecom/server"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	server.Server
}

func main() {
	port, is_port_set := os.LookupEnv("GO_PORT")
	host, is_host_set := os.LookupEnv("GO_HOST")
	mode, _ := os.LookupEnv("SERVER_MODE")

	var server Server = Server{Server: server.Server{}}

	if len(mode) != 0 {
		fmt.Println("SERVER_MODE:", mode)
	}

	if !is_port_set {
		panic("GO_PORT variable is to be set")
	}

	if !is_host_set {
		panic("GO_HOST variable is to be set")
	}

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusNotFound)
	})

	// For validating the admin slug
	app.Use("/:admin/*", server.ValidateSlug)

	api := app.Group("/api")
	server.LoadApis(api)

	app.Static("/", "./dist")

	app.Listen(fmt.Sprintf("%s:%s", host, port))
}
