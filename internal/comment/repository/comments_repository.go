package repository

import (
	"database/sql"
	"forum/internal/entity"
)

type CommentsRepository struct {
	db *sql.DB
}

func NewCommentsRepository(db *sql.DB) *CommentsRepository {
	return &CommentsRepository{db}
}

func (cr *CommentsRepository) FetchById(id int) (entity.Comment, error) {
	comment := entity.Comment{}
	tx, err := cr.db.Begin()
	if err != nil {
		return comment, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("SELECT * FROM comments WHERE id = ?;")
	if err != nil {
		return comment, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return comment, err
	}
	if rows.Next() {
		rows.Scan(&comment.Id, &comment.Post.Id, &comment.User.Id, &comment.Date, &comment.Content)
	}
	if err = tx.Commit(); err != nil {
		return entity.Comment{}, err
	}
	return comment, nil
}

func (cr *CommentsRepository) FetchByPostId(id int) ([]entity.Comment, error) {
	comments := []entity.Comment{}
	tx, err := cr.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("SELECT * FROM comments WHERE post_id = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		comment := entity.Comment{}
		rows.Scan(&comment.Id, &comment.Post.Id, &comment.User.Id, &comment.Date, &comment.Content)
		comments = append(comments, comment)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *CommentsRepository) FetchByUserId(id int) ([]entity.Comment, error) {
	comments := []entity.Comment{}
	tx, err := cr.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("SELECT * FROM comments WHERE user_id = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		comment := entity.Comment{}
		rows.Scan(&comment.Id, &comment.Post.Id, &comment.User.Id, &comment.Date, &comment.Content)
		comments = append(comments, comment)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *CommentsRepository) Store(comment entity.Comment) (int64, error) {
	tx, err := cr.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`INSERT INTO comments(post_id, user_id, date, content) VALUES(?,?,?,?);`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(comment.Post.Id, comment.User.Id, comment.Date, comment.Content)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
