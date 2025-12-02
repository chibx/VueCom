package admin

import (
	"errors"
	"mime/multipart"
	"strings"
	"vuecom/gateway/api/v1/handlers"
	"vuecom/gateway/internal/v1/types"
	dbModels "vuecom/shared/models/db"

	cldApi "github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DoesOwnerExist(ctx *fiber.Ctx, api *types.Api) (bool, error) {
	db := api.Deps.DB
	backendUser, err := gorm.G[dbModels.BackendUser](db).Select("role").Where("role = 'owner'").Limit(1).Find(ctx.Context())
	if err != nil {
		return false, err
	}

	if len(backendUser) == 0 {
		return false, nil
	}

	if backendUser[0].Role != "owner" {
		return false, nil
	}

	return true, nil
}

// TODO: Validate the business name and the admin route to avoid clashes with url and also storage buckets

func InitializeApp(ctx *fiber.Ctx, api *types.Api) error {
	db := api.Deps.DB
	cld := api.Deps.Cld
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Error initializing app. Try again")
	var appData = new(dbModels.AppData)
	// cache
	_data, err := gorm.G[dbModels.AppData](db).First(ctx.Context())
	if err != nil {
		return err500
	}
	if _data.Name != "" {
		return fiber.NewError(fiber.StatusBadRequest, "An active app was found!!")
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return err500
	}

	appData, logoFile, err := validateInitializeProps(form)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	_, err = cld.Admin.CreateFolder(ctx.Context(), admin.CreateFolderParams{
		Folder: appData.Name,
	})

	if err != nil {
		return err500
	}

	fileIO, err := logoFile.Open()
	if err != nil {
		return err500
	}

	upploadRes, err := cld.Upload.Upload(ctx.Context(), fileIO, uploader.UploadParams{
		Folder:      appData.Name,
		Overwrite:   cldApi.Bool(true),
		DisplayName: appData.Name + "_logo",
		PublicID:    appData.Name + "_logo",
	})

	if err != nil {
		return err500
	}

	appData.LogoUrl = upploadRes.SecureURL

	err = gorm.G[dbModels.AppData](db).Create(ctx.Context(), appData)
	if err != nil {
		return err500
	}

	return nil
}

func RegisterOwner(ctx *fiber.Ctx, api *types.Api) error {
	db := api.Deps.DB
	owner := &dbModels.BackendUser{}

	err := gorm.G[dbModels.BackendUser](db).Create(ctx.Context(), owner)
	if err != nil {
		return err
	}
	return nil
}

func validateInitializeProps(form *multipart.Form) (appData *dbModels.AppData, file *multipart.FileHeader, err error) {
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

	if logoFile.Size > handlers.MAX_IMAGE_UPLOAD {
		return nil, nil, errors.New("uploaded logo must not be more than 5MB in size")
	}

	appData = &dbModels.AppData{
		Name:       name,
		AdminRoute: adminRoute,
	}

	return appData, logoFile, nil
}

func fieldIsMissing(field string) string {
	return field + " is either missing or invalid"
}
