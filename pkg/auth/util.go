package auth

func validate(o interface{}) error {
	if i, ok := o.(interface {
		Valid() error
	}); ok {
		return i.Valid()
	}

	return nil
}
