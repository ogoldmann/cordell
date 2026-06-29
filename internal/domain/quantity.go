package domain

// Quantity represents a positive item quantity in the custody domain.
type Quantity int

// NewQuantity creates a valid Quantity from an integer value.
func NewQuantity(value int) (Quantity, error) {
	if value <= 0 {
		return 0, ErrInvalidQuantity
	}

	return Quantity(value), nil
}

// Int returns the quantity as a regular integer.
func (q Quantity) Int() int {
	return int(q)
}
