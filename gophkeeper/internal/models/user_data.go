package models

import (
	"time"

	"github.com/google/uuid"
)

// Общие данные для всех записей
// Метаданные для сущностей
type Metadata map[string]interface{}

// Базовые данные для всех пользовательских данных
type BaseUserData struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	Metadata
}

func NewBaseUserData(name string, userID uuid.UUID, metadata Metadata) BaseUserData {
	return BaseUserData{
		ID:       uuid.New(),
		UserID:   userID,
		Name:     name,
		Metadata: metadata,
	}
}

// Аутентификационные данные
type UserAuthInfo struct {
	BaseUserData
	Login    string `db:"login" json:"login" validate:"required"`
	Password string `db:"password" json:"password" validate:"required"`
}

func NewUserAuthInfo(name string, userID uuid.UUID, metadata Metadata, login, password string) UserAuthInfo {
	return UserAuthInfo{
		BaseUserData: NewBaseUserData(name, userID, metadata),
		Login:        login,
		Password:     password,
	}
}

// Текстовые данные
type UserTextData struct {
	BaseUserData
	Data string `db:"data" json:"data"`
}

func NewUserTextData(name string, userID uuid.UUID, metadata Metadata, text string) UserTextData {
	return UserTextData{
		BaseUserData: NewBaseUserData(name, userID, metadata),
		Data:         text,
	}
}

// Бинарные данные
type UserFileData struct {
	BaseUserData
	PathToFile string `db:"path_to_file" json:"path_to_file"`
}

func NewUserFileData(name string, userID uuid.UUID, metadata Metadata, pathToFile string) UserFileData {
	return UserFileData{
		BaseUserData: NewBaseUserData(name, userID, metadata),
		PathToFile:   pathToFile,
	}
}

// Банковская карта
type UserBankCard struct {
	BaseUserData
	Number     string `db:"number" json:"number" validate:"required,creditcard"`
	CardHolder string `db:"card_holder" json:"card_holder"`
	ExpireDate string `db:"expire_date" json:"expire_date" validate:"required,datetime=01/2006"`
	CSC        string `db:"csc" json:"csc" validate:"len=3"`
}

func NewUserBankCard(
	name string,
	userID uuid.UUID,
	metadata Metadata,
	number, cardHolder, expireDate, csc string,
) UserBankCard {
	return UserBankCard{
		BaseUserData: NewBaseUserData(name, userID, metadata),
		Number:       number,
		CardHolder:   cardHolder,
		ExpireDate:   expireDate,
		CSC:          csc,
	}
}
