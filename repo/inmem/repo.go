package inmem

import (
	"context"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/nacobas/customer/customer"
	"github.com/nacobas/customer/registry"
)

func NewRepo() registry.Repo {
	return &repo{
		mtx:  sync.RWMutex{},
		data: map[uint32]customer.Customer{},
	}
}

type repo struct {
	mtx  sync.RWMutex
	data map[uint32]customer.Customer
}

func (r *repo) Get(ctx context.Context, id uint32) (*customer.Customer, error) {
	const op string = "inmem.repo.Get"

	r.mtx.RLock()
	defer r.mtx.RUnlock()

	c, ok := r.data[id]
	if !ok {
		return nil, errors.Wrap(registry.ErrNotFound, op)
	}

	return &c, nil
}

func (r *repo) Insert(ctx context.Context, c *customer.Customer) error {
	const op string = "inmem.repo.New"

	r.mtx.Lock()
	defer r.mtx.Unlock()

	_, ok := r.data[c.ID]
	if ok {
		return errors.Wrap(registry.ErrUsedID, op)
	}

	r.data[c.ID] = *c

	return nil
}

func (r *repo) Update(ctx context.Context, c *customer.Customer) error {
	const op string = "inmem.repo.Update"

	r.mtx.Lock()
	defer r.mtx.Unlock()

	_, ok := r.data[c.ID]
	if !ok {
		return errors.Wrap(registry.ErrNotFound, op)
	}

	r.data[c.ID] = *c

	return nil
}
