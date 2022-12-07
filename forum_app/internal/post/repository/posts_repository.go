package repository

import (
	"context"
	"database/sql"
	"forum_app/internal/entity"
	"log"
	"time"
)

type PostsRepository struct {
	db       *sql.DB
	errorLog *log.Logger
}

func NewPostsRepository(db *sql.DB, errorLog *log.Logger) *PostsRepository {
	return &PostsRepository{db, errorLog}
}

func (pr *PostsRepository) FetchById(ctx context.Context, id int) (entity.Post, error) {
	post := entity.Post{}
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		pr.errorLog.Println(err)
		return post, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM posts WHERE id = ?;")
	if err != nil {
		pr.errorLog.Println(err)
		return post, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		pr.errorLog.Println(err)
		return post, err
	}
	if rows.Next() {
		rows.Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content)
	}
	if err = tx.Commit(); err != nil {
		pr.errorLog.Println(err)
		return entity.Post{}, err
	}
	return post, nil
}

func (pr *PostsRepository) FetchAll(ctx context.Context) ([]entity.Post, error) {
	posts := []entity.Post{}
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		pr.errorLog.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM posts;")
	if err != nil {
		pr.errorLog.Println(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		pr.errorLog.Println(err)
		return nil, err
	}
	for rows.Next() {
		post := entity.Post{}
		rows.Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content)
		posts = append(posts, post)
	}
	if err = tx.Commit(); err != nil {
		pr.errorLog.Println(err)
		return nil, err
	}
	return posts, nil
}

func (pr *PostsRepository) FetchByUserId(ctx context.Context, id int) ([]entity.Post, error) {
	posts := []entity.Post{}
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		pr.errorLog.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM posts WHERE user_id = ?;")
	if err != nil {
		pr.errorLog.Println(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		pr.errorLog.Println(err)
		return nil, err
	}
	for rows.Next() {
		post := entity.Post{}
		rows.Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content)
		posts = append(posts, post)
	}
	if err = tx.Commit(); err != nil {
		pr.errorLog.Println(err)
		return nil, err
	}
	return posts, nil
}

func (pr *PostsRepository) FetchByCategoryId(ctx context.Context, id int) ([]entity.Post, error) {
	posts := []entity.Post{}
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		pr.errorLog.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT post_id FROM post_categories WHERE category_id = ?")
	if err != nil {
		pr.errorLog.Println(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		pr.errorLog.Println(err)
		return nil, err
	}
	for rows.Next() {
		post := entity.Post{}
		rows.Scan(&post.Id)
		posts = append(posts, post)
	}
	if err = tx.Commit(); err != nil {
		pr.errorLog.Println(err)
		return nil, err
	}
	return posts, nil
}

func (pr *PostsRepository) Store(ctx context.Context, post entity.Post) (int64, error) {
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		pr.errorLog.Println(err)
		return 0, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO posts(user_id, date, title, content) 
		VALUES(?,?,?,?);`)
	if err != nil {
		pr.errorLog.Println(err)
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, post.User.Id, time.Now().Format("2006-01-02"), post.Title, post.Content)
	if err != nil {
		pr.errorLog.Println(err)
		return 0, err
	}
	stmt_cat, err := tx.PrepareContext(ctx, `INSERT INTO post_categories(post_id, category_id) VALUES(?,?);`)
	if err != nil {
		pr.errorLog.Println(err)
		return 0, err
	}
	defer stmt_cat.Close()
	post_id, err := res.LastInsertId()
	if err != nil {
		pr.errorLog.Println(err)
		return 0, err
	}
	for _, category := range post.Category {
		_, err = stmt_cat.ExecContext(ctx, post_id, category.Id)
		if err != nil {
			pr.errorLog.Println(err)
			return 0, err
		}
	}
	if err = tx.Commit(); err != nil {
		pr.errorLog.Println(err)
		return 0, err
	}
	return post_id, nil
}

// for future use
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
