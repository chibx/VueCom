package request

type OnlyID struct {
	ID int `json:"id" params:"id"`
}

var IMAGE_FORMATS = []string{"image/jpeg", "image/png", "image/jpg"}
