package models

import "time"

type AppData struct {
	Name       string `json:"app_name" gorm:"" redis:"name"`
	AdminRoute string `json:"-" gorm:"" redis:"admin_route"`
	LogoUrl    string `json:"app_logo" gorm:"" redis:"logo_url"`
	// Plan       int    `json:"app_plan"`
}

type OnlyID struct {
	ID int `json:"id"`
}

// JWT Format sent back to the client dashboard
type BackendJWTPayload struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	Exp    int    `json:"exp"`
}

// JWT Format sent back to the client
type CustomerJWTPayload struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	Exp    int    `json:"exp"`
}

// type ApiCustomer struct {
// 	sharedUserProps
// }

// // Base Backend Panel User
// type ApiBackendUser struct {
// 	sharedUserProps
// 	Role string `json:"role" gorm:"type:varchar(50)"`
// }

type ApiProducts struct {
}

type CustomerReviews struct {
	Id        uint      `gorm:"primary"`
	Text      string    `gorm:"text"`
	Rating    int8      `gorm:"rating"`
	CreatedAt time.Time ``
	EditTimes int8      `gorm:""`
	UserId    uint      `gorm:""`
	ProductId uint      `gorm:""`
}

type f struct {
}
