package entities

type BaseEntity[T comparable] struct {
	id T
}

func (e *BaseEntity[T]) ID() T {
	return e.id
}

func (e *BaseEntity[T]) SetID(id T) {
	e.id = id
}

func (e *BaseEntity[T]) Equal(other *BaseEntity[T]) bool {
	if e == nil || other == nil {
		return false
	}

	var zero T
	if e.id == zero {
		return false // ID no inicializado.
	}

	return e.id == other.id
}
