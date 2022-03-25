package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dianbanjiu/online_share/db"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TextMessage struct {
	Id        int            `json:"id" gorm:"primary key"`
	Content   string         `json:"content" gorm:"type:text;not null" validate:"required"`
	CreatedAt datatypes.Date `json:"created_at" gorm:"type:datetime"`
}

func (message *TextMessage) AfterCreate(db *gorm.DB) (err error) {
	err = db.Model(message).UpdateColumn("created_at", time.Now()).Error
	return
}

type ResponseMessage struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (message *TextMessage) Validate() error {
	valid := validator.New()
	return valid.Struct(message)
}

// GetTextHistory godoc
// @Summary     获取历史文本记录
// @Description  获取历史文本记录
// @Tags         text
// @Accept       json
// @Produce      json
// @Param        page   query      int  true  "页数"
// @Param        size   query      int  true  "每页的条目数"
// @Success      200  {object}  ResponseMessage{message=[]TextMessage}
// @Failure      400  {object} ResponseMessage
// @Failure      500  {object} ResponseMessage
// @Router       /v1/text [get]
func GetTextHistory(ctx *gin.Context) {
	var responseMessage ResponseMessage
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

	var messages = make([]TextMessage, 0)
	if err := db.DB.Offset((page - 1) * size).Limit(size).Order("created_at desc").Find(&messages).Error; err != nil {
		responseMessage.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, responseMessage)
		return
	}
	responseMessage.Data = messages
	ctx.JSON(http.StatusOK, responseMessage)
}

// AddText godoc
// @Summary     新增文本
// @Description  新增文本
// @Tags         text
// @Accept       json
// @Produce      json
// @Param        message   body      TextMessage  true  "文本内容"
// @Success      200  {object}  ResponseMessage
// @Failure      400  {object} ResponseMessage
// @Failure      500  {object} ResponseMessage
// @Router       /v1/text [post]
func AddText(ctx *gin.Context) {
	var message TextMessage
	var responseMessage ResponseMessage
	err := ctx.BindJSON(&message)
	if err != nil {
		responseMessage.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, responseMessage)
		return
	}
	if err := message.Validate(); err != nil {
		responseMessage.Message = "the param is not valid. "
		ctx.JSON(http.StatusBadRequest, responseMessage)
		return
	}

	if err := db.DB.Create(&message).Error; err != nil {
		responseMessage.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, responseMessage)
		return
	}

	responseMessage.Data = "ok"
	ctx.JSON(http.StatusOK, responseMessage)
}

// DeleteTextRecord godoc
// @Summary     删除文本记录
// @Description  删除文本记录
// @Tags         text
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "文本记录的 ID"
// @Success      200  {object}  ResponseMessage{message=TextMessage}
// @Failure      400  {object} ResponseMessage
// @Failure      500  {object} ResponseMessage
// @Router       /v1/text/:id [delete]
func DeleteTextRecord(ctx *gin.Context) {
	var responseMessage ResponseMessage
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responseMessage.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, responseMessage)
		return
	}

	var message TextMessage
	message.Id = id
	if err := db.DB.Delete(&message).Error; err != nil {
		responseMessage.Message = err.Error()
		ctx.JSON(http.StatusInternalServerError, responseMessage)
		return
	}
	responseMessage.Data = message
	ctx.JSON(http.StatusOK, responseMessage)
}
