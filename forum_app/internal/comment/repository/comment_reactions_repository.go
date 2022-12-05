package repository

import (
	"context"
	"database/sql"
	"errors"
	"forum_app/internal/entity"
	"log"
)

type CommentReactionsRepository struct {
	db       *sql.DB
	errorLog *log.Logger
}

func NewCommentReactionsRepository(db *sql.DB, errorLog *log.Logger) *CommentReactionsRepository {
	return &CommentReactionsRepository{db, errorLog}
}

func (crr *CommentReactionsRepository) FetchByCommentId(ctx context.Context, id int, like bool) ([]entity.Reaction, error) {
	reactions := []entity.Reaction{}
	tx, err := crr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		crr.errorLog.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT user_id, date, like FROM comment_reactions WHERE comment_id = ? AND like = ?;")
	if err != nil {
		crr.errorLog.Println(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id, like)
	if err != nil {
		crr.errorLog.Println(err)
		return nil, err
	}
	for rows.Next() {
		reaction := entity.Reaction{}
		rows.Scan(&reaction.User.Id, &reaction.Date, &reaction.Like)
		reactions = append(reactions, reaction)
	}
	if err = tx.Commit(); err != nil {
		crr.errorLog.Println(err)
		return nil, err
	}
	return reactions, nil
}

func (crr *CommentReactionsRepository) StoreReaction(ctx context.Context, commentReaction entity.CommentReaction) error {
	tx, err := crr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		crr.errorLog.Println(err)
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO comment_reactions(comment_id, user_id, date, like) VALUES(?, ?, ?, ?)`)
	if err != nil {
		crr.errorLog.Println(err)
		return err
	}
	defer stmt.Close()
	if _, err = stmt.ExecContext(ctx, commentReaction.Comment.Id, commentReaction.Reaction.User.Id, commentReaction.Reaction.Date, commentReaction.Reaction.Like); err != nil {
		crr.errorLog.Println(err)
		return err
	}
	if err = tx.Commit(); err != nil {
		crr.errorLog.Println(err)
		return err
	}
	return nil
}

func (crr *CommentReactionsRepository) UpdateReaction(ctx context.Context, commentReaction entity.CommentReaction) error {
	tx, err := crr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		crr.errorLog.Println(err)
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `UPDATE comment_reactions SET like = ? WHERE comment_id = ? AND user_id = ?;`)
	if err != nil {
		crr.errorLog.Println(err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, commentReaction.Like, commentReaction.Post.Id, commentReaction.Reaction.User.Id)
	if err != nil {
		crr.errorLog.Println(err)
		return err
	}
	rAffected, err := res.RowsAffected()
	if rAffected > 1 {
		crr.errorLog.Println(errors.New("more than one row has been affected"))
		return errors.New("more than one row has been affected")
	}
	if err != nil {
		crr.errorLog.Println(err)
		return err
	}
	if err = tx.Commit(); err != nil {
		crr.errorLog.Println(err)
		return err
	}
	return nil
}

func (crr *CommentReactionsRepository) DeleteReaction(ctx context.Context, commentReaction entity.CommentReaction) error {
	tx, err := crr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		crr.errorLog.Println(err)
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM comment_reactions WHERE comment_id = ? AND user_id = ?;`)
	if err != nil {
		crr.errorLog.Println(err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, commentReaction.Post.Id, commentReaction.Reaction.User.Id)
	if err != nil {
		crr.errorLog.Println(err)
		return err
	}
	rAffected, err := res.RowsAffected()
	if rAffected > 1 {
		crr.errorLog.Println(errors.New("more than one row has been affected"))
		return errors.New("more than one row has been affected")
	}
	if err != nil {
		crr.errorLog.Println(err)
		return err
	}
	if err = tx.Commit(); err != nil {
		crr.errorLog.Println(err)
		return err
	}
	return nil
}

func (crr *CommentReactionsRepository) FetchByUserId(ctx context.Context, userId int, like bool) ([]entity.CommentReaction, error) {
	commentReactions := []entity.CommentReaction{}
	tx, err := crr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		crr.errorLog.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT comment_id, date, like FROM comment_reactions WHERE user_id = ? AND like = ?;")
	if err != nil {
		crr.errorLog.Println(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userId, like)
	if err != nil {
		crr.errorLog.Println(err)
		return nil, err
	}
	for rows.Next() {
		commentReaction := entity.CommentReaction{}
		rows.Scan(&commentReaction.Comment.Id, &commentReaction.Reaction.Date, &commentReaction.Reaction.Like)
		commentReactions = append(commentReactions, commentReaction)
	}
	if err = tx.Commit(); err != nil {
		crr.errorLog.Println(err)
		return nil, err
	}
	return commentReactions, nil
}
