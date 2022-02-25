package persistence

import (
	"github.com/exchange-diary/domain/entity"
	"github.com/exchange-diary/domain/repository"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// RoomGorm is a db representation of entity.Room
type RoomGorm struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"column:name;not null"`
	Code   string `gorm:"column:code;not null"`
	Hint   string `gorm:"column:hint;not null"`
	Theme  string `gorm:"column:theme;not null"`
	Period uint8  `gorm:"column:period;not null"`

	MasterID      uint             `gorm:"column:master_id"`
	Master        AccountGorm `gorm:"foreignKey:MasterID"`
	TurnAccountID uint             `gorm:"column:turn_account_id"`
	TurnAccount   AccountGorm `gorm:"foreignKey:TurnAccountID"`

	// TODO: json field
	// Orders        []uint

	BaseGormModel
}

// TableName define gorm table name
func (RoomGorm) TableName() string {
	return "rooms"
}

// RoomGorms define list of RoomGorm
type RoomGorms []RoomGorm

// RoomRepository is a impl of domain/repository/roomRepository.go RoomRepository interface
type RoomRepository struct {
	db *gorm.DB
}

// NewRoomRepository ...
// Domain layer의 RoomRepository interface를 만족시키는 repository impl.
// gorm connection을 들고 가지고 있다.
func NewRoomRepository(db *gorm.DB) repository.RoomRepository {
	return &RoomRepository{db: db}
}

// Create func inserts a row to db
func (rr *RoomRepository) Create(room *entity.Room) (*entity.Room, error) {
	dto := RoomGorm{}
	copier.Copy(&dto, &room)
	if err := rr.db.Create(&dto).Error; err != nil {
		return nil, err
	}
	newRoom := new(entity.Room)
	copier.Copy(&newRoom, &dto)
	return newRoom, nil
}

// GetByID func find a row by entity's ID from db
func (rr *RoomRepository) GetByID(id uint) (*entity.Room, error) {
	dto := RoomGorm{ID: id}
	if err := rr.db.First(&dto).Error; err != nil {
		return nil, err
	}
	room := new(entity.Room)
	copier.Copy(&room, &dto)
	return room, nil
}

// GetAll func get all row from db table
func (rr *RoomRepository) GetAll(limit, offset uint) (*entity.Rooms, error) {
	dto := RoomGorms{}
	rr.db.Limit(int(limit)).Offset(int(offset)).Find(&dto)
	rooms := new(entity.Rooms)
	copier.Copy(&rooms, &dto)
	return rooms, nil
}

// GetAllByAccountID func finds all Rooms which account is joined from db table
func (rr *RoomRepository) GetAllByAccountID(accountID, limit, offset uint) (*entity.Rooms, error) {
	return &entity.Rooms{}, nil
}

// Update func update a room fields
func (rr *RoomRepository) Update(room *entity.Room) (*entity.Room, error) {
	return &entity.Room{}, nil
}

// Delete func delete a room
func (rr *RoomRepository) Delete(room *entity.Room) error {
	dto := RoomGorm{}
	copier.Copy(&dto, &room)
	if err := rr.db.Delete(&dto).Error; err != nil {
		return err
	}
	return nil
}
