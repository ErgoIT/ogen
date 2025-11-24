package validate

import (
	"testing"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDecimal_Set(t *testing.T) {
	assert.False(t, Decimal{}.Set())
	assert.True(t, Decimal{MaxSet: true}.Set())
	assert.True(t, Decimal{MinSet: true}.Set())
	assert.True(t, Decimal{MultipleOfSet: true}.Set())
}

func TestDecimal_Setters(t *testing.T) {
	for _, tc := range []struct {
		do       func(*Decimal)
		expected Decimal
	}{
		{
			do: func(i *Decimal) {
				i.SetMultipleOf(decimal.Ten)
			},
			expected: Decimal{
				MultipleOf:    decimal.Ten,
				MultipleOfSet: true,
			},
		},
		{
			do: func(i *Decimal) {
				i.SetExclusiveMaximum(decimal.Ten)
			},
			expected: Decimal{
				Max:          decimal.Ten,
				MaxExclusive: true,
				MaxSet:       true,
			},
		},
		{
			do: func(i *Decimal) {
				i.SetExclusiveMinimum(decimal.Ten)
			},
			expected: Decimal{
				Min:          decimal.Ten,
				MinExclusive: true,
				MinSet:       true,
			},
		},
		{
			do: func(i *Decimal) {
				i.SetMaximum(decimal.Ten)
			},
			expected: Decimal{
				Max:    decimal.Ten,
				MaxSet: true,
			},
		},
		{
			do: func(i *Decimal) {
				i.SetMinimum(decimal.Ten)
			},
			expected: Decimal{
				Min:    decimal.Ten,
				MinSet: true,
			},
		},
	} {
		var r Decimal
		tc.do(&r)
		assert.Equal(t, tc.expected, r)
	}
}

func TestDecimal_Validate(t *testing.T) {
	for _, tc := range []struct {
		Name      string
		Validator Decimal
		Value     decimal.Decimal
		Valid     bool
	}{
		{Name: "Zero", Valid: true},
		{
			Name:      "MaxOk",
			Validator: Decimal{Max: decimal.Ten, MaxSet: true},
			Value:     decimal.Ten,
			Valid:     true,
		},
		{
			Name:      "MaxErr",
			Validator: Decimal{Max: decimal.Ten, MaxSet: true},
			Value:     decimal.MustNew(11, 0),
			Valid:     false,
		},
		{
			Name:      "MaxExclErr",
			Validator: Decimal{Max: decimal.Ten, MaxSet: true, MaxExclusive: true},
			Value:     decimal.Ten,
			Valid:     false,
		},
		{
			Name:      "MinOk",
			Validator: Decimal{Min: decimal.Ten, MinSet: true},
			Value:     decimal.Ten,
			Valid:     true,
		},
		{
			Name:      "MinErr",
			Validator: Decimal{Min: decimal.Ten, MinSet: true},
			Value:     decimal.MustNew(9, 0),
			Valid:     false,
		},
		{
			Name:      "MinExclErr",
			Validator: Decimal{Min: decimal.Ten, MinSet: true, MinExclusive: true},
			Value:     decimal.Ten,
			Valid:     false,
		},
		{
			Name:      "MultipleOfOk",
			Validator: Decimal{MultipleOf: decimal.Ten, MultipleOfSet: true},
			Value:     decimal.MustNew(20, 0),
			Valid:     true,
		},
		{
			Name:      "MultipleOfErr",
			Validator: Decimal{MultipleOf: decimal.Ten, MultipleOfSet: true},
			Value:     decimal.MustNew(13, 0),
			Valid:     false,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			valid := tc.Validator.Validate(tc.Value) == nil
			assert.Equal(t, tc.Valid, valid, "%v: %+v",
				tc.Validator,
				tc.Value,
			)
		})
	}
}
