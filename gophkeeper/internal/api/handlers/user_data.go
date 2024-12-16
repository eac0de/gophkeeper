package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/eac0de/gophkeeper/internal/services"
	"github.com/eac0de/gophkeeper/shared/pkg/httperror"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type UserDataHandlers struct {
	userDataService *services.UserDataService
}

func NewUserDataHandlers(
	userDataService *services.UserDataService,
) *UserDataHandlers {
	return &UserDataHandlers{
		userDataService: userDataService,
	}
}

func (ah *UserDataHandlers) InsertUserAuthInfo(c *gin.Context) {
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	var requestData struct {
		Name     *string                `json:"name"`
		Login    *string                `json:"login"`
		Password *string                `json:"password"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	if requestData.Name == nil || requestData.Login == nil || requestData.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "name,login and password are required"})
		return
	}
	userAuthInfo, err := ah.userDataService.InsertUserAuthInfo(c.Request.Context(), userID, *requestData.Name, *requestData.Login, *requestData.Password, requestData.Metadata)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.JSON(http.StatusCreated, userAuthInfo)
}

func (ah *UserDataHandlers) GetUserAuthInfo(c *gin.Context) {
	dataID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data id"})
		return
	}
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	userAuthInfo, err := ah.userDataService.GetUserAuthInfo(c.Request.Context(), dataID, userID)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.JSON(http.StatusOK, userAuthInfo)
}

func (ah *UserDataHandlers) DeleteUserAuthInfo(c *gin.Context) {
	dataID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data id"})
		return
	}
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	err = ah.userDataService.DeleteUserAuthInfo(c.Request.Context(), dataID, userID)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.String(http.StatusNoContent, "")
}

func (ah *UserDataHandlers) UpdateUserAuthInfo(c *gin.Context) {
	dataID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data id"})
		return
	}
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	var requestData struct {
		Name     *string                `json:"name"`
		Login    *string                `json:"login"`
		Password *string                `json:"password"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	if requestData.Name == nil || requestData.Login == nil || requestData.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "name,login and password are required"})
		return
	}
	userAuthInfo, err := ah.userDataService.UpdateUserAuthInfo(
		c.Request.Context(),
		userID,
		dataID,
		*requestData.Name,
		*requestData.Login,
		*requestData.Password,
		requestData.Metadata,
	)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.JSON(http.StatusOK, userAuthInfo)
}

func (ah *UserDataHandlers) InsertUserTextData(c *gin.Context) {
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	var requestData struct {
		Name     *string                `json:"name"`
		TextData *string                `json:"text_data"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	if requestData.Name == nil || requestData.TextData == nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "name and text_data are required"})
		return
	}
	userTextData, err := ah.userDataService.InsertUserTextData(c.Request.Context(), userID, *requestData.Name, *requestData.TextData, requestData.Metadata)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.JSON(http.StatusCreated, userTextData)
}

func (ah *UserDataHandlers) GetUserTextData(c *gin.Context) {
	dataID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data id"})
		return
	}
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	userTextData, err := ah.userDataService.GetUserTextData(c.Request.Context(), dataID, userID)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.JSON(http.StatusOK, userTextData)
}

func (ah *UserDataHandlers) DeleteUserTextData(c *gin.Context) {
	dataID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data id"})
		return
	}
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	err = ah.userDataService.DeleteUserTextData(c.Request.Context(), dataID, userID)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.String(http.StatusNoContent, "")
}

func (ah *UserDataHandlers) UpdateUserTextData(c *gin.Context) {
	dataID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data id"})
		return
	}
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	var requestData struct {
		Name     *string                `json:"name"`
		TextData *string                `json:"text_data"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	if requestData.Name == nil || requestData.TextData == nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "name and text_data are required"})
		return
	}
	userTextData, err := ah.userDataService.UpdateUserTextData(
		c.Request.Context(),
		userID,
		dataID,
		*requestData.Name,
		*requestData.TextData,
		requestData.Metadata,
	)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.JSON(http.StatusOK, userTextData)
}

func (ah *UserDataHandlers) InsertUserFileData(c *gin.Context) {
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	var requestData struct {
		Metadata map[string]interface{} `json:"metadata"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	dir := fmt.Sprintf("../user_files/%s", userID.String())
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := file.Filename
	pathToFile := fmt.Sprintf("%s/%s", dir, filename)
	count := 0
	for {
		_, err := os.Stat(pathToFile)
		if err != nil {
			if os.IsNotExist(err) {
				break
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		count++
		filename = fmt.Sprintf("%s(%d)", file.Filename, count)
		pathToFile = fmt.Sprintf("%s/%s", dir, filename)
	}

	if err := c.SaveUploadedFile(file, pathToFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userFileData, err := ah.userDataService.InsertUserFileData(
		c.Request.Context(),
		userID,
		file.Filename,
		pathToFile,
		requestData.Metadata,
	)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.JSON(http.StatusCreated, userFileData)
}

func (ah *UserDataHandlers) GetUserFileData(c *gin.Context) {
	dataID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data id"})
		return
	}
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	userFileData, err := ah.userDataService.GetUserFileData(c.Request.Context(), dataID, userID)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.JSON(http.StatusOK, userFileData)
}

func (ah *UserDataHandlers) DeleteUserFileData(c *gin.Context) {
	dataID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data id"})
		return
	}
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	err = ah.userDataService.DeleteUserFileData(c.Request.Context(), dataID, userID)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.String(http.StatusNoContent, "")
}

func (ah *UserDataHandlers) UpdateUserFileData(c *gin.Context) {
	dataID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data id"})
		return
	}
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	var requestData struct {
		Name     *string                `json:"name"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	if requestData.Name == nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "name is required"})
		return
	}
	userFileData, err := ah.userDataService.GetUserFileData(c.Request.Context(), dataID, userID)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	dir, _ := strings.CutSuffix(userFileData.PathToFile, userFileData.Name)
	newPathToFile := fmt.Sprintf("%s/%s", dir, *requestData.Name)
	if err := os.Rename(userFileData.PathToFile, newPathToFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = ah.userDataService.UpdateUserFileData(
		c.Request.Context(),
		userID,
		dataID,
		*requestData.Name,
		newPathToFile,
		requestData.Metadata,
	)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.JSON(http.StatusOK, userFileData)
}

func (ah *UserDataHandlers) InsertUserBankCard(c *gin.Context) {
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	var requestData struct {
		Name       *string `json:"name"`
		Number     *string `json:"number"`
		CardHolder *string `json:"card_holder"`
		ExpireDate *string `json:"expire_date"`
		CSC        *string `json:"csc"`
		Metadata   map[string]interface{}
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	if requestData.Name == nil || requestData.Number == nil || requestData.CardHolder == nil || requestData.ExpireDate == nil || requestData.CSC == nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "name,number,card_holder,expire_date and csc are required"})
		return
	}
	userBankCard, err := ah.userDataService.InsertUserBankCard(
		c.Request.Context(),
		userID,
		*requestData.Name,
		*requestData.Number,
		*requestData.CardHolder,
		*requestData.ExpireDate,
		*requestData.CSC,
		requestData.Metadata,
	)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.JSON(http.StatusCreated, userBankCard)
}

func (ah *UserDataHandlers) GetUserBankCard(c *gin.Context) {
	dataID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data id"})
		return
	}
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	userBankCard, err := ah.userDataService.GetUserBankCard(c.Request.Context(), dataID, userID)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.JSON(http.StatusOK, userBankCard)
}

func (ah *UserDataHandlers) DeleteUserBankCard(c *gin.Context) {
	dataID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data id"})
		return
	}
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	err = ah.userDataService.DeleteUserBankCard(c.Request.Context(), dataID, userID)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.String(http.StatusNoContent, "")
}

func (ah *UserDataHandlers) UpdateUserBankCard(c *gin.Context) {
	dataID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data id"})
		return
	}
	userID := c.MustGet(gin.AuthUserKey).(uuid.UUID)
	var requestData struct {
		Name       *string `json:"name"`
		Number     *string `json:"number"`
		CardHolder *string `json:"card_holder"`
		ExpireDate *string `json:"expire_date"`
		CSC        *string `json:"csc"`
		Metadata   map[string]interface{}
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	if requestData.Name == nil || requestData.Number == nil || requestData.CardHolder == nil || requestData.ExpireDate == nil || requestData.CSC == nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "name,number,card_holder,expire_date and csc are required"})
		return
	}
	userBankCard, err := ah.userDataService.UpdateUserBankCard(
		c.Request.Context(),
		userID,
		dataID,
		*requestData.Name,
		*requestData.Number,
		*requestData.CardHolder,
		*requestData.ExpireDate,
		*requestData.CSC,
		requestData.Metadata,
	)
	if err != nil {
		msg, statusCode := httperror.GetMessageAndStatusCode(err)
		c.JSON(statusCode, gin.H{"detail": msg})
		return
	}
	c.JSON(http.StatusOK, userBankCard)
}
