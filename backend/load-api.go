package main

import (
	"vuecom/server/api"

	"github.com/gofiber/fiber/v2"
)

// Potentially Long Function | Just stack all the routes in here
func (s *Server) LoadApis(_api fiber.Router) {
	api := api.Api{S: s.Server}
	_api.Get("/:name", api.ApiHandler)

}
