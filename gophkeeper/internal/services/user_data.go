package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/eac0de/gophkeeper/internal/models"

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

	GetUserTextDataList(ctx context.Context, userID uuid.UUID, offset int) ([]models.UserTextData, error)
	GetUserFileDataList(ctx context.Context, userID uuid.UUID, offset int) ([]models.UserFileData, error)
	GetUserAuthInfoList(ctx context.Context, userID uuid.UUID, offset int) ([]models.UserAuthInfo, error)
	GetUserBankCardList(ctx context.Context, userID uuid.UUID, offset int) ([]models.UserBankCard, error)

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
	userTextData, err := models.NewUserTextData(name, userID, metadata, text)
	if err != nil {

	}
	err = uds.store.InsertUserTextData(ctx, &userTextData)
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
	ext string,
) (*models.UserFileData, error) {
	userFileData, err := models.NewUserFileData(name, userID, pathToFile, ext)
	if err != nil {

	}
	err = uds.store.InsertUserFileData(ctx, &userFileData)
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
	userAuthInfo, err := models.NewUserAuthInfo(name, userID, metadata, login, password)
	if err != nil {
		return nil, err
	}
	err = uds.store.InsertUserAuthInfo(ctx, &userAuthInfo)
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
	userBankCard, err := models.NewUserBankCard(name, userID, metadata, number, cardHolder, expireDate, csc)
	if err != nil {
		return nil, err
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
	err = models.Validate(userTextData)
	if err != nil {
		return nil, err
	}
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
	if name != userFileData.Name {
		dir := fmt.Sprintf("../user_files/%s", userID.String())
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
		clearName := name
		pathToFile := fmt.Sprintf("%s/%s%s", dir, name, userFileData.Ext)
		count := 0
		for {
			_, err := os.Stat(pathToFile)
			if err != nil {
				if os.IsNotExist(err) {
					break
				} else {
					return nil, err
				}
			}
			count++
			name = fmt.Sprintf("%s(%d)", clearName, count)
			pathToFile = fmt.Sprintf("%s/%s%s", dir, name, userFileData.Ext)
			if name == userFileData.Name {
				break
			}
		}
		if name != userFileData.Name {
			if err := os.Rename(userFileData.PathToFile, pathToFile); err != nil {
				return nil, err
			}
			userFileData.Name = name
			userFileData.PathToFile = pathToFile
		}
	}
	userFileData.Metadata = metadata
	userFileData.UpdatedAt = time.Now()
	err = models.Validate(userFileData)
	if err != nil {
		return nil, err
	}
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
	userAuthInfo, err := uds.store.GetUserAuthInfo(ctx, ID, userID)
	if err != nil {
		return nil, err
	}
	userAuthInfo.Name = name
	userAuthInfo.Login = login
	userAuthInfo.Password = password
	userAuthInfo.Metadata = metadata
	userAuthInfo.UpdatedAt = time.Now()
	err = models.Validate(userAuthInfo)
	if err != nil {
		return nil, err
	}
	err = uds.store.UpdateUserAuthInfo(ctx, userAuthInfo)
	if err != nil {
		return nil, err
	}
	return userAuthInfo, nil
}

func (uds *UserDataService) UpdateUserBankCard(
	ctx context.Context,
	userID uuid.UUID,
	ID uuid.UUID,
	name string,
	number, cardHolder, expireDate, csc string,
	metadata map[string]interface{},
) (*models.UserBankCard, error) {
	userBankCard, err := uds.store.GetUserBankCard(ctx, ID, userID)
	if err != nil {
		return nil, err
	}
	userBankCard.Name = name
	userBankCard.Number = number
	userBankCard.CardHolder = cardHolder
	userBankCard.ExpireDate = expireDate
	userBankCard.CSC = csc
	userBankCard.Metadata = metadata
	userBankCard.UpdatedAt = time.Now()
	err = models.Validate(userBankCard)
	if err != nil {
		return nil, err
	}
	err = uds.store.UpdateUserBankCard(ctx, userBankCard)
	if err != nil {
		return nil, err
	}
	return userBankCard, nil
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

func (uds *UserDataService) GetUserTextDataList(ctx context.Context, userID uuid.UUID, offset int) ([]models.UserTextData, error) {
	return uds.store.GetUserTextDataList(ctx, userID, offset)
}

func (uds *UserDataService) GetUserFileDataList(ctx context.Context, userID uuid.UUID, offset int) ([]models.UserFileData, error) {
	return uds.store.GetUserFileDataList(ctx, userID, offset)
}

func (uds *UserDataService) GetUserAuthInfoList(ctx context.Context, userID uuid.UUID, offset int) ([]models.UserAuthInfo, error) {
	return uds.store.GetUserAuthInfoList(ctx, userID, offset)
}

func (uds *UserDataService) GetUserBankCardList(ctx context.Context, userID uuid.UUID, offset int) ([]models.UserBankCard, error) {
	return uds.store.GetUserBankCardList(ctx, userID, offset)
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
