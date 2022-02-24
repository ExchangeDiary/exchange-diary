package persistence

// AccountGormModel is a db representation of entity.Account
type AccountGormModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"column:name;not null"`
	BaseGormModel
}

// TableName define gorm table name
func (AccountGormModel) TableName() string {
	return "accounts"
}

// AccountGormModels is a type that represents list of AccountGormModel
type AccountGormModels []AccountGormModel
