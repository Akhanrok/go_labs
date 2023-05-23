package models

import "time"

type ShoppingList struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId"`
	Name      string    `json:"name"`
	StoreName string    `json:"storeName"`
	Quantity  string    `json:"quantity"`
	Link      string    `json:"link"`
	CreatedAt time.Time `json:"createdAt"`
}

func GetShoppingListByID(id uint) (*ShoppingList, error) {
	return nil, nil
}

func SaveShoppingList(shoppingList *ShoppingList) error {
	return nil
}
