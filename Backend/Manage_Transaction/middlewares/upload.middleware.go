package middlewares

import (
	global_types "chatify/types"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UploadFile() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var content_type string = ctx.Get("Content-Type")

		if content_type == "" || strings.Contains(content_type, "multipart/form-data") != true {
			return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
				Message:      "Invalid content type. Expected multipart/form-data.",
				ErrorSection: "Authorization | validate multipart/form-data",
			})
		}

		// form, err := ctx.MultipartForm()
		// if err != nil {
		// 	return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
		// 		Message:      "Error parsing multipart form",
		// 		ErrorSection: "Authorization | parsing multipart/form-data",
		// 	})
		// }

		// var request_file []*multipart.FileHeader = form.File["files"]
		// var all_file []global_types.IFile

		// for i := range request_file {
		// 	var extension string = filepath.Ext(request_file[i].Filename)
		// 	var file_name string = uuid.New().String() + extension

		// 	all_file = append(all_file, global_types.IFile{
		// 		FileName: request_file[i].Filename,
		// 		FilePath: utils.GetFileDirectory(file_name),
		// 		FileType: request_file[i].Header["Content-Type"][0],
		// 	})
		// }

		// ctx.Locals("files", all_file)
		return ctx.Next()
	}
}
