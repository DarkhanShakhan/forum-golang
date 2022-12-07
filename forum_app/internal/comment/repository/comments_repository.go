package repository

import (
	"context"
	"database/sql"
	"forum_app/internal/entity"
	"log"
	"time"
)

type CommentsRepository struct {
	db       *sql.DB
	errorLog *log.Logger
}

func NewCommentsRepository(db *sql.DB, errorLog *log.Logger) *CommentsRepository {
	return &CommentsRepository{db, errorLog}
}

func (cr *CommentsRepository) FetchById(ctx context.Context, id int) (entity.Comment, error) {
	comment := entity.Comment{}
	tx, err := cr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		cr.errorLog.Println(err)
		return comment, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM comments WHERE id = ?;")
	if err != nil {
		cr.errorLog.Println(err)
		return comment, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		cr.errorLog.Println(err)
		return comment, err
	}
	if rows.Next() {
		rows.Scan(&comment.Id, &comment.Post.Id, &comment.User.Id, &comment.Date, &comment.Content)
	}
	if err = tx.Commit(); err != nil {
		cr.errorLog.Println(err)
		return entity.Comment{}, err
	}
	return comment, nil
}

func (cr *CommentsRepository) FetchByPostId(ctx context.Context, id int) ([]entity.Comment, error) {
	comments := []entity.Comment{}
	tx, err := cr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		cr.errorLog.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT c.id, c.user_id, u.name, c.date, c.content FROM comments AS c LEFT JOIN users AS u ON c.user_id=u.id WHERE post_id = ?;")
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
		comment := entity.Comment{}
		rows.Scan(&comment.Id, &comment.User.Id, &comment.User.Name, &comment.Date, &comment.Content)
		comments = append(comments, comment)
	}
	if err = tx.Commit(); err != nil {
		cr.errorLog.Println(err)
		return nil, err
	}
	return comments, nil
}

func (cr *CommentsRepository) FetchByUserId(ctx context.Context, id int) ([]entity.Comment, error) {
	comments := []entity.Comment{}
	tx, err := cr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		cr.errorLog.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT c.id, c.post_id, p.title, c.date, c.content FROM comments AS c LEFT JOIN posts AS p ON c.post_id=p.id WHERE c.user_id = ?;")
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
	if rows.Next() {
		comment := entity.Comment{}
		rows.Scan(&comment.Id, &comment.Post.Id, &comment.Post.Title, &comment.Date, &comment.Content)
		comments = append(comments, comment)
	}
	if err = tx.Commit(); err != nil {
		cr.errorLog.Println(err)
		return nil, err
	}
	return comments, nil
}

func (cr *CommentsRepository) Store(ctx context.Context, comment entity.Comment) (int64, error) {
	tx, err := cr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		cr.errorLog.Println(err)
		return 0, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO comments(post_id, user_id, date, content) VALUES(?,?,?,?);`)
	if err != nil {
		cr.errorLog.Println(err)
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, comment.Post.Id, comment.User.Id, time.Now().Format("2006-01-02"), comment.Content)
	if err != nil {
		cr.errorLog.Println(err)
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		cr.errorLog.Println(err)
		return 0, err
	}
	return res.LastInsertId()
}
