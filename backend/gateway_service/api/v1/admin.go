package v1

import (
	"vuecom/shared/models"
	"vuecom/shared/models/db"

	"github.com/gofiber/fiber/v2"
)

func (api *Api) DoesOwnerExist(ctx *fiber.Ctx) (bool, error) {
	backendUser := &db.BackendUser{}
	result := api.DB.Select("role").Where("role = 'owner'").First(backendUser)
	if result.Error != nil {
		return false, result.Error
	}

	if backendUser.Role != "owner" {
		return false, nil
	}

	return true, nil
}

//  Name       string `json:"app_name"`
// 	AdminRoute string `json:"admin_route"`
// 	Plan       int    `json:"app_plan"`
// 	ImageUrl   string `json:"app_mage_url"`

func (api *Api) InitializeApp(ctx *fiber.Ctx) error {
	appData := models.CreateAppData{}
	err := ctx.BodyParser(&appData)
	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	fields := form.Value
	files := form.File
	name := fields["name"]
	adminRoute := fields["admin_route"]
	plan := fields["plan"]

	switch {
	case len(name) == 0:
	case len(adminRoute) == 0:
	case len(plan) == 0:
	}

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return nil
}

func RegisterOwner(api *Api, owner *db.BackendUser) error {
	result := api.DB.Create(owner)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
