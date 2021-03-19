package customer

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewWithRandom(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc string
		info Info
		want *Customer
		err  error
	}{
		{
			desc: "new person",
			info: testPerson(t),
			want: &Customer{State: 1, Info: testPerson(t)},
			err:  nil,
		},
		{
			desc: "new org",
			info: testOrg(t),
			want: &Customer{State: 1, Info: testOrg(t)},
			err:  nil,
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			got := NewWithRandomID(tC.info)

			assert.Equal(t, tC.want.State, got.State)
			assert.Equal(t, tC.want.Info, got.Info)
			assert.Equal(t, tC.want.Type(), got.Type())
			assert.NotZero(t, got.ID)

		})
	}
}

func TestUpdateInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc string
		c    *Customer
		info Info
		err  error
	}{
		{
			desc: "update person info",
			c:    &Customer{State: 1, Info: testPerson(t)},
			info: &PersonInfo{"new-name", "new-name", "new-SSN", parseDate(t, "2020-01-01"), "FI"},
			err:  nil,
		},
		{
			desc: "update org info",
			c:    &Customer{State: 1, Info: testOrg(t)},
			info: &OrganizationInfo{"new-name", "new", "new-id", parseDate(t, "2020-01-01"), "FI"},
			err:  nil,
		},
		{
			desc: "try update person with org",
			c:    &Customer{State: 1, Info: testPerson(t)},
			info: &OrganizationInfo{"new-name", "new", "new-id", parseDate(t, "2020-01-01"), "FI"},
			err:  ErrTypeNotEqual,
		},
		{
			desc: "try org with person",
			c:    &Customer{State: 1, Info: testOrg(t)},
			info: &PersonInfo{"new-name", "new-name", "new-SSN", parseDate(t, "2020-01-01"), "FI"},
			err:  ErrTypeNotEqual,
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			err := tC.c.UpdateInfo(tC.info)

			if tC.err == nil {
				assert.Equal(t, tC.c.Info, tC.info)
			} else {
				assert.True(t, errors.Is(err, tC.err), "Expected error should be found in the chain")
			}
		})
	}
}

func TestCustomerValidations(t *testing.T) {
	t.Parallel()

	v := NewValidator()

	testCases := []struct {
		desc string
		c    *Customer
		err  error
	}{
		{
			desc: "valid private customer",
			c:    &Customer{ID: 1, State: 1, Info: testPerson(t)},
			err:  nil,
		},
		{
			desc: "valid org customer",
			c:    &Customer{ID: 1, State: 1, Info: testOrg(t)},
			err:  nil,
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {

			got := v.Struct(tC.c)

			assert.Equal(t, tC.err, got)

		})
	}
}

func TestPersonValidations(t *testing.T) {
	t.Parallel()

	v := NewValidator()

	testCases := []struct {
		desc string
		c    *PersonInfo
		err  error
	}{
		{
			desc: "valid person info",
			c:    testPerson(t),
			err:  nil,
		},
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.desc, func(t *testing.T) {

			got := v.Struct(tC.c)

			assert.Equal(t, tC.err, got)

		})
	}
}

func testPerson(t *testing.T) *PersonInfo {
	return &PersonInfo{"given-name", "family-name", "SSN", parseDate(t, "1970-01-01"), "US"}
}

func testOrg(t *testing.T) *OrganizationInfo {
	return &OrganizationInfo{"org-name", "Ltd", "legal-id", parseDate(t, "1970-01-01"), "US"}
}

func parseDate(t *testing.T, datestr string) date.Date {
	d, err := date.ParseDate(datestr)
	if err != nil {
		t.Fatalf("Failed to parse date from: %s, error: %v", datestr, err)
	}
	return d
}
