package v1

import (
	"errors"
	"mime/multipart"
	"strings"
	"vuecom/shared/models"
	"vuecom/shared/models/db"

	cldApi "github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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

// TODO: Validate the business name and the admin route to avoid clashes with url and also storage buckets

func (api *Api) InitializeApp(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	appData, logoFile, err := validateInitializeProps(form)
	if err != nil {
		return ctx.Status(404).SendString(err.Error())
	}
	_, err = api.Cld.Admin.CreateFolder(ctx.Context(), admin.CreateFolderParams{
		Folder: appData.Name,
	})

	err500 := fiber.NewError(500, "Error initializing app. Try again")

	if err != nil {
		return err500
	}

	fileIO, err := logoFile.Open()
	if err != nil {
		return err500
	}

	_, err = api.Cld.Upload.Upload(ctx.Context(), fileIO, uploader.UploadParams{
		Folder:      appData.Name,
		Overwrite:   cldApi.Bool(true),
		DisplayName: appData.Name + "_logo",
		PublicID:    appData.Name + "_logo",
	})

	result := api.DB.Create(appData)
	if result.Error != nil {
		return err500
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

func validateInitializeProps(form *multipart.Form) (appData *models.AppData, file *multipart.FileHeader, err error) {
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

	name := strings.TrimSpace(nameField[0])
	adminRoute := strings.TrimSpace(adminRouteField[0])
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

	appData = &models.AppData{
		Name:       name,
		AdminRoute: adminRoute,
	}

	return appData, logoFile, nil
}

func fieldIsMissing(field string) string {
	return field + " is either missing or invalid"
}
