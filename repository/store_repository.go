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

	res, err := this.Db.Query("SELECT * FROM Stores")
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var (
			id        int
			createdAt time.Time
			updatedAt time.Time
			name      string
		)

		err := res.Scan(&id, &createdAt, &updatedAt, &name)
		if err != nil {
			log.Println("Error scanning result:", err)
			return nil, err
		}

		stores = append(stores, models.Store{
			Id:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Name:      name,
		})
	}

	if err = res.Err(); err != nil {
		log.Println("Error during row iteration:", err)
		return nil, err
	}

	return stores, nil
}
