package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/pbalan/ontrack/src/graph/generated"
	"github.com/pbalan/ontrack/src/graph/model"
	"github.com/pbalan/ontrack/src/utils"
	log "github.com/sirupsen/logrus"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input model.NewUser) (*model.User, error) {
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

	r.DB.Create(&user)

	return &user, nil
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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}
