package repository

import (
	"database/sql"
	"forum/internal/entity"
)

type CategoriesRepository struct {
	db *sql.DB
}

func NewCategoriesRepository(db *sql.DB) *CategoriesRepository {
	return &CategoriesRepository{db}
}

func FetchByPostId(id int) ([]entity.Category, error) {
	return nil, nil
}

func FetchById(id int) (entity.Category, error) {
	return entity.Category{}, nil
}
