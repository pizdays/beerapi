package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/exercise/beer/model"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var ErrBeerNotFound = errors.New("Beer not found")

type BeerService struct {
	mongoDB *mongo.Collection
	mariaDB *gorm.DB
}

func NewBeerService(mongoDB *mongo.Collection, mariaDB *gorm.DB) BeerService {
	return BeerService{mongoDB: mongoDB, mariaDB: mariaDB}

}

func (b BeerService) Find(beerId uint) (model.Beer, error) {

	var user model.Beer
	if err := b.mariaDB.First(&user, "id = ?", beerId).Error; err != nil {
		return model.Beer{}, fmt.Errorf("service.User.Find: %w", ErrBeerNotFound)
	}

	return user, nil

}

func (b BeerService) FindAll(search string, limit, offset int) ([]model.Beer, int64, error) {

	var (
		beer  []model.Beer
		count int64
	)

	if err := b.mariaDB.Model(&beer).Count(&count).Error; err != nil {
		return []model.Beer{}, 0, fmt.Errorf("service.Beer.FindAll: %w", err)
	}

	if search == "" {
		if err := b.mariaDB.Model(&beer).Limit(limit).Offset(offset).Scan(&beer).Error; err != nil {
			return []model.Beer{}, 0, fmt.Errorf("service.Beer.FindAll: %w", err)
		}

	} else {

		if err := b.mariaDB.Model(&beer).Where("name LIKE ?", search+"%").Limit(limit).Offset(offset).Find(&beer).Error; err != nil {
			return []model.Beer{}, 0, fmt.Errorf("service.Beer.FindAll: %w", err)
		}

	}

	return beer, count, nil

}

func (b BeerService) Create(beer model.Beer) (model.Beer, error) {

	err := b.mariaDB.Create(&beer).Error
	if err != nil {
		return model.Beer{}, fmt.Errorf("service.Beer.Create: %w", err)
	}

	return beer, nil

}

func (b BeerService) CreateLog(logs model.Logs) (*mongo.InsertOneResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logResult, err := b.mongoDB.InsertOne(ctx, logs)

	if err != nil {
		return nil, fmt.Errorf("service.Beer.CreateLog: %w", err)
	}

	return logResult, nil

}

func (b BeerService) Update(beer model.Beer) (model.Beer, error) {

	err := b.mariaDB.
		Where("id = ?", beer.ID).
		Save(&beer).Error
	if err != nil {
		return model.Beer{}, fmt.Errorf("service.Beer.Update: %w", err)
	}

	return beer, nil

}

func (b BeerService) Delete(beerId uint) error {

	err := b.mariaDB.Delete(&model.Beer{}, beerId).
		Error
	if err != nil {
		return fmt.Errorf("service.Beer.Delete: %w", err)
	}

	return nil

}

func (b BeerService) CreateLogOnMongoDB(log model.Logs) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := b.mongoDB.InsertOne(ctx, log)
	id := res.InsertedID

	if err != nil {
		return nil, fmt.Errorf("service.Beer.CreateLogOnMongoDB: %w", err)

	}
	return id, nil

}
