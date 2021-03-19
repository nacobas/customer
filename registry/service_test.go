package registry_test

import (
	"context"
	"testing"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/cockroachdb/errors"
	"github.com/nacobas/customer/customer"
	"github.com/nacobas/customer/registry"
	"github.com/nacobas/customer/repo/inmem"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	repo := inmem.NewRepo()

	svc := registry.NewService(repo)

	testCases := []struct {
		desc string
		info customer.Info
		want *customer.Customer
		err  error
	}{
		{
			desc: "new person",
			info: testPerson(t),
			want: &customer.Customer{State: 1, Info: testPerson(t)},
			err:  nil,
		},
		{
			desc: "new org",
			info: testOrg(t),
			want: &customer.Customer{State: 1, Info: testOrg(t)},
			err:  nil,
		},
		{
			desc: "invalid person info - missing SSN",
			info: &customer.PersonInfo{
				GivenName:   "given-name",
				FamilyName:  "family-name",
				SSN:         "",
				DateOfBirth: parseDate(t, "1970-01-01"),
				Citizenship: "US",
			},
			want: nil,
			err:  registry.ErrValidation,
		},
		{
			desc: "invalid organisation info - invalid Country",
			info: &customer.OrganizationInfo{
				Name:                "org-name",
				Form:                "Ltd",
				LeagalID:            "legal-id",
				RegistrationDate:    parseDate(t, "1970-01-01"),
				RegistrationCountry: "EUR",
			},
			want: nil,
			err:  registry.ErrValidation,
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			got, err := svc.New(context.Background(), tC.info)

			if tC.err == nil {
				assert.NotZero(t, got.ID, "New ID shoul not be zero")
				assert.Equal(t, tC.want.State, got.State, "customer state should equal")
				assert.Equal(t, tC.want.Info, got.Info, "customer.Info should be equal")
				assert.Nil(t, err, "error should be nil")
			} else {
				assert.True(t, errors.Is(err, tC.err), "Expected error should be found in the chain")
				assert.Nil(t, got, "customer should be nil")
			}

		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	repo := inmem.NewRepoWithSeed(seed(t))

	svc := registry.NewService(repo)

	testCases := []struct {
		desc string
		id   uint32
		want *customer.Customer
		err  error
	}{
		{
			desc: "get person",
			id:   1,
			want: &customer.Customer{ID: 1, State: 1, Info: testPerson(t)},
			err:  nil,
		},
		{
			desc: "get org",
			id:   2,
			want: &customer.Customer{ID: 2, State: 2, Info: testOrg(t)},
			err:  nil,
		},
		{
			desc: "not found",
			id:   3,
			want: nil,
			err:  registry.ErrNotFound,
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			got, err := svc.Get(context.Background(), tC.id)

			if tC.err == nil {
				assert.Equal(t, tC.want, got, "customer should equal")
				assert.Nil(t, err, "error should be nil")
			} else {
				assert.True(t, errors.Is(err, tC.err), "Expected error should be found in the chain")
				assert.Nil(t, got, "customer should be nil")
			}

		})
	}
}

func TestUpdateInfo(t *testing.T) {
	t.Parallel()

	repo := inmem.NewRepoWithSeed(seed(t))

	svc := registry.NewService(repo)

	testCases := []struct {
		desc string
		id   uint32
		info customer.Info
		want *customer.Customer
		err  error
	}{
		{
			desc: "update person",
			id:   1,
			info: &customer.PersonInfo{
				GivenName:   "new-given-name",
				FamilyName:  "family-name",
				SSN:         "SSN",
				DateOfBirth: parseDate(t, "1970-01-01"),
				Citizenship: "US"},
			want: &customer.Customer{ID: 1, State: 1, Info: &customer.PersonInfo{
				GivenName:   "new-given-name",
				FamilyName:  "family-name",
				SSN:         "SSN",
				DateOfBirth: parseDate(t, "1970-01-01"),
				Citizenship: "US"}},
			err: nil,
		},
		{
			desc: "update org",
			id:   2,
			info: &customer.OrganizationInfo{
				Name:                "new-org-name",
				Form:                "Ltd",
				LeagalID:            "legal-id",
				RegistrationDate:    parseDate(t, "1970-01-01"),
				RegistrationCountry: "US"},
			want: &customer.Customer{ID: 2, State: 2, Info: &customer.OrganizationInfo{
				Name:                "new-org-name",
				Form:                "Ltd",
				LeagalID:            "legal-id",
				RegistrationDate:    parseDate(t, "1970-01-01"),
				RegistrationCountry: "US"}},
			err: nil,
		},
		{
			desc: "not found",
			id:   3,
			info: testPerson(t),
			want: nil,
			err:  registry.ErrNotFound,
		},
		{
			desc: "customer info validation",
			id:   2,
			info: &customer.OrganizationInfo{
				Name:                "org-name",
				Form:                "Ltd",
				LeagalID:            "",
				RegistrationDate:    parseDate(t, "1970-01-01"),
				RegistrationCountry: "US",
			},
			want: nil,
			err:  registry.ErrValidation,
		},
		{
			desc: "wrong type of customer info",
			id:   1,
			info: testOrg(t),
			want: nil,
			err:  registry.ErrExpected,
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			got, err := svc.UpdateInfo(context.Background(), tC.id, tC.info)

			if tC.err == nil {
				assert.Equal(t, tC.want, got, "customer should equal")
				assert.Nil(t, err, "error should be nil")
			} else {
				assert.True(t, errors.Is(err, tC.err), "Expected error should be found in the chain")
				assert.Nil(t, got, "customer should be nil")
			}

		})
	}
}

func TestSetState(t *testing.T) {
	t.Parallel()

	repo := inmem.NewRepoWithSeed(seed(t))

	svc := registry.NewService(repo)

	testCases := []struct {
		desc  string
		id    uint32
		state customer.State
		err   error
	}{
		{
			desc:  "set person to active",
			id:    1,
			state: 2,
			err:   nil,
		},
		{
			desc:  "set org to passive",
			id:    2,
			state: 3,
			err:   nil,
		},
		{
			desc:  "not found",
			id:    3,
			state: 1,
			err:   registry.ErrNotFound,
		},
		{
			desc:  "invalid state value",
			id:    1,
			state: 0,
			err:   registry.ErrValidation,
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			err := svc.SetState(context.Background(), tC.id, tC.state)

			assert.Truef(t, errors.Is(err, tC.err), "Expecting %v , got: %v", tC.err, err)

		})
	}
}

func seed(t *testing.T) []customer.Customer {

	return []customer.Customer{
		{ID: 1, State: 1, Info: testPerson(t)},
		{ID: 2, State: 2, Info: testOrg(t)},
	}

}

func testPerson(t *testing.T) *customer.PersonInfo {
	return &customer.PersonInfo{
		GivenName:   "given-name",
		FamilyName:  "family-name",
		SSN:         "SSN",
		DateOfBirth: parseDate(t, "1970-01-01"),
		Citizenship: "US"}
}

func testOrg(t *testing.T) *customer.OrganizationInfo {
	return &customer.OrganizationInfo{
		Name:                "org-name",
		Form:                "Ltd",
		LeagalID:            "legal-id",
		RegistrationDate:    parseDate(t, "1970-01-01"),
		RegistrationCountry: "US"}
}

func parseDate(t *testing.T, datestr string) date.Date {
	d, err := date.ParseDate(datestr)
	if err != nil {
		t.Fatalf("Failed to parse date from: %s, error: %v", datestr, err)
	}
	return d
}
