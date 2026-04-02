module github.com/chibx/vuecom/backend/services/payment

go 1.25.0

require (
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/redis/go-redis/v9 v9.18.0
	go.uber.org/zap v1.27.1
	google.golang.org/grpc v1.80.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
)

require (
	github.com/chibx/vuecom/backend/shared v0.0.0-00010101000000-000000000000
	github.com/goccy/go-json v0.10.6
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.9.1 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260401024825-9d38bb4040a9 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
// vuecom/shared v1.0.0
// github.com/chibx/vuecom/backend/shared v0.0.1
)

replace (
	github.com/chibx/vuecom/backend/services/catalog => ../catalog
	github.com/chibx/vuecom/backend/services/inventory => ../inventory
	github.com/chibx/vuecom/backend/services/orders => ../orders
	github.com/chibx/vuecom/backend/services/payment => ../payment
	github.com/chibx/vuecom/backend/shared => ../../shared
)
