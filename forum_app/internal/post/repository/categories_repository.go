package repository

import (
	"context"
	"database/sql"
	"errors"
	"forum_app/internal/entity"
)

type CategoriesRepository struct {
	db *sql.DB
}

func NewCategoriesRepository(db *sql.DB) *CategoriesRepository {
	return &CategoriesRepository{db}
}

func (cr *CategoriesRepository) FetchAllCategories(ctx context.Context) ([]entity.Category, error) {
	categories := []entity.Category{}
	tx, err := cr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		category := entity.Category{}
		rows.Scan(&category.Id, &category.Title)
		categories = append(categories, category)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (cr *CategoriesRepository) FetchByPostId(ctx context.Context, id int) ([]entity.Category, error) {
	categories := []entity.Category{}
	tx, err := cr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT c.id, c.title FROM categories as c INNER JOIN post_categories as pc ON c.id = pc.category_id WHERE pc.post_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		category := entity.Category{}
		rows.Scan(&category.Id, &category.Title)
		categories = append(categories, category)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (cr *CategoriesRepository) Store(ctx context.Context, post entity.Post) error {
	tx, err := cr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO post_categories(post_id, category_id) VALUES(?, ?);")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, category := range post.Category {
		res, err := stmt.ExecContext(ctx, post.Id, category.Id)
		if err != nil {
			return err
		}
		if rAffected, err := res.RowsAffected(); rAffected > 1 || err != nil {
			if err != nil {
				return err
			}
			return errors.New(`more than one row has been affected`)
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
