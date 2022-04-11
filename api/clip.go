package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.design/x/clipboard"
)

func GetTheLastClip(ctx *gin.Context) {
	text := clipboard.Read(clipboard.FmtText)
	ctx.JSON(http.StatusOK, string(text))
}

func PushNewClip(ctx *gin.Context) {
	text, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "post data is invaild. ")
		return
	}
	clipboard.Write(clipboard.FmtText, text)
	ctx.JSON(http.StatusOK, "success")
}
