package repository

import (
	"database/sql"
	appdb "kshoplistSrv/database"
	"kshoplistSrv/models"
	"log"
	"time"
)

type ItemRepository struct {
	Db *sql.DB
}

func NewItemRepository() ItemRepository {
	return ItemRepository{Db: appdb.Db}
}

func (this ItemRepository) GetAll() []models.Item {
	var items []models.Item

	res, err := this.Db.Query("select * from Items")
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		var (
			id        int
			createdAt time.Time
			updatedAt time.Time
			storeId   int
			name      string
		)

		err := res.Scan(&id, &createdAt, &updatedAt, &storeId, &name)
		if err != nil {
			log.Fatal(err)
		}

		items = append(items, models.Item{
			Id:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			StoreId:   storeId,
			Name:      name,
		})
	}

	return items
}

func (this ItemRepository) Post(item models.Item) {
	stmt, err := this.Db.Prepare("INSERT INTO Items (StoreID, Name) VALUES (?, ?)")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(item.StoreId, item.Name)
	if err != nil {
		log.Fatal(err)
	}

	newItemId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("[ItemRepository]: Inserted item id:", newItemId)
	}
}

func (this ItemRepository) Put(item models.Item) {
	query := "UPDATE Items SET Name=? WHERE ID=?"
	stmt, err := this.Db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(item.Name, item.Id)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("[ItemRepository]: Updated item id:", item.Id)
	}
}

func (this ItemRepository) Delete(item models.Item) {
	stmt, err := this.Db.Prepare("DELETE FROM Items WHERE (ID = ?)")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(item.Id)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("[ItemRepository]: Deleted item id:", item.Id)
	}
}
