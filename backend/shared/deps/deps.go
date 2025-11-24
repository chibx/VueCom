package deps

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Deps struct {
	DB    *gorm.DB
	Redis *redis.Client
	Cld   *cloudinary.Cloudinary
}
