package v1

import (
	"errors"
	"mime/multipart"
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
// 	LogoUrl   string `json:"app_logo"`

func (api *Api) InitializeApp(ctx *fiber.Ctx) error {
	// appData := models.CreateAppData{}
	// err := ctx.BodyParser(&appData)
	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	appData, logoFile, err := validateInitializeProps(form)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}

	// if err != nil {
	// 	return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	// }

	return nil
}

func RegisterOwner(api *Api, owner *db.BackendUser) error {
	result := api.DB.Create(owner)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func validateInitializeProps(form *multipart.Form) (appData *models.CreateAppData, file *multipart.FileHeader, err error) {
	fields := form.Value
	files := form.File
	nameField := fields["name"]
	adminRouteField := fields["admin_route"]
	logoFileField := files["app_logo"]

	switch {
	case len(nameField) == 0:
		return nil, nil, errors.New(fieldIsMissing("`Name` field"))
	case len(adminRouteField) == 0:
		return nil, nil, errors.New(fieldIsMissing("`Admin Route` field"))
	case len(logoFileField) == 0:
		return nil, nil, errors.New(fieldIsMissing("`Application logo`"))
	}

	name := nameField[0]
	adminRoute := adminRouteField[0]
	logoFile := logoFileField[0]

	switch {
	case len(name) <= 3:
		return nil, nil, errors.New("`Name` should have more than 3 characters")
	case len(adminRoute) < 8:
		return nil, nil, errors.New("`Admin Route` should have at least 8 characters and should be contain alphanumeric characters")
	}

	if logoFile.Size > MAX_IMAGE_UPLOAD {
		return nil, nil, errors.New("The uploaded logo must not be more than 5MB in size")
	}

	appData = &models.CreateAppData{
		Name:       name,
		AdminRoute: adminRoute,
	}

	return appData, logoFile, nil
}

func fieldIsMissing(field string) string {
	return field + " is either missing or invalid"
}
