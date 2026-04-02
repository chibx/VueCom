module github.com/chibx/vuecom/backend/services/gateway

go 1.25.0

// require (
// github.com/chibx/vuecom/backend/services/catalog v0.0.0
// github.com/chibx/vuecom/backend/services/orders v0.0.0
// github.com/chibx/vuecom/backend/services/inventory v0.0.0
// github.com/chibx/vuecom/backend/services/payment v0.0.0
// )

require (
	github.com/chibx/vuecom/backend/services/catalog v0.0.0-00010101000000-000000000000
	github.com/chibx/vuecom/backend/services/inventory v0.0.0-00010101000000-000000000000
	github.com/chibx/vuecom/backend/services/orders v0.0.0-00010101000000-000000000000
	github.com/chibx/vuecom/backend/services/payment v0.0.0-00010101000000-000000000000
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/joho/godotenv v1.5.1
	github.com/pquerna/otp v1.5.0
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
)

require (
	// vuecom/shared v0.0.1
	github.com/chibx/vuecom/backend/shared v0.0.1
	github.com/gabriel-vasile/mimetype v1.4.13
	github.com/go-playground/validator/v10 v10.30.2
	github.com/go-redis/redis_rate/v10 v10.0.1
	github.com/goccy/go-json v0.10.6
	github.com/golang-jwt/jwt/v5 v5.3.1
	go.uber.org/zap v1.27.1
)

require (
	github.com/boombuler/barcode v1.1.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260401024825-9d38bb4040a9 // indirect
)

require (
	github.com/andybalholm/brotli v1.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudinary/cloudinary-go/v2 v2.15.0
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gofiber/fiber/v2 v2.52.12
	github.com/google/uuid v1.6.0
	github.com/gorilla/schema v1.4.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.9.1 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.18.5 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.21 // indirect
	github.com/redis/go-redis/v9 v9.18.0
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.69.0
	golang.org/x/crypto v0.49.0
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
)

replace (
	github.com/chibx/vuecom/backend/services/catalog => ../catalog
	github.com/chibx/vuecom/backend/services/inventory => ../inventory
	github.com/chibx/vuecom/backend/services/orders => ../orders
	github.com/chibx/vuecom/backend/services/payment => ../payment
	github.com/chibx/vuecom/backend/shared => ../../shared
)
