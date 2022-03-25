package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const BaseDir = "./upload/"

// GetFileHistory godoc
// @Summary     获取历史文件
// @Description  获取历史文件
// @Tags         file
// @Accept       json
// @Produce      json
// @Param        page   query      int  true  "页数"
// @Param        size   query      int  true  "每页的条目数"
// @Success      200  {object}  ResponseMessage{message=[]string}
// @Failure      400  {object} ResponseMessage
// @Failure      500  {object} ResponseMessage
// @Router       /v1/file [get]
func GetFileHistory(ctx *gin.Context) {
	var responseMessage ResponseMessage
	files, err := filepath.Glob(BaseDir + "/*")
	if err != nil {
		responseMessage.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, responseMessage)
		return
	}
	pageStr := ctx.DefaultQuery("page", "1")
	sizeStr := ctx.DefaultQuery("size", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		responseMessage.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, responseMessage)
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		responseMessage.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, responseMessage)
		return
	}
	for i := 0; i < len(files); i++ {
		files[i] = filepath.Base(files[i])
	}

	start := (page - 1) * size
	end := page * size
	if start >= len(files) {
		responseMessage.Data = []string{}
	} else {
		responseMessage.Data = files[start:end]
	}
	ctx.JSON(http.StatusOK, responseMessage)
}

// UploadFile godoc
// @Summary     上传文件
// @Description  上传文件
// @Tags         file
// @Accept       mpfd
// @Produce      json
// @Param        file   body   multipart.File     true  "文件"
// @Success      200  {object}  ResponseMessage
// @Failure      400  {object} ResponseMessage
// @Failure      500  {object} ResponseMessage
// @Router       /v1/file [post]
func UploadFile(ctx *gin.Context) {
	var responseMessage ResponseMessage
	file, err := ctx.FormFile("file")
	if err != nil {
		responseMessage.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, responseMessage)
		return
	}
	dst := BaseDir + file.Filename
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		responseMessage.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, responseMessage)
		return
	}
	responseMessage.Message = "ok"
	ctx.JSON(http.StatusOK, responseMessage)
}

// DeleteFileRecord godoc
// @Summary     删除文件
// @Description  删除文件
// @Tags         file
// @Accept       json
// @Produce      json
// @Param        name   path      string  true  "文件名"
// @Success      200  {object}  ResponseMessage{message=string}
// @Failure      400  {object} ResponseMessage
// @Failure      500  {object} ResponseMessage
// @Router       /v1/file/:name [delete]
func DeleteFileRecord(ctx *gin.Context) {
	var responseMessage ResponseMessage
	fileName := BaseDir + strings.TrimSpace(ctx.Param("name"))
	if fileName == BaseDir {
		responseMessage.Message = "file name cannot be empty. "
		ctx.JSON(http.StatusBadRequest, responseMessage)
		return
	}

	err := os.Remove(fileName)
	if err != nil {
		responseMessage.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, responseMessage)
		return
	}

	responseMessage.Data = "ok"
	ctx.JSON(http.StatusOK, responseMessage)
}
