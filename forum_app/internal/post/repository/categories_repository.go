package repository

import (
	"context"
	"database/sql"
	"forum_app/internal/entity"
	"log"
)

type CategoriesRepository struct {
	db       *sql.DB
	errorLog *log.Logger
}

func NewCategoriesRepository(db *sql.DB, errorLog *log.Logger) *CategoriesRepository {
	return &CategoriesRepository{db, errorLog}
}

func (cr *CategoriesRepository) FetchById(ctx context.Context, id int) (entity.Category, error) {
	category := entity.Category{}
	tx, err := cr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		cr.errorLog.Println(err)
		return entity.Category{}, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM categories WHERE id = ?;")
	if err != nil {
		cr.errorLog.Println(err)
		return entity.Category{}, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		cr.errorLog.Println(err)
		return entity.Category{}, err
	}
	if rows.Next() {
		rows.Scan(&category.Id, &category.Title)
	}
	if err = tx.Commit(); err != nil {
		cr.errorLog.Println(err)
		return entity.Category{}, err
	}
	return category, nil
}

func (cr *CategoriesRepository) FetchAllCategories(ctx context.Context) ([]entity.Category, error) {
	categories := []entity.Category{}
	tx, err := cr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		cr.errorLog.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM categories")
	if err != nil {
		cr.errorLog.Println(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		cr.errorLog.Println(err)
		return nil, err
	}
	for rows.Next() {
		category := entity.Category{}
		rows.Scan(&category.Id, &category.Title)
		categories = append(categories, category)
	}
	if err = tx.Commit(); err != nil {
		cr.errorLog.Println(err)
		return nil, err
	}
	return categories, nil
}

func (cr *CategoriesRepository) FetchByPostId(ctx context.Context, id int) ([]entity.Category, error) {
	categories := []entity.Category{}
	tx, err := cr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		cr.errorLog.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT c.id, c.title FROM categories as c INNER JOIN post_categories as pc ON c.id = pc.category_id WHERE pc.post_id = ?")
	if err != nil {
		cr.errorLog.Println(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		cr.errorLog.Println(err)
		return nil, err
	}
	for rows.Next() {
		category := entity.Category{}
		rows.Scan(&category.Id, &category.Title)
		categories = append(categories, category)
	}
	if err = tx.Commit(); err != nil {
		cr.errorLog.Println(err)
		return nil, err
	}
	return categories, nil
}
