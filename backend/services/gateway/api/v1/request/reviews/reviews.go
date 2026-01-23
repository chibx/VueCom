package reviews

type CustomerCreateReviews struct {
	Text      string `json:"text"`
	Rating    int8   `json:"rating"`
	ProductId uint   `json:"product_id"`
}

type CustomerEditReviews struct {
	Id     uint   `json:"id"`
	Text   string `json:"text"`
	Rating int8   `json:"rating"`
}
