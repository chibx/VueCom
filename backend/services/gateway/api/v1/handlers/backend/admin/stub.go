package admin

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
