package schemes

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
	Name      string    `db:"name" json:"name" validate:"required"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	Metadata Metadata `db:"metadata" json:"metadata"`
}

func NewBaseUserData(name string, userID uuid.UUID, metadata Metadata) BaseUserData {
	return BaseUserData{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  metadata,
	}
}

type UserAuthInfo struct {
	BaseUserData
	Login    string `db:"login" json:"login" validate:"required"`
	Password string `db:"password" json:"password" validate:"required"`
}

// Текстовые данные
type UserTextData struct {
	BaseUserData
	Data string `db:"data" json:"data" validate:"required"`
}

// Бинарные данные
type UserFileData struct {
	BaseUserData
	PathToFile string `db:"path_to_file" json:"-"`
}

// Банковская карта
type UserBankCard struct {
	BaseUserData
	Number     string `db:"number" json:"number" validate:"required,credit_card"`
	CardHolder string `db:"card_holder" json:"card_holder"`
	ExpireDate string `db:"expire_date" json:"expire_date" validate:"required,datetime=01/06"`
	CSC        string `db:"csc" json:"csc" validate:"len=3"`
}
