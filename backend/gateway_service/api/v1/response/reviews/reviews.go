package reviews

import "time"

type CustomerReviews struct {
	Id        uint      `gorm:"primary"`
	Text      string    `gorm:"text"`
	Rating    int8      `gorm:"rating"`
	CreatedAt time.Time ``
	EditTimes int8      `gorm:""`
	UserId    uint      `gorm:""`
	ProductId uint      `gorm:""`
}
