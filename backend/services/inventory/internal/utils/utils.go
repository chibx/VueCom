package utils

import (
	"github.com/chibx/vuecom/backend/services/inventory/internal/global"
	"go.uber.org/zap"
)

func FailOnError(err error, msg string) {
	if err != nil {
		global.Logger.Fatal(msg, zap.Error(err))
	}
}
