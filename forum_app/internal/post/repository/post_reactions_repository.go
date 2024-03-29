package repository

import (
	"context"
	"database/sql"
	"errors"
	"forum_app/internal/entity"
	"log"
	"time"
)

type PostReactionsRepository struct {
	db       *sql.DB
	errorLog *log.Logger
}

func NewPostReactionsRepository(db *sql.DB, errorLog *log.Logger) *PostReactionsRepository {
	return &PostReactionsRepository{db, errorLog}
}

func (rr *PostReactionsRepository) FetchByPostId(ctx context.Context, id int, like bool) ([]entity.Reaction, error) {
	reactions := []entity.Reaction{}
	tx, err := rr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		rr.errorLog.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT user_id, date, like FROM post_reactions WHERE post_id = ? and like = ?")
	if err != nil {
		rr.errorLog.Println(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id, like)
	if err != nil {
		rr.errorLog.Println(err)
		return nil, err
	}
	for rows.Next() {
		reaction := entity.Reaction{}
		rows.Scan(&reaction.User.Id, &reaction.Date, &reaction.Like)
		reactions = append(reactions, reaction)
	}
	if err = tx.Commit(); err != nil {
		rr.errorLog.Println(err)
		return nil, err
	}
	return reactions, nil
}

func (rr *PostReactionsRepository) StoreReaction(ctx context.Context, postReaction entity.PostReaction) error {
	tx, err := rr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		rr.errorLog.Println(err)
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO post_reactions(post_id, user_id, date, like) VALUES(?, ?, ?, ?)`)
	if err != nil {
		rr.errorLog.Println(err)
		return err
	}
	defer stmt.Close()
	postReaction.Reaction.Date = time.Now().Format("2006-01-02")
	if _, err = stmt.ExecContext(ctx, postReaction.Post.Id, postReaction.Reaction.User.Id, postReaction.Reaction.Date, postReaction.Reaction.Like); err != nil {
		rr.errorLog.Println(err)
		return err
	}
	if err = tx.Commit(); err != nil {
		rr.errorLog.Println(err)
		return err
	}
	return nil
}

func (rr *PostReactionsRepository) UpdateReaction(ctx context.Context, postReaction entity.PostReaction) error {
	tx, err := rr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		rr.errorLog.Println(err)
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `UPDATE post_reactions SET like = ?, date = ? WHERE post_id = ? AND user_id = ?;`)
	if err != nil {
		rr.errorLog.Println(err)
		return err
	}
	defer stmt.Close()
	postReaction.Reaction.Date = time.Now().Format("2006-01-02")
	res, err := stmt.ExecContext(ctx, postReaction.Like, postReaction.Reaction.Date, postReaction.Post.Id, postReaction.Reaction.User.Id)
	if err != nil {
		rr.errorLog.Println(err)
		return err
	}
	rAffected, err := res.RowsAffected()
	if rAffected > 1 {
		rr.errorLog.Println(errors.New("more than one row has been affected"))
		return errors.New("more than one row has been affected")
	} else if rAffected == 0 {
		rr.errorLog.Println(errors.New("no row has been affected"))
		return errors.New("no row has been affected")
	}
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rr *PostReactionsRepository) DeleteReaction(ctx context.Context, postReaction entity.PostReaction) error {
	tx, err := rr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		rr.errorLog.Println(err)
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM post_reactions WHERE post_id = ? AND user_id = ? AND like = ?;`)
	if err != nil {
		rr.errorLog.Println(err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, postReaction.Post.Id, postReaction.Reaction.User.Id, postReaction.Reaction.Like)
	if err != nil {
		rr.errorLog.Println(err)
		return err
	}
	rAffected, err := res.RowsAffected()
	if rAffected > 1 {
		rr.errorLog.Println(errors.New("more than one row has been affected"))
		return errors.New("more than one row has been affected")
	}
	if err != nil {
		rr.errorLog.Println(err)
		return err
	}
	if err = tx.Commit(); err != nil {
		rr.errorLog.Println(err)
		return err
	}
	return nil
}

func (rr *PostReactionsRepository) FetchByUserId(ctx context.Context, userId int, like bool) ([]entity.PostReaction, error) {
	postReactions := []entity.PostReaction{}
	tx, err := rr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		rr.errorLog.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT post_id, date, like FROM post_reactions where user_id = ? and like = ?;")
	if err != nil {
		rr.errorLog.Println(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userId, like)
	if err != nil {
		rr.errorLog.Println(err)
		return nil, err
	}
	for rows.Next() {
		postReaction := entity.PostReaction{}
		rows.Scan(&postReaction.Post.Id, &postReaction.Reaction.Date, &postReaction.Reaction.Like)
		postReactions = append(postReactions, postReaction)
	}
	if err = tx.Commit(); err != nil {
		rr.errorLog.Println(err)
		return nil, err
	}
	return postReactions, nil
}
