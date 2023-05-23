package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Akhanrok/go_labs/database"
	"github.com/Akhanrok/go_labs/models"
)

type ShoppingListHandler struct {
	db *database.Database
}

func NewShoppingListHandler(db *database.Database) *ShoppingListHandler {
	return &ShoppingListHandler{db: db}
}

func (h *ShoppingListHandler) CreateShoppingList(w http.ResponseWriter, r *http.Request) {
	var list models.ShoppingList
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	list.CreatedAt = time.Now()

	err = h.db.CreateShoppingList(&list)
	if err != nil {
		http.Error(w, "Failed to create shopping list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ShoppingListHandler) AddProductToList(w http.ResponseWriter, r *http.Request) {
	listID := parseUintPathParam(r, "listID")
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product.ShoppingListID = listID

	err = h.db.AddProductToList(&product)
	if err != nil {
		http.Error(w, "Failed to add product to shopping list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ShoppingListHandler) RemoveProductFromList(w http.ResponseWriter, r *http.Request) {
	listID := parseUintPathParam(r, "listID")
	productID := parseUintPathParam(r, "productID")

	err := h.db.RemoveProductFromList(listID, productID)
	if err != nil {
		http.Error(w, "Failed to remove product from shopping list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ShoppingListHandler) SortShoppingListByType(w http.ResponseWriter, r *http.Request) {
	listID := parseUintPathParam(r, "listID")

	err := h.db.SortShoppingListByType(listID)
	if err != nil {
		http.Error(w, "Failed to sort shopping list by product type", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func parseUintPathParam(r *http.Request, param string) uint {
	value := r.URL.Query().Get(param)
	u, _ := strconv.ParseUint(value, 10, 64)
	return uint(u)
}
