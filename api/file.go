package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const BaseDir = "./upload/"
const baseLength int64 = 512

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

// DownloadFileByName godoc
// @Summary     根据文件名获取指定文件
// @Description  根据文件名获取指定文件
// @Tags         file
// @Accept       json
// @Produce      octet-stream
// @Param        name   path      string  true  "文件名"
// @Success      200  {object}
// @Failure      400  {object} ResponseMessage
// @Failure      416  {object} ResponseMessage
// @Failure      500  {object} ResponseMessage
// @Router       /v1/file/:name [get]
func DownloadFileByName(ctx *gin.Context) {
	var responseMessage ResponseMessage
	name := ctx.Param("name")
	if name == "" {
		responseMessage.Message = "file name cannot be empty. "
		ctx.JSON(http.StatusBadRequest, responseMessage)
		return
	}
	srcPath := BaseDir + name
	file, err := os.Open(srcPath)
	if os.IsNotExist(err) {
		responseMessage.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, responseMessage)
		return
	}
	defer file.Close()
	fs, _ := file.Stat()

	var start, end int64
	r := ctx.GetHeader("Range")
	if r != "" {
		_, err := fmt.Sscanf(r, "bytes=%d-%d", &start, &end)
		if err != nil {
			responseMessage.Message = err.Error()
			ctx.JSON(http.StatusRequestedRangeNotSatisfiable, responseMessage)
			return
		}

		if start > end || start >= fs.Size() || start < 0 {
			responseMessage.Message = "range is not valid. "
			ctx.JSON(http.StatusRequestedRangeNotSatisfiable, responseMessage)
			return
		}
		if end >= fs.Size() {
			end = fs.Size()
		}
	} else {
		end = fs.Size() - 1
	}

	setDownloadHeader(ctx, fs.Size(), name)
	var data []byte
	for {
		if end-start < baseLength {
			data = make([]byte, end-start)
		} else {
			data = make([]byte, baseLength)
		}
		_, err = file.ReadAt(data, start)
		if err != nil {
			responseMessage.Message = err.Error()
			ctx.JSON(http.StatusInternalServerError, responseMessage)
			return
		}

		_, err = ctx.Writer.Write(data)
		if err != nil {
			responseMessage.Message = err.Error()
			ctx.JSON(http.StatusInternalServerError, responseMessage)
			return
		}
		start += baseLength
		if start > end {
			break
		}

	}

}

func setDownloadHeader(ctx *gin.Context, length int64, filename string) {
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Accept-Ranges", "bytes")
	ctx.Header("Content-Length", strconv.FormatInt(length, 10))
}

// DownloadFileByNameHead godoc
// @Summary     获取请求头
// @Description  文件下载接口支持断点续传，此接口可用于获取请求头
// @Tags         file
// @Accept       json
// @Produce      octet-stream
// @Param        name   path      string  true  "文件名"
// @Success      200  {object}
// @Failure      400  {object} ResponseMessage
// @Failure      416  {object} ResponseMessage
// @Failure      500  {object} ResponseMessage
// @Router       /v1/file/:name [head]
func DownloadFileByNameHead(ctx *gin.Context) {
	name := ctx.Param("name")
	src := BaseDir + name
	fs, err := os.Stat(src)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
	}
	setDownloadHeader(ctx, fs.Size(), name)
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
