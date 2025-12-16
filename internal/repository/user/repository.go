package user

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/entity"
	postUserRegister "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserRegister/v0"
)

type Repository struct {
	db *db.PGConnect
}

func New(db *db.PGConnect) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) UsersCreate(ctx context.Context, userName string, user entity.User) (resp *postUserRegister.Response, err error) {
	err = r.db.QueryWithOneResultJSON(
		ctx,
		&resp,
		"select orders.users_create(_username=>$1, _user=>$2);",
		userName, user,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Repository) UsersGetByUserName(ctx context.Context, userName string) (user *entity.User, err error) {
	err = r.db.QueryWithOneResultJSON(
		ctx,
		&user,
		"select orders.users_get_by_username(_username=>$1);",
		userName,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
