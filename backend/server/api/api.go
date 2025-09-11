package api

import (
	"fmt"
	"vuecom/server"

	"github.com/gofiber/fiber/v2"
)

type Api struct {
	S server.Server
}

func (a *Api) ApiHandler(ctx *fiber.Ctx) error {
	fmt.Println(ctx.Params("name"))
	return nil
}
