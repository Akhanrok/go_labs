package models

import "time"

type Product struct {
	ID             uint      `json:"id"`
	ShoppingListID uint      `json:"shoppingListId"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	CreatedAt      time.Time `json:"createdAt"`
}

func GetProductByID(id uint) (*Product, error) {
	return nil, nil
}

func SaveProduct(product *Product) error {
	return nil
}
