package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eac0de/gophkeeper/shared/pkg/httperror"
	gpvalidator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Общие данные для всех записей
// Метаданные для сущностей
type Metadata map[string]interface{}

var validator *gpvalidator.Validate

func init() {
	validator = gpvalidator.New()

	// Регистрация кастомного валидатора
	validator.RegisterValidation("card_number", func(fl gpvalidator.FieldLevel) bool {
		// Простая проверка: длина карты должна быть от 13 до 19 символов
		cardNumber := fl.Field().String()
		if len(cardNumber) < 13 || len(cardNumber) > 19 {
			return false
		}

		// Реализация алгоритма Луна
		return luhnCheck(cardNumber)
	})
}

// luhnCheck реализует алгоритм Луна
func luhnCheck(cardNumber string) bool {
	var sum int
	alt := false
	for i := len(cardNumber) - 1; i >= 0; i-- {
		n := int(cardNumber[i] - '0')
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}

// Базовые данные для всех пользовательских данных
type BaseUserData struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"-"`
	Name      string    `db:"name" json:"name" validate:"required"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	Metadata Metadata `db:"metadata" json:"metadata"`
}

func NewBaseUserData(name string, userID uuid.UUID, metadata Metadata) BaseUserData {
	return BaseUserData{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  metadata,
	}
}

// Аутентификационные данные
type UserAuthInfo struct {
	BaseUserData
	Login    string `db:"login" json:"login" validate:"required"`
	Password string `db:"password" json:"password" validate:"required"`
}

func NewUserAuthInfo(name string, userID uuid.UUID, metadata Metadata, login, password string) (UserAuthInfo, error) {
	userAuthInfo := UserAuthInfo{
		BaseUserData: NewBaseUserData(name, userID, metadata),
		Login:        login,
		Password:     password,
	}
	return userAuthInfo, Validate(userAuthInfo)
}

// Текстовые данные
type UserTextData struct {
	BaseUserData
	Data string `db:"data" json:"data" validate:"required"`
}

func NewUserTextData(name string, userID uuid.UUID, metadata Metadata, text string) (UserTextData, error) {
	userTxtData := UserTextData{
		BaseUserData: NewBaseUserData(name, userID, metadata),
		Data:         text,
	}
	return userTxtData, Validate(userTxtData)
}

// Бинарные данные
type UserFileData struct {
	BaseUserData
	PathToFile string `db:"path_to_file" json:"-"`
	Ext        string `db:"ext" json:"ext"`
}

func NewUserFileData(name string, userID uuid.UUID, pathToFile, ext string) (UserFileData, error) {
	userFileData := UserFileData{
		BaseUserData: NewBaseUserData(name, userID, Metadata{}),
		PathToFile:   pathToFile,
		Ext:          ext,
	}
	return userFileData, Validate(userFileData)
}

// Банковская карта
type UserBankCard struct {
	BaseUserData
	Number     string `db:"number" json:"number" validate:"required,card_number"`
	CardHolder string `db:"card_holder" json:"card_holder" validate:"required"`
	ExpireDate string `db:"expire_date" json:"expire_date" validate:"required,datetime=01/06"`
	CSC        string `db:"csc" json:"csc" validate:"required,len=3"`
}

func NewUserBankCard(
	name string,
	userID uuid.UUID,
	metadata Metadata,
	number, cardHolder, expireDate, csc string,
) (UserBankCard, error) {
	userBankCard := UserBankCard{
		BaseUserData: NewBaseUserData(name, userID, metadata),
		Number:       number,
		CardHolder:   cardHolder,
		ExpireDate:   expireDate,
		CSC:          csc,
	}
	return userBankCard, Validate(userBankCard)
}

func Validate(item interface{}) error {
	err := validator.Struct(item)
	if err != nil {
		var msg string
		for _, err := range err.(gpvalidator.ValidationErrors) {
			msg += fmt.Sprintf("Field: '%s', Condition: '%s'\n", err.Field(), err.Tag())
		}
		return httperror.New(err, msg, http.StatusUnprocessableEntity)
	}
	return nil
}
