package persistence

import (
	"github.com/ExchangeDiary_Server/exchange-diary/domain/entity"
	"github.com/ExchangeDiary_Server/exchange-diary/domain/repository"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type RoomModel struct {
	BaseModel
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


func (roomRepository *RoomRepository) Create(room *entity.Room) (*entity.Room, error) {
	roomModel := RoomModel{}
	copier.Copy(&roomModel, &room)
	if err := roomRepository.db.Create(&roomModel).Error; err != nil {
		return nil, err
	}
	newRoom := new(entity.Room)
	copier.Copy(&newRoom, &roomModel)
	return newRoom, nil
}

func (roomRepository *RoomRepository) GetByID(id int) (*entity.Room, error) {
	return &entity.Room{}, nil
}

func (roomRepository *RoomRepository) GetAll() (*entity.Rooms, error) {
	return &entity.Rooms{}, nil
}

func (roomRepository *RoomRepository) GetAllByAccountId(accountId int) (*entity.Rooms, error) {
	return &entity.Rooms{}, nil
}

func (roomRepository *RoomRepository) Update(room *entity.Room) (*entity.Room, error) {
	return &entity.Room{}, nil
}

func (roomRepository *RoomRepository) Delete(room *entity.Room) error {
	return nil
}