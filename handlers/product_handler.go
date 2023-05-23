package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Akhanrok/go_labs/database"
	"github.com/Akhanrok/go_labs/models"
)

type ProductHandler struct {
	db *database.Database
}

func NewProductHandler(db *database.Database) *ProductHandler {
	return &ProductHandler{db: db}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.db.CreateProduct(&product)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	productID := parseUintPathParam(r, "productID")

	product, err := h.db.GetProduct(productID)
	if err != nil {
		http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID := parseUintPathParam(r, "productID")

	existingProduct, err := h.db.GetProduct(productID)
	if err != nil {
		http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
		return
	}

	var updatedProduct models.Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingProduct.Name = updatedProduct.Name
	existingProduct.Type = updatedProduct.Type

	err = h.db.UpdateProduct(existingProduct)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := parseUintPathParam(r, "productID")

	err := h.db.DeleteProduct(productID)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
