package ccvalidator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var mockNow = func() time.Time {
	return time.Date(2024, 4, 1, 1, 1, 1, 0, time.UTC)
}

func Test_luhnValidate(t *testing.T) {
	type args struct {
		number         string
		year           int
		month          int
		allowTestCards bool
	}
	testCases := []struct {
		name  string
		args  args
		valid bool
		err   error
	}{
		{
			name: "positive_amex",
			args: args{
				number: "370193274384838",
				month:  int(mockNow().Month()),
				year:   mockNow().Year(),
			},
			valid: true,
		},
		{
			name: "positive_atm",
			args: args{
				number: "6120 1007 6595 6750",
				month:  int(mockNow().Month()),
				year:   mockNow().Year(),
			},
			valid: true,
		},
		{
			name: "positive_visa_1",
			args: args{
				number: "4004669383770674",
				month:  int(mockNow().Month()),
				year:   mockNow().Year(),
			},
			valid: true,
		},
		{
			name: "positive_visa_2",
			args: args{
				number: "4004666596547",
				month:  int(mockNow().Month()),
				year:   mockNow().Year(),
			},
			valid: true,
		},
		{
			name: "positive_mastercard_1",
			args: args{
				number: "5513440000154150",
				month:  int(mockNow().Month()),
				year:   mockNow().Year(),
			},
			valid: true,
		},
		{
			name: "positive_mastercard_2",
			args: args{
				number: "5513440000456183",
				month:  int(mockNow().Month()),
				year:   mockNow().Year(),
			},
			valid: true,
		},
		{
			name: "negative_invalid_expr_year",
			args: args{
				number: "5513440000456183",
				month:  int(mockNow().Month()),
				year:   9999,
			},
			valid: false,
			err:   ErrInvalidExpYear,
		},
		{
			name: "negative_invalid_expr_month",
			args: args{
				number: "5513440000456183",
				month:  13,
				year:   mockNow().Year(),
			},
			valid: false,
			err:   ErrInvalidExpMonth,
		},
		{
			name: "negative_invalid_card_number",
			args: args{
				number: "551344000045618322",
				month:  int(mockNow().Month()),
				year:   mockNow().Year(),
			},
			valid: false,
			err:   ErrInvalidCardNumber,
		},
		{
			name: "negative_expired",
			args: args{
				number: "5513440000456183",
				month:  int(mockNow().Month()) - 1,
				year:   mockNow().Year() - 1,
			},
			valid: false,
			err:   ErrCardExpired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewLuhnValidator(tc.args.allowTestCards).
				Validate(tc.args.number, 111, tc.args.month, tc.args.year)
			if tc.err != nil && err != nil {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
