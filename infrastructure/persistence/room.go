package persistence

import (
	"encoding/json"
	"time"

	"github.com/ExchangeDiary/exchange-diary/domain/entity"
	"github.com/ExchangeDiary/exchange-diary/domain/repository"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/logger"
	"github.com/jinzhu/copier"
	"gorm.io/datatypes"
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

	MasterID      uint       `gorm:"column:master_id"`
	Master        MemberGorm `gorm:"column:master_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TurnAccountID uint       `gorm:"column:turn_account_id"`
	DueAt         time.Time  `gorm:"column:due_at"`
	TurnAccount   MemberGorm `gorm:"column:turn_account_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Orders datatypes.JSON `gorm:"column:orders"`

	BaseGormModel
}

// RoomGorms define list of RoomGorm
type RoomGorms []RoomGorm

// RoomRepository is a impl of domain/repository/roomRepository.go RoomRepository interface
type RoomRepository struct {
	db *gorm.DB
}

// TableName define gorm table name
func (RoomGorm) TableName() string {
	return "rooms"
}

// NewRoomRepository ...
// Domain layer의 RoomRepository interface를 만족시키는 repository impl.
// gorm connection을 들고 가지고 있다.
func NewRoomRepository(db *gorm.DB) repository.RoomRepository {
	return &RoomRepository{db: db}
}

func masterOrMember(masterID uint, memberRoomIDs []uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if memberRoomIDs != nil {
			return db.Where("id IN (?)", memberRoomIDs).Or("master_id = ?", masterID)
		}
		return db.Where("master_id = ?", masterID)
	}
}

func paginate(limit, offset uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset)).Limit(int(limit))
	}
}

// ToDTO : entity.Room -> RoomGorm
func ToDTO(dto *RoomGorm, room *entity.Room) *RoomGorm {
	copier.Copy(&dto, &room)
	ordersJSON, _ := room.OrdersToJSON()
	dto.Orders = datatypes.JSON(ordersJSON)
	return dto
}

// ToEntity : RoomGorm -> entity.Room
func ToEntity(dto *RoomGorm) *entity.Room {
	var orders []uint

	room := new(entity.Room)
	copier.Copy(&room, &dto)

	err := json.Unmarshal([]byte(dto.Orders), &orders)
	if err != nil {
		logger.Error(err.Error())
	}
	room.Orders = orders
	return room
}

// Create func inserts a row to db
func (rr *RoomRepository) Create(room *entity.Room) (*entity.Room, error) {
	dto := ToDTO(&RoomGorm{}, room)
	if err := rr.db.Create(&dto).Error; err != nil {
		return nil, err
	}
	return ToEntity(dto), nil
}

// GetByID func find a row by entity's ID from db
func (rr *RoomRepository) GetByID(id uint) (*entity.Room, error) {
	dto := RoomGorm{ID: id}
	if err := rr.db.First(&dto).Error; err != nil {
		return nil, err
	}
	return ToEntity(&dto), nil
}

// GetAll rooms from masterID or roomIDs
func (rr *RoomRepository) GetAll(accountID uint, roomIDs []uint, limit, offset uint) (*entity.Rooms, error) {
	dto := RoomGorms{}
	rr.db.Scopes(paginate(limit, offset), masterOrMember(accountID, roomIDs)).Order(" updated_at desc ").Find(&dto)
	rooms := entity.Rooms{}
	for _, roomGorm := range dto {
		rooms = append(rooms, *ToEntity(&roomGorm))
	}
	return &rooms, nil
}

// Update func update a room fields
func (rr *RoomRepository) Update(room *entity.Room) (*entity.Room, error) {
	dto := ToDTO(&RoomGorm{}, room)
	if err := rr.db.Save(&dto).Error; err != nil {
		return nil, err
	}
	return ToEntity(dto), nil
}

// Delete func delete a room
func (rr *RoomRepository) Delete(room *entity.Room) error {
	dto := ToDTO(&RoomGorm{}, room)
	if err := rr.db.Delete(&dto).Error; err != nil {
		return err
	}
	return nil
}
