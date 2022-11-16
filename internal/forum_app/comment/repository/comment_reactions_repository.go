package repository

import (
	"database/sql"
	"errors"
	"forum/internal/forum_app/entity"
)

type CommentReactionsRepository struct {
	db *sql.DB
}

func NewCommentReactionsRepository(db *sql.DB) *CommentReactionsRepository {
	return &CommentReactionsRepository{db}
}

func (crr *CommentReactionsRepository) FetchByCommentId(id int, like bool) ([]entity.Reaction, error) {
	reactions := []entity.Reaction{}
	tx, err := crr.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("SELECT user_id, date, like FROM comment_reactions WHERE comment_id = ? AND like = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id, like)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		reaction := entity.Reaction{}
		rows.Scan(&reaction.User.Id, &reaction.Date, &reaction.Like)
		reactions = append(reactions, reaction)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return reactions, nil
}

func (crr *CommentReactionsRepository) StoreReaction(commentReaction entity.CommentReaction) error {
	tx, err := crr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`INSERT INTO comment_reactions(comment_id, user_id, date, like) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(commentReaction.Comment.Id, commentReaction.Reaction.User.Id, commentReaction.Reaction.Date, commentReaction.Reaction.Like); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (crr *CommentReactionsRepository) UpdateReaction(commentReaction entity.CommentReaction) error {
	tx, err := crr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`UPDATE comment_reactions SET like = ? WHERE comment_id = ? AND user_id = ?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(commentReaction.Like, commentReaction.Post.Id, commentReaction.Reaction.User.Id)
	if err != nil {
		return err
	}
	rAffected, err := res.RowsAffected()
	if rAffected > 1 {
		return errors.New("more than one row has been affected")
	}
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (crr *CommentReactionsRepository) DeleteReaction(commentReaction entity.CommentReaction) error {
	tx, err := crr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`DELETE FROM comment_reactions WHERE comment_id = ? AND user_id = ?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(commentReaction.Post.Id, commentReaction.Reaction.User.Id)
	if err != nil {
		return err
	}
	rAffected, err := res.RowsAffected()
	if rAffected > 1 {
		return errors.New("more than one row has been affected")
	}
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (crr *CommentReactionsRepository) FetchByUserId(userId int, like bool) ([]entity.CommentReaction, error) {
	commentReactions := []entity.CommentReaction{}
	tx, err := crr.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("SELECT comment_id, date, like FROM comment_reactions WHERE user_id = ? AND like = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId, like)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		commentReaction := entity.CommentReaction{}
		rows.Scan(&commentReaction.Comment.Id, &commentReaction.Reaction.Date, &commentReaction.Reaction.Like)
		commentReactions = append(commentReactions, commentReaction)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return commentReactions, nil
}
