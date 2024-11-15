package database

import (
	"database/sql"
	"fmt"

	"github.com/Yashupadhyaya/inventory-service/models"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDatabase(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}
	return db.Ping()
}

func GetInventoryByProductID(productID string) (*models.Inventory, error) {
	inventory := &models.Inventory{}
	err := db.QueryRow("SELECT product_id, quantity FROM inventory WHERE product_id = $1", productID).
		Scan(&inventory.ProductID, &inventory.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}
	return inventory, nil
}

func UpdateInventory(productID string, quantity int) error {
	_, err := db.Exec("UPDATE inventory SET quantity = $1 WHERE product_id = $2", quantity, productID)
	return err
}
