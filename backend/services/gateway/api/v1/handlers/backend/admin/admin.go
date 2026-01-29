package admin

import (
	"errors"
	"mime/multipart"

	"github.com/chibx/vuecom/backend/shared/errors/server"
	appModels "github.com/chibx/vuecom/backend/shared/models/db/appdata"
	"github.com/valyala/fasthttp"

	backendusers "github.com/chibx/vuecom/backend/services/gateway/api/v1/request/backend_users"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"

	// userModels "github.com/chibx/vuecom/backend/shared/models/db/users"

	cldApi "github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func DoesOwnerExist(ctx *fiber.Ctx, api *types.Api) (bool, error) {
	// db := api.Deps.DB
	// has, err := db.BackendUsers().HasAdmin(ctx.Context())
	// if err != nil {
	// 	return false, err
	// }

	// return has, nil
	return api.HasAdmin, nil
}

// TODO: Validate the business name and the admin route to avoid clashes with url and also storage buckets

func InitializeApp(api *types.Api) fiber.Handler {
	logger := utils.Logger()

	return func(ctx *fiber.Ctx) error {
		if api.IsAppInit {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "An active app was found!!\nIf you want to initialize a new app, please connect new database")
		}

		db := api.Deps.DB
		cld := api.Deps.Cld
		err500 := fiber.NewError(fiber.StatusInternalServerError, "Error initializing application. Try again later")
		var appData = new(appModels.AppData)
		// cache
		/* I might not need this since i check on server startup and i might have a background task refresh the local variable */
		// ----------------------------------------------------
		// _data, err := db.AppData().GetAppData(ctx.Context())
		// if err != nil {
		// 	logger.Error("Error getting AppData upon app initialization", zap.Error(err))
		// 	return response.FromFiberError(ctx, err500)
		// }

		// if _data != nil {
		// 	api.IsAppInit = true
		// 	api.AppName = _data.AppName
		// 	// api.AdminSlug = _data.AdminRoute // Just store the values if the IsAppInit guard does not do anything
		// 	return response.WriteResponse(ctx, fiber.StatusBadRequest, "An active app was found!!")
		// }
		// ----------------------------------------------------

		appData, logoFile, err := validateInitializeProps(ctx)
		if err != nil {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, err.Error())
		}
		_, err = cld.Admin.CreateFolder(ctx.Context(), admin.CreateFolderParams{
			Folder: appData.AppName,
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
			Folder:      appData.AppName,
			Overwrite:   cldApi.Bool(true),
			DisplayName: appData.AppName + "_logo",
			PublicID:    appData.AppName + "_logo",
		})

		if err != nil {
			logger.Error("Error uploading application logo upon initialization", zap.Error(err))
			return response.FromFiberError(ctx, err500)
		}

		appData.LogoUrl = upploadRes.SecureURL
		appData.Settings = appModels.GetDefaultAppSettings()

		err = db.AppData().CreateAppData(ctx.Context(), appData)
		if err != nil {
			logger.Error("Error creating app AppData", zap.Error(err))
			return response.FromFiberError(ctx, err500)
		}

		api.IsAppInit = true
		api.AppName = appData.AppName
		// api.AdminSlug = appData.AdminRoute

		return response.WriteResponse(ctx, fiber.StatusOK, "App initialized successfully")
	}
}

func RegisterOwner(api *types.Api) fiber.Handler {
	logger := utils.Logger()

	return func(ctx *fiber.Ctx) error {
		if api.HasAdmin {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "An existing owner was found!!")
		}

		var err error
		var db = api.Deps.DB
		// var cld = api.Deps.Cld
		err500 := fiber.NewError(fiber.StatusInternalServerError, "An error occurred, please try again")

		userExists, err := DoesOwnerExist(ctx, api)
		if err != nil {
			logger.Error("Error checking for existing users", zap.Error(err))
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.WriteResponse(ctx, fiber.StatusBadRequest, "Owner does not exist")
			}
			return response.FromFiberError(ctx, err500)
		}
		if userExists {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "An existing owner was found!!")
		}

		var reqUser *backendusers.CreateBackendUserRequest
		err = ctx.BodyParser(&reqUser)
		if err != nil {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Invalid request body")
		}

		err = reqUser.Validate()
		if err != nil {
			logger.Error("Validation Error: Owner Registration", zap.Error(err))
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "One or more fields do not satisfy the requirements")
		}

		backUser, err := reqUser.ToDBBackendUser(ctx.Context(), api, ctx)
		if err != nil {
			var serverErr *server.ServerErr
			if errors.As(err, &serverErr) {
				return response.WriteResponse(ctx, serverErr.Code, serverErr.Message)
			}

			return response.FromFiberError(ctx, err500)
		}

		// reqUserImage, err := ctx.FormFile("image")
		// if err != nil {
		// 	return response.WriteResponse(ctx, fiber.StatusBadRequest, "Invalid form data")
		// }
		// if reqUserImage != nil {
		// 	if reqUserImage.Size > constants.MaxImageUpload {
		// 		return response.WriteResponse(ctx, fiber.StatusBadRequest, "uploaded image must not be more than 5MB in size")
		// 	}
		// 	fileIO, err := reqUserImage.Open()
		// 	if err != nil {
		// 		return response.FromFiberError(ctx, err500)
		// 	}

		// 	if reqUserImage.Size > constants.MaxImageUpload {
		// 		return response.WriteResponse(ctx, fiber.StatusBadRequest, "Uploaded file must not be more than 5MB in size")
		// 	}
		// 	_, err = utils.IsSupportedImage(fileIO)
		// 	if err != nil {
		// 		return response.WriteResponse(ctx, fiber.StatusBadRequest, "Uploaded image must be either a jpeg, jpg or png image")
		// 	}
		// 	result, err := cld.Upload.Upload(ctx.Context(), fileIO, uploader.UploadParams{
		// 		Folder:      request.GetBackendFolder(api),
		// 		Overwrite:   cldApi.Bool(true),
		// 		DisplayName: backUser.FullName,
		// 		PublicID:    backUser.FullName,
		// 	})
		// 	if err != nil {
		// 		return response.FromFiberError(ctx, err500)
		// 	}
		// 	backUser.Image = &result.SecureURL
		// }

		// err = gorm.G[dbModels.BackendUser](db).Create(ctx.Context(), backUser)
		err = db.BackendUsers().CreateUser(ctx.Context(), backUser)
		if err != nil {
			return response.WriteResponse(ctx, fiber.StatusInternalServerError, "Error registering owner")
		}

		api.HasAdmin = true
		return response.WriteResponse(ctx, fiber.StatusOK, "Owner registered successfully")
	}
}

func validateInitializeProps(ctx *fiber.Ctx) (*appModels.AppData, *multipart.FileHeader, error) {
	var err error
	var appData *appModels.AppData

	name := ctx.FormValue("name")
	logoFile, err := ctx.FormFile("app_logo")
	if err != nil {
		if errors.Is(err, fasthttp.ErrMissingFile) {
			return nil, nil, errors.New("Logo is missing")
		}

		return nil, nil, errors.New("Something went wrong, please try again")
	}

	if logoFile == nil {
		return nil, nil, errors.New(fieldIsMissing("`Application logo`"))
	}

	if len(name) <= 3 {
		return nil, nil, errors.New("`Name` should have more than 3 characters")
	}

	if logoFile.Size > constants.MaxImageUpload {
		return nil, nil, errors.New("uploaded logo must not be more than 5MB in size")
	}

	unknownErr := errors.New("Something went wrong, please try again")
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
	appData = &appModels.AppData{
		AppName: name,
	}

	return appData, logoFile, nil
}

func fieldIsMissing(field string) string {
	return field + " is either missing or invalid"
}
