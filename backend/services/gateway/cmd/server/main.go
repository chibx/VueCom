package main

import (
	"fmt"
	"time"

	v1 "github.com/chibx/vuecom/backend/services/gateway/api/v1"
	"github.com/chibx/vuecom/backend/services/gateway/config"
	"github.com/chibx/vuecom/backend/services/gateway/internal/global"
	igrpc "github.com/chibx/vuecom/backend/services/gateway/internal/grpc"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

func main() {
	config := config.GetConfig()
	v1_api := &types.Api{Config: config, Deps: &types.Deps{}}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		StreamRequestBody:     true,
	})

	app.Use(helmet.New())
	initServer(app, v1_api)
	stopConns := igrpc.InitClients()
	logger := global.Logger

	defer func() {
		_ = logger.Sync()
	}()

	defer stopConns()

	v1.LoadRoutes(app, v1_api)

	logger.Info(fmt.Sprintf("Server is listening on http://%s:%s\n", config.Host, config.Port), zap.Time("now", time.Now()))
	logger.Fatal("Error starting server:", zap.Error(app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))))

}
