package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "tea",
		Price: 2.40,
		SKU:   "abs-sba-bas",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
