package admin

import (
	"errors"
	"mime/multipart"
	"strings"
	"vuecom/gateway/api/v1/handlers"
	"vuecom/gateway/api/v1/request"
	backendusers "vuecom/gateway/api/v1/request/backend_users"
	"vuecom/gateway/api/v1/response"
	"vuecom/gateway/internal/types"
	"vuecom/gateway/internal/utils"
	"vuecom/shared/errors/server"
	dbModels "vuecom/shared/models/db"

	cldApi "github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func DoesOwnerExist(ctx *fiber.Ctx, api *types.Api) (bool, error) {
	db := api.Deps.DB
	_, err := db.BackendUsers().GetAdmin(ctx.Context())
	if err != nil {
		return false, err
	}

	return true, nil
}

// TODO: Validate the business name and the admin route to avoid clashes with url and also storage buckets

func InitializeApp(ctx *fiber.Ctx, api *types.Api) error {
	logger := api.Deps.Logger
	if api.IsAppInit {
		return response.NewResponse(ctx, fiber.StatusBadRequest, "An active app was found!!\nIf you want to initialize a new app, please connect new database")
	}

	db := api.Deps.DB
	cld := api.Deps.Cld
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Error initializing application. Try again later")
	var appData = new(dbModels.AppData)
	// cache
	_data, err := db.AppData().GetAppData(ctx.Context())
	if err != nil {
		logger.Error("Error getting AppData upon app initialization", zap.Error(err))
		return response.FromFiberError(ctx, err500)
	}

	if _data != nil {
		api.IsAppInit = true
		api.AppName = _data.Name
		api.AdminSlug = _data.AdminRoute // Just store the values if the IsAppInit guard does not do anything
		return response.NewResponse(ctx, fiber.StatusBadRequest, "An active app was found!!")
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return response.NewResponse(ctx, fiber.StatusBadRequest, "Invalid form data")
	}

	appData, logoFile, err := validateInitializeProps(form)
	if err != nil {
		return response.NewResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	_, err = cld.Admin.CreateFolder(ctx.Context(), admin.CreateFolderParams{
		Folder: appData.Name,
	})

	if err != nil {
		logger.Error("Error creating application folder", zap.Error(err))
		return response.FromFiberError(ctx, err500)
	}

	fileIO, err := logoFile.Open()
	if err != nil {
		logger.Error("Error opening logo file during initialization", zap.Error(err))
		return response.FromFiberError(ctx, err500)
	}

	upploadRes, err := cld.Upload.Upload(ctx.Context(), fileIO, uploader.UploadParams{
		Folder:      appData.Name,
		Overwrite:   cldApi.Bool(true),
		DisplayName: appData.Name + "_logo",
		PublicID:    appData.Name + "_logo",
	})

	if err != nil {
		logger.Error("Error uploading application logo upon initialization", zap.Error(err))
		return response.FromFiberError(ctx, err500)
	}

	appData.LogoUrl = upploadRes.SecureURL

	// err = gorm.G[dbModels.AppData](db).Create(ctx.Context(), appData)
	err = db.AppData().CreateAppData(appData, ctx.Context())
	if err != nil {
		logger.Error("Error creating app AppData", zap.Error(err))
		return response.FromFiberError(ctx, err500)
	}

	api.IsAppInit = true
	api.AppName = appData.Name
	api.AdminSlug = appData.AdminRoute

	return response.NewResponse(ctx, fiber.StatusOK, "App initialized successfully")
}

func RegisterOwner(ctx *fiber.Ctx, api *types.Api) error {
	logger := api.Deps.Logger
	if api.HasAdmin {
		return response.NewResponse(ctx, fiber.StatusBadRequest, "An existing owner was found!!")
	}

	var err error
	var db = api.Deps.DB
	var cld = api.Deps.Cld

	userExists, err := DoesOwnerExist(ctx, api)
	if err != nil {
		logger.Error("Error checking for existing users", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewResponse(ctx, fiber.StatusBadRequest, "Owner does not exist")
		}
		return response.NewResponse(ctx, fiber.StatusInternalServerError, "An Error occurred, please try again")
	}
	if userExists {
		return response.NewResponse(ctx, fiber.StatusBadRequest, "An existing owner was found!!")
	}

	var reqUser *backendusers.CreateBackendUserRequest
	err = ctx.BodyParser(&reqUser)
	if err != nil {
		return response.NewResponse(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	err = reqUser.Validate()
	if err != nil {
		logger.Error("Validation Error: Owner Registration", zap.Error(err))
		return response.NewResponse(ctx, fiber.StatusBadRequest, "One or more fields do not satisfy the requirements")
	}

	backUser, err := reqUser.ToDBBackendUser(api, ctx.Context())
	if err != nil {
		var serverErr *server.ServerErr
		if errors.As(err, &serverErr) {
			return response.NewResponse(ctx, serverErr.Code, serverErr.Message)
		}

		return response.NewResponse(ctx, fiber.StatusInternalServerError, "An error occurred, please try again")
	}

	reqUserImage, err := ctx.FormFile("image")
	if err != nil {
		return response.NewResponse(ctx, fiber.StatusBadRequest, "Invalid form data")
	}
	if reqUserImage != nil {
		if reqUserImage.Size > handlers.MAX_IMAGE_UPLOAD {
			return response.NewResponse(ctx, fiber.StatusBadRequest, "uploaded image must not be more than 5MB in size")
		}
		fileIO, err := reqUserImage.Open()
		if err != nil {
			return response.NewResponse(ctx, fiber.StatusInternalServerError, "An error occurred, please try again")
		}

		_, err = utils.IsSupportedImage(fileIO)
		if err != nil {
			return response.NewResponse(ctx, fiber.StatusBadRequest, "uploaded image must be either a jpeg, jpg or png image")
		}
		result, err := cld.Upload.Upload(ctx.Context(), fileIO, uploader.UploadParams{
			Folder:      request.GetBackendFolder(api),
			Overwrite:   cldApi.Bool(true),
			DisplayName: backUser.FullName,
			PublicID:    backUser.FullName,
		})
		if err != nil {
			return response.NewResponse(ctx, fiber.StatusInternalServerError, "An error occurred, please try again")
		}
		backUser.Image = &result.SecureURL
	}

	// err = gorm.G[dbModels.BackendUser](db).Create(ctx.Context(), backUser)
	err = db.BackendUsers().CreateUser(backUser, ctx.Context())
	if err != nil {
		return response.NewResponse(ctx, fiber.StatusInternalServerError, "Error registering owner")
	}

	api.HasAdmin = true
	return response.NewResponse(ctx, fiber.StatusOK, "Owner registered successfully")
}

func validateInitializeProps(form *multipart.Form) (appData *dbModels.AppData, file *multipart.FileHeader, err error) {
	fields := form.Value
	files := form.File
	nameField := fields["name"]
	adminRouteField := fields["admin_route"]
	logoFileField := files["app_logo"]

	switch {
	case nameField == nil:
		return nil, nil, errors.New(fieldIsMissing("`Name` field"))
	case adminRouteField == nil:
		return nil, nil, errors.New(fieldIsMissing("`Admin Route` field"))
	case logoFileField == nil:
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
	// fmt.Println(logoFile.Filename, logoFile.Header)

	unknownErr := errors.New("unknown Error Occured! Try again")
	logo, err := logoFile.Open()
	if err != nil {
		return nil, nil, unknownErr
	}

	isSupported, err := utils.IsSupportedImage(logo)
	if err != nil {
		return nil, nil, err
	}
	if !isSupported {
		return nil, nil, errors.New("uploaded logo must be either a jpeg, jpg or png image")
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
