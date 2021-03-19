package registry

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
	"github.com/nacobas/customer/customer"
)

var (
	ErrNotFound   = errors.New("Not found")
	ErrValidation = errors.New("Input validation failed")
	ErrExpected   = errors.New("Expected error")
	ErrUnexpected = errors.New("Unexpected error")
)

func NewService(r Repo) Service {

	return &service{r, customer.NewValidator()}
}

type Service interface {
	Get(ctx context.Context, id uint32) (*customer.Customer, error)
	New(ctx context.Context, i customer.Info) (*customer.Customer, error)
	UpdateInfo(ctx context.Context, id uint32, i customer.Info) (*customer.Customer, error)
	SetState(ctx context.Context, id uint32, s customer.State) error
}

type Repo interface {
	Get(ctx context.Context, id uint32) (*customer.Customer, error)
	Insert(ctx context.Context, c *customer.Customer) error
	Update(ctx context.Context, c *customer.Customer) error
}

type service struct {
	repo     Repo
	validate *validator.Validate
}

func (svc *service) Get(ctx context.Context, id uint32) (*customer.Customer, error) {
	const op string = "registry.Service.Get"

	c, err := svc.repo.Get(ctx, id)
	if err != nil {
		return nil, errors.Mark(errors.Wrap(err, op), ErrNotFound)
	}

	return c, nil
}

func (svc *service) New(ctx context.Context, i customer.Info) (*customer.Customer, error) {
	const op string = "registry.Service.New"

	if err := svc.validate.Struct(i); err != nil {
		return nil, errors.Mark(errors.Wrap(err, op), ErrValidation)
	}

	c := customer.NewWithRandomID(i)

	if err := svc.repo.Insert(ctx, c); err != nil {
		return nil, errors.Mark(errors.Wrap(err, op), ErrUnexpected)
	}

	return c, nil
}

func (svc *service) UpdateInfo(ctx context.Context, id uint32, i customer.Info) (*customer.Customer, error) {
	const op string = "registry.Service.UpdateInfo"

	if err := svc.validate.Struct(i); err != nil {
		return nil, errors.Mark(errors.Wrap(err, op), ErrValidation)
	}

	c, err := svc.repo.Get(ctx, id)
	if err != nil {
		return nil, errors.Mark(errors.Wrap(err, op), ErrNotFound)
	}

	if err := c.UpdateInfo(i); err != nil {
		return nil, errors.Mark(errors.Wrap(err, op), ErrExpected)
	}

	if err = svc.repo.Update(ctx, c); err != nil {
		return nil, errors.Mark(errors.Wrap(err, op), ErrUnexpected)
	}

	return c, nil
}

func (svc *service) SetState(ctx context.Context, id uint32, s customer.State) error {
	const op string = "registry.Service.SetState"

	if err := svc.validate.Var(s, "min=1,max=3"); err != nil {
		return errors.Mark(errors.Wrap(err, op), ErrValidation)
	}

	c, err := svc.repo.Get(ctx, id)
	if err != nil {
		return errors.Mark(errors.Wrap(err, op), ErrNotFound)
	}

	c.State = s

	return errors.Mark(errors.Wrap(svc.repo.Update(ctx, c), op), ErrUnexpected)
}
