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

func (u *PostsUsecase) FetchById(ctx context.Context, id int) (entity.Post, error) {
	post, err := u.postsRepo.FetchById(ctx, id)
	if err != nil {
		u.errorLog.Println(err)
		return entity.Post{}, err
	}
	u.fetchPostDetails(ctx, &post)
	return post, nil
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

func (u *PostsUsecase) FetchAll(ctx context.Context) ([]entity.Post, error) {
	posts, err := u.postsRepo.FetchAll(ctx)
	if err != nil {
		u.errorLog.Println(err)
		return nil, err
	}
	for ix := range posts {
		u.fetchPostDetails(ctx, &posts[ix])
		posts[ix].Comments, posts[ix].Likes, posts[ix].Dislikes = nil, nil, nil
	}
	return posts, nil
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

func (u *PostsUsecase) FetchCategoryPosts(ctx context.Context, id int) (entity.Category, error) {
	var err error
	category, err := u.categoriesRepo.FetchById(ctx, id)
	if err != nil {
		u.errorLog.Println(err)
		return entity.Category{}, err
	}
	category.Posts, err = u.postsRepo.FetchByCategoryId(ctx, category.Id)
	if err != nil {
		u.errorLog.Println(err)
		return category, err
	}
	for ix, post := range category.Posts {
		category.Posts[ix], err = u.postsRepo.FetchById(ctx, post.Id)
		if err != nil {
			u.errorLog.Println(err)
		}
		category.Posts[ix].User, err = u.usersRepo.FetchById(ctx, category.Posts[ix].User.Id)
		if err != nil {
			u.errorLog.Println(err)
		}

	}
	category.CountTotals()
	return category, nil
}

// func (u *PostsUsecase) FetchAllSorted() ([]entity.Post, error) {
// 	posts, err := u.postsRepo.FetchAllSorted()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return posts, nil
// }

func (u *PostsUsecase) Store(ctx context.Context, post entity.Post) (int64, error) {
	id, err := u.postsRepo.Store(ctx, post)
	if err != nil {
		u.errorLog.Println(err)
		return 0, err
	}
	return id, nil
}

func (u *PostsUsecase) StorePostReaction(ctx context.Context, postReaction entity.PostReaction) error {
	return u.postReactionsRepo.StoreReaction(ctx, postReaction)
}

func (u *PostsUsecase) UpdatePostReaction(ctx context.Context, postReaction entity.PostReaction) error {
	return u.postReactionsRepo.UpdateReaction(ctx, postReaction)
}

func (u *PostsUsecase) DeletePostReaction(ctx context.Context, postReaction entity.PostReaction) error {
	return u.postReactionsRepo.DeleteReaction(ctx, postReaction)
}

// for future use
// func (u *PostsUsecase) Update(post entity.Post) error {
// 	err := u.postsRepo.Update(post)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (u *PostsUsecase) DeleteById(id int) error {
// 	err := u.postsRepo.Delete(id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
