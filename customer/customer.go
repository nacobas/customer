package customer

import (
	"math/rand"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/Azure/go-autorest/autorest/date"
)

var (
	ErrTypeNotEqual = errors.New("Type not equal")
)

func New(id uint32, i Info) *Customer {
	return &Customer{
		ID:    id,
		State: Prospect,
		Info:  i,
	}
}

func NewWithRandomID(i Info) *Customer {

	return New(NewRandomID(), i)
}

func NewRandomID() uint32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return r.Uint32()
}

type Customer struct {
	ID    uint32 `validate:"required"`
	State State  `validate:"min=1,max=3"`
	Info  `validate:"required"`
}

func (c *Customer) UpdateInfo(i Info) error {

	if c.Type() != i.Type() {
		return ErrTypeNotEqual
	}

	c.Info = i

	return nil
}

type Info interface {
	Type() CustomerType
}

type PersonInfo struct {
	GivenName   string    `validate:"person-name"`
	FamilyName  string    `validate:"person-name"`
	SSN         string    `validate:"required"`
	DateOfBirth date.Date `validate:"required,before"`
	Citizenship string    `validate:"required,iso3166_1_alpha2"`
}

func (pi *PersonInfo) Type() CustomerType {
	return Private
}

type OrganizationInfo struct {
	Name                string    `validate:"org-name"`
	Form                string    `validate:"required"`
	LeagalID            string    `validate:"required"`
	RegistrationDate    date.Date `validate:"required,before"`
	RegistrationCountry string    `validate:"required,iso3166_1_alpha2"`
}

func (oi *OrganizationInfo) Type() CustomerType {
	return Organization
}

type CustomerType int32

const (
	Private CustomerType = iota + 1
	Organization
)

type State int32

const (
	Prospect State = iota + 1
	Active
	Passive
)
