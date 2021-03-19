package inmem

import (
	"context"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/nacobas/customer/customer"
	"github.com/nacobas/customer/registry"
)

var (
	ErrUsedID   = errors.New("ID allready in use")
	ErrConflict = errors.New("Unique Data conflict")
	ErrNotFound = errors.New("Not found")
)

func NewRepo() registry.Repo {
	return &repo{
		mtx:  sync.RWMutex{},
		data: map[uint32]customer.Customer{},
	}
}

func NewRepoWithSeed(seed []customer.Customer) registry.Repo {

	var data = map[uint32]customer.Customer{}

	for _, c := range seed {
		data[c.ID] = c
	}

	return &repo{
		mtx:  sync.RWMutex{},
		data: data,
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
		return errors.Wrap(ErrUsedID, op)
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
		return errors.Wrap(ErrNotFound, op)
	}

	r.data[c.ID] = *c

	return nil
}
