package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/eac0de/gophkeeper/internal/models"
	"github.com/eac0de/gophkeeper/shared/pkg/httperror"
	"github.com/go-playground/validator/v10"

	"github.com/google/uuid"
)

type IUserDataStore interface {
	InsertUserTextData(ctx context.Context, data *models.UserTextData) error
	InsertUserFileData(ctx context.Context, data *models.UserFileData) error
	InsertUserAuthInfo(ctx context.Context, data *models.UserAuthInfo) error
	InsertUserBankCard(ctx context.Context, data *models.UserBankCard) error

	UpdateUserTextData(ctx context.Context, data *models.UserTextData) error
	UpdateUserFileData(ctx context.Context, data *models.UserFileData) error
	UpdateUserAuthInfo(ctx context.Context, data *models.UserAuthInfo) error
	UpdateUserBankCard(ctx context.Context, data *models.UserBankCard) error

	GetUserTextData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserTextData, error)
	GetUserFileData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserFileData, error)
	GetUserAuthInfo(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserAuthInfo, error)
	GetUserBankCard(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserBankCard, error)

	DeleteUserTextData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error
	DeleteUserFileData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error
	DeleteUserAuthInfo(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error
	DeleteUserBankCard(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error
}

type UserDataService struct {
	store IUserDataStore
}

func NewUserDataService(userDataStore IUserDataStore) *UserDataService {
	return &UserDataService{
		store: userDataStore,
	}
}

func (uds *UserDataService) InsertUserTextData(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	text string,
	metadata map[string]interface{},
) (*models.UserTextData, error) {
	userTextData := models.NewUserTextData(name, userID, metadata, text)
	err := uds.store.InsertUserTextData(ctx, &userTextData)
	if err != nil {
		return nil, err
	}

	return &userTextData, nil
}

func (uds *UserDataService) InsertUserFileData(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	pathToFile string,
) (*models.UserFileData, error) {
	userFileData := models.NewUserFileData(name, userID, pathToFile)
	err := uds.store.InsertUserFileData(ctx, &userFileData)
	if err != nil {
		return nil, err
	}

	return &userFileData, nil
}

func (uds *UserDataService) InsertUserAuthInfo(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	login, password string,
	metadata map[string]interface{},
) (*models.UserAuthInfo, error) {
	userAuthInfo := models.NewUserAuthInfo(name, userID, metadata, login, password)
	err := uds.store.InsertUserAuthInfo(ctx, &userAuthInfo)
	if err != nil {
		return nil, err
	}
	return &userAuthInfo, nil
}

func (uds *UserDataService) InsertUserBankCard(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	number, cardHolder, expireDate, csc string,
	metadata map[string]interface{},
) (*models.UserBankCard, error) {
	userBankCard := models.NewUserBankCard(name, userID, metadata, number, cardHolder, expireDate, csc)
	v := validator.New()
	err := v.Struct(userBankCard)
	if err != nil {
		msg := ""
		for _, err := range err.(validator.ValidationErrors) {
			msg += fmt.Sprintf("Field: '%s', Condition: '%s'\n", err.Field(), err.Tag())
		}
		return nil, httperror.New(err, msg, http.StatusBadRequest)
	}
	err = uds.store.InsertUserBankCard(ctx, &userBankCard)
	if err != nil {
		return nil, err
	}
	return &userBankCard, nil
}

func (uds *UserDataService) UpdateUserTextData(
	ctx context.Context,
	userID uuid.UUID,
	ID uuid.UUID,
	name string,
	text string,
	metadata map[string]interface{},
) (*models.UserTextData, error) {
	userTextData, err := uds.store.GetUserTextData(ctx, ID, userID)
	if err != nil {
		return nil, err
	}
	userTextData.Name = name
	userTextData.Data = text
	userTextData.Metadata = metadata
	userTextData.UpdatedAt = time.Now()
	err = uds.store.UpdateUserTextData(ctx, userTextData)
	if err != nil {
		return nil, err
	}
	return userTextData, nil
}

func (uds *UserDataService) UpdateUserFileData(
	ctx context.Context,
	userID uuid.UUID,
	ID uuid.UUID,
	name string,
	metadata map[string]interface{},
) (*models.UserFileData, error) {
	userFileData, err := uds.store.GetUserFileData(ctx, ID, userID)
	if err != nil {
		return nil, err
	}
	dir, _ := strings.CutSuffix(userFileData.PathToFile, userFileData.Name)
	newPathToFile := fmt.Sprintf("%s/%s", dir, name)
	if err := os.Rename(userFileData.PathToFile, newPathToFile); err != nil {
		return nil, err
	}
	userFileData.Name = name
	userFileData.PathToFile = newPathToFile
	userFileData.Metadata = metadata
	userFileData.UpdatedAt = time.Now()
	err = uds.store.UpdateUserFileData(ctx, userFileData)
	if err != nil {
		return nil, err
	}
	return userFileData, nil
}

func (uds *UserDataService) UpdateUserAuthInfo(
	ctx context.Context,
	userID uuid.UUID,
	ID uuid.UUID,
	name string,
	login, password string,
	metadata map[string]interface{},
) (*models.UserAuthInfo, error) {
	userAuthInfo := models.NewUserAuthInfo(name, userID, metadata, login, password)
	userAuthInfo.ID = ID
	err := uds.store.UpdateUserAuthInfo(ctx, &userAuthInfo)
	if err != nil {
		return nil, err
	}
	return &userAuthInfo, nil
}

func (uds *UserDataService) UpdateUserBankCard(
	ctx context.Context,
	userID uuid.UUID,
	ID uuid.UUID,
	name string,
	number, cardHolder, expireDate, csc string,
	metadata map[string]interface{},
) (*models.UserBankCard, error) {
	userBankCard := models.NewUserBankCard(name, userID, metadata, number, cardHolder, expireDate, csc)
	userBankCard.ID = ID
	err := uds.store.UpdateUserBankCard(ctx, &userBankCard)
	if err != nil {
		return nil, err
	}
	return &userBankCard, nil
}

func (uds *UserDataService) GetUserTextData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserTextData, error) {
	return uds.store.GetUserTextData(ctx, dataID, userID)
}

func (uds *UserDataService) GetUserFileData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserFileData, error) {
	return uds.store.GetUserFileData(ctx, dataID, userID)
}

func (uds *UserDataService) GetUserAuthInfo(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserAuthInfo, error) {
	return uds.store.GetUserAuthInfo(ctx, dataID, userID)
}

func (uds *UserDataService) GetUserBankCard(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserBankCard, error) {
	return uds.store.GetUserBankCard(ctx, dataID, userID)
}

func (uds *UserDataService) DeleteUserTextData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error {
	return uds.store.DeleteUserTextData(ctx, dataID, userID)
}

func (uds *UserDataService) DeleteUserFileData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error {
	userFileData, err := uds.store.GetUserFileData(ctx, dataID, userID)
	if err != nil {
		return err
	}
	if err := os.Remove(userFileData.PathToFile); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return uds.store.DeleteUserFileData(ctx, dataID, userID)
}

func (uds *UserDataService) DeleteUserAuthInfo(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error {
	return uds.store.DeleteUserAuthInfo(ctx, dataID, userID)
}

func (uds *UserDataService) DeleteUserBankCard(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error {
	return uds.store.DeleteUserBankCard(ctx, dataID, userID)
}
