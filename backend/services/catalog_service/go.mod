module vuecom/catalog

go 1.25.0

replace vuecom/shared => ../../shared

require (
	gorm.io/gorm v1.31.1
	vuecom/shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.31.0 // indirect
)
