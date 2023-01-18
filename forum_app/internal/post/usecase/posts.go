package usecase

import (
	"context"
	"forum_app/internal/entity"
	"log"
)

type PostsUsecase struct {
	postsRepo         PostsRepository
	postReactionsRepo PostReactionsRepository
	commentsRepo      CommentsRepository
	categoriesRepo    CategoriesRepository
	usersRepo         UsersRepository
	errorLog          *log.Logger
}

func NewPostsUsecase(postsRepo PostsRepository, postReactionsRepo PostReactionsRepository, commentsRepo CommentsRepository, categoriesRepo CategoriesRepository, usersRepo UsersRepository, errorLog *log.Logger) *PostsUsecase {
	return &PostsUsecase{
		postsRepo:         postsRepo,
		postReactionsRepo: postReactionsRepo,
		commentsRepo:      commentsRepo,
		categoriesRepo:    categoriesRepo,
		usersRepo:         usersRepo,
		errorLog:          errorLog,
	}
}

func (u *PostsUsecase) FetchById(ctx context.Context, id int, postRes chan entity.PostResult) {
	post, err := u.postsRepo.FetchById(ctx, id)
	if err != nil {
		u.errorLog.Println(err)
		postRes <- entity.PostResult{Err: err}
		return
	}
	if post.Id == 0 {
		u.errorLog.Println(entity.ErrPostNotFound)
		postRes <- entity.PostResult{Err: entity.ErrPostNotFound}
		return
	}
	u.fetchPostDetails(ctx, &post)
	postRes <- entity.PostResult{Post: post}
}

func (u *PostsUsecase) fetchPostDetails(ctx context.Context, post *entity.Post) {
	var (
		err           error
		user          = make(chan entity.User)
		comments      = make(chan []entity.Comment)
		categories    = make(chan []entity.Category)
		likes         = make(chan []entity.Reaction)
		dislikes      = make(chan []entity.Reaction)
		errUser       = make(chan error)
		errComments   = make(chan error)
		errCategories = make(chan error)
		errLikes      = make(chan error)
		errDislikes   = make(chan error)
	)
	go u.fetchUser(ctx, post.User.Id, user, errUser)
	go u.fetchCategories(ctx, post.Id, categories, errCategories)
	go u.fetchComments(ctx, post.Id, comments, errComments)
	go u.fetchLikes(ctx, post.Id, likes, errLikes)
	go u.fetchDislikes(ctx, post.Id, dislikes, errDislikes)
	for i := 0; i < 5; i++ {
		select {
		case post.User = <-user:
			if err = <-errUser; err != nil {
				u.errorLog.Println(err)
			}
		case post.Category = <-categories:
			if err = <-errCategories; err != nil {
				u.errorLog.Println(err)
			}
		case post.Comments = <-comments:
			if err = <-errComments; err != nil {
				u.errorLog.Println(err)
			}
		case post.Likes = <-likes:
			if err = <-errLikes; err != nil {
				u.errorLog.Println(err)
			}
		case post.Dislikes = <-dislikes:
			if err = <-errDislikes; err != nil {
				u.errorLog.Println(err)
			}
		}
	}
	post.CountTotals()
}

func (u *PostsUsecase) FetchAll(ctx context.Context, postsRes chan entity.PostsResult) {
	posts, err := u.postsRepo.FetchAll(ctx)
	if err != nil {
		postsRes <- entity.PostsResult{Err: err}
	}
	for ix := range posts {
		u.fetchPostDetails(ctx, &posts[ix])
		posts[ix].Comments, posts[ix].Likes, posts[ix].Dislikes = nil, nil, nil
	}
	postsRes <- entity.PostsResult{Posts: posts}
}

func (u *PostsUsecase) fetchUser(ctx context.Context, id int, user chan entity.User, errUser chan error) {
	tempUser, err := u.usersRepo.FetchById(ctx, id)
	user <- tempUser
	errUser <- err
}

func (u *PostsUsecase) fetchComments(ctx context.Context, id int, comments chan []entity.Comment, errComments chan error) {
	tempComments, err := u.commentsRepo.FetchByPostId(ctx, id)
	comments <- tempComments
	errComments <- err
}

func (u *PostsUsecase) fetchCategories(ctx context.Context, id int, categories chan []entity.Category, errCategories chan error) {
	tempCategories, err := u.categoriesRepo.FetchByPostId(ctx, id)
	categories <- tempCategories
	errCategories <- err
}

func (u *PostsUsecase) fetchLikes(ctx context.Context, id int, likes chan []entity.Reaction, errLikes chan error) {
	tempLikes, err := u.postReactionsRepo.FetchByPostId(ctx, id, true)
	likes <- tempLikes
	errLikes <- err
}

func (u *PostsUsecase) fetchDislikes(ctx context.Context, id int, dislikes chan []entity.Reaction, errDislikes chan error) {
	tempDislikes, err := u.postReactionsRepo.FetchByPostId(ctx, id, false)
	dislikes <- tempDislikes
	errDislikes <- err
}

func (u *PostsUsecase) FetchCategoryPosts(ctx context.Context, id int, catRes chan entity.CatResult) {
	var err error
	category, err := u.categoriesRepo.FetchById(ctx, id)
	if err != nil {
		catRes <- entity.CatResult{Err: err}
		return
	}
	category.Posts, err = u.postsRepo.FetchByCategoryId(ctx, category.Id)
	if err != nil {
		catRes <- entity.CatResult{Err: err}
		return
	}
	if category.Id == 0 {
		catRes <- entity.CatResult{Err: entity.ErrCategoryNotFound}
		return
	}
	for ix := range category.Posts {
		category.Posts[ix], err = u.postsRepo.FetchById(ctx, category.Posts[ix].Id)
		u.fetchPostDetails(ctx, &category.Posts[ix])
		category.Posts[ix].Comments = nil
		category.Posts[ix].Likes = nil
		category.Posts[ix].Dislikes = nil
	}
	category.CountTotals()
	catRes <- entity.CatResult{Cat: category}
}

func (u *PostsUsecase) Store(ctx context.Context, post entity.Post, res chan entity.Result) {
	id, err := u.postsRepo.Store(ctx, post)
	if err != nil {
		res <- entity.Result{Err: err}
	}
	res <- entity.Result{Id: id}
}

func (u *PostsUsecase) StorePostReaction(ctx context.Context, postReaction entity.PostReaction, err chan error) {
	err <- u.postReactionsRepo.StoreReaction(ctx, postReaction)
}

func (u *PostsUsecase) UpdatePostReaction(ctx context.Context, postReaction entity.PostReaction, err chan error) {
	err <- u.postReactionsRepo.UpdateReaction(ctx, postReaction)
}

func (u *PostsUsecase) DeletePostReaction(ctx context.Context, postReaction entity.PostReaction, err chan error) {
	err <- u.postReactionsRepo.DeleteReaction(ctx, postReaction)
}
