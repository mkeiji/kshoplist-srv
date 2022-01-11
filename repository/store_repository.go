package repository

import (
	"database/sql"
	appdb "kshoplistSrv/database"
	"kshoplistSrv/models"
	"log"
	"time"
)

type StoreRepository struct {
	Db *sql.DB
}

func NewStoreRepository() StoreRepository {
	return StoreRepository{Db: appdb.Db}
}

func (this StoreRepository) GetAll() ([]models.Store, error) {
	var stores []models.Store

	res, err := this.Db.Query("select * from Stores")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for res.Next() {
		var (
			id        int
			createdAt time.Time
			updatedAt time.Time
			name      string
		)

		err := res.Scan(&id, &createdAt, &updatedAt, &name)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		stores = append(stores, models.Store{
			Id:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Name:      name,
		})
	}

	return stores, nil
}
