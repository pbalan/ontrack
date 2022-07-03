package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"strconv"

	"github.com/pbalan/ontrack/src/graph/generated"
	"github.com/pbalan/ontrack/src/graph/model"
	"github.com/pbalan/ontrack/src/repositories"
	"github.com/pbalan/ontrack/src/services"
	"github.com/pbalan/ontrack/src/utils"
	log "github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input model.NewUser) (interface{}, error) {
	// Check Email
	_, err := repositories.UserGetByEmail(r.DB, input.Email)
	if err == nil {
		// if err != record not found
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	password, err := utils.HashPassword(input.Password)
	if err != nil {
		log.Fatal("Unable to generate password")
	}
	user := model.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  password,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}
	createdUser, err := repositories.CreateUser(r.DB, &user)
	if err != nil {
		return nil, err
	}

	token, err := services.JwtGenerate(ctx, strconv.FormatInt(createdUser.ID, 10))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, username *string, email *string, password string) (interface{}, error) {
	var getUser *model.User
	var err error

	if email == nil && username == nil {
		return nil, &gqlerror.Error{
			Message: "Username/Email not specified",
		}
	}
	if utils.DerefString(email) == "" {
		getUser, err = repositories.UserGetByUsername(r.DB, utils.DerefString(username))
	}
	if utils.DerefString(username) == "" {
		getUser, err = repositories.UserGetByEmail(r.DB, utils.DerefString(email))
	}
	if err != nil {
		// if user not found
		if err == gorm.ErrRecordNotFound {
			return nil, &gqlerror.Error{
				Message: "Email not found",
			}
		}
		return nil, err
	}

	if err := utils.CheckPasswordHash(password, getUser.Password); err == false {
		return nil, errors.New("unable to login")
	}

	token, err := services.JwtGenerate(ctx, strconv.FormatInt(getUser.ID, 10))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	r.DB.Find(&users)

	return users, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
