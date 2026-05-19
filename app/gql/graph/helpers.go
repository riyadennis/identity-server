package graph

import (
	"context"
	"errors"

	"github.com/riyadennis/identity-server/app/gql/graph/model"
	"github.com/riyadennis/identity-server/business"
	"github.com/riyadennis/identity-server/business/store"
	"github.com/riyadennis/identity-server/business/validation"
)

func (r *mutationResolver) insertUser(ctx context.Context, input model.RegisterInput, createdBy string) (*model.RegisterResponse, error) {
	u := &store.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
		Terms:     input.Terms,
		CreatedBy: createdBy,
		Active:    true,
	}
	if input.Company != nil {
		u.Company = *input.Company
	}
	if input.PostCode != nil {
		u.PostCode = *input.PostCode
	}

	if err := validation.ValidateUser(u); err != nil {
		return nil, err
	}

	existing, err := r.Store.Read(ctx, u.Email)
	if err != nil {
		return nil, err
	}
	if existing.Email == u.Email {
		return nil, errors.New("email already exists")
	}

	u.Password, err = business.EncryptPassword(u.Password)
	if err != nil {
		return nil, err
	}

	created, err := r.Store.Insert(ctx, u)
	if err != nil {
		return nil, err
	}

	return &model.RegisterResponse{
		ID:        &created.ID,
		FirstName: &created.FirstName,
		LastName:  &created.LastName,
		Email:     &created.Email,
		Company:   &created.Company,
		PostCode:  &created.PostCode,
		Terms:     &created.Terms,
		CreatedAt: &created.CreatedAt,
	}, nil
}
