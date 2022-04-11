package api

import (
	"online_share/common"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PushFile(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, "file upload failed")
		log.Panicln(err)
		return
	}
	files := form.File["files"]
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", common.DefaultSaveDir, file.Filename)
		if err := ctx.SaveUploadedFile(file, path); err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
	}
	ctx.JSON(http.StatusOK, "success")
}
