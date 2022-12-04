package repository

import (
	"context"
	"database/sql"
	"forum_app/internal/entity"
)

type PostsRepository struct {
	db *sql.DB
}

func NewPostsRepository(db *sql.DB) *PostsRepository {
	return &PostsRepository{db}
}

func (pr *PostsRepository) FetchById(ctx context.Context, id int) (entity.Post, error) {
	post := entity.Post{}
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return post, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM posts WHERE id = ?;")
	if err != nil {
		return post, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return post, err
	}
	if rows.Next() {
		rows.Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content)
	}
	if err = tx.Commit(); err != nil {
		return entity.Post{}, err
	}
	return post, nil
}

func (pr *PostsRepository) FetchAll(ctx context.Context) ([]entity.Post, error) {
	posts := []entity.Post{}
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM posts;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		post := entity.Post{}
		rows.Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content)
		posts = append(posts, post)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (pr *PostsRepository) FetchByUserId(ctx context.Context, id int) ([]entity.Post, error) {
	posts := []entity.Post{}
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM posts WHERE user_id = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		post := entity.Post{}
		rows.Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content)
		posts = append(posts, post)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (pr *PostsRepository) FetchByCategoryId(ctx context.Context, id int) ([]entity.Post, error) {
	posts := []entity.Post{}
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT post_id FROM post_categories WHERE category_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		post := entity.Post{}
		rows.Scan(&post.Id)
		posts = append(posts, post)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (pr *PostsRepository) Store(ctx context.Context, post entity.Post) (int64, error) {
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO posts(user_id, date, title, content) VALUES(?,?,?,?);`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, post.User.Id, post.Date, post.Title, post.Content)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

//for future use
// func (pr *PostsRepository) Update(entity.Post) error {
// 	return nil
// }

// func (pr *PostsRepository) Delete(id int) error {
// 	result, err := pr.db.Exec(DELETE_QUERY+POSTS+BY_ID, id)
// 	if err != nil {
// 		return err
// 	}
// 	nbr, err := result.RowsAffected()
// 	if nbr > 1 {
// 		return errors.New("more than one row has been affected")
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
