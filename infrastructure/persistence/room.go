package persistence

import (
	"github.com/exchange-diary/domain/entity"
	"github.com/exchange-diary/domain/repository"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type RoomModel struct {
	ID int `gorm:"primary_key"`
	BaseDateModel
	entity.Room
}

func (_ *RoomModel) TableName() string {
	return "rooms"
}

type RoomModels []RoomModel

type RoomRepository struct {
	db *gorm.DB
}

// Domain layer의 RoomRepository interface를 만족시키는 repository impl.
// gorm connection을 들고 가지고 있다.
func NewRoomRepository(db *gorm.DB) repository.RoomRepository {
	return &RoomRepository{db: db}
}

func (rr *RoomRepository) Create(room *entity.Room) (*entity.Room, error) {
	roomModel := RoomModel{}
	copier.Copy(&roomModel, &room)
	if err := rr.db.Create(&roomModel).Error; err != nil {
		return nil, err
	}
	newRoom := new(entity.Room)
	copier.Copy(&newRoom, &roomModel)
	return newRoom, nil
}

func (rr *RoomRepository) GetByID(id int) (*entity.Room, error) {
	roomModel := RoomModel{ID: id}
	if err := rr.db.First(&roomModel).Error; err != nil {
		return nil, err
	}
	room := new(entity.Room)
	copier.Copy(&room, &roomModel)
	return room, nil
}

func (rr *RoomRepository) GetAll(limit, offset int) (*entity.Rooms, error) {
	roomModels := RoomModels{}
	rr.db.Limit(limit).Offset(offset).Find(&roomModels)
	rooms := new(entity.Rooms)
	copier.Copy(&rooms, &roomModels)
	return rooms, nil
}

func (rr *RoomRepository) GetAllByAccountId(accountId, limit, offset int) (*entity.Rooms, error) {
	return &entity.Rooms{}, nil
}

func (rr *RoomRepository) Update(room *entity.Room) (*entity.Room, error) {
	return &entity.Room{}, nil
}

func (rr *RoomRepository) Delete(room *entity.Room) error {
	roomModel := RoomModel{}
	copier.Copy(&roomModel, &room)
	if err := rr.db.Delete(&roomModel).Error; err != nil {
		return err
	}
	return nil
}
