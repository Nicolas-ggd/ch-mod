package models

type Type string

// Defined variables to identify token type
const (
	Auth       Type = "auth"
	Validation Type = "validation"
)

type UserToken struct {
	Hash   []byte `json:"hash"`
	UserID uint   `json:"user_id" binding:"required" gorm:"primaryKey;autoIncrement:false;foreignKey:user_id"`
	Type   Type   `json:"locale" binding:"required" gorm:"primaryKey;autoIncrement:false"`
}

type TokenClaim struct {
	UserId uint    `json:"user_id"`
	Role   *string `json:"role"`
	Time   *string `json:"time"`
}
