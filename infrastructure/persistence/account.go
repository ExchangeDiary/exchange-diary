package persistence

// AccountGorm is a db representation of entity.Account
type AccountGorm struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"column:name;not null"`
	BaseGormModel
}

// TableName define gorm table name
func (AccountGorm) TableName() string {
	return "accounts"
}

// AccountGorms is a type that represents list of AccountGorm
type AccountGorms []AccountGorm
