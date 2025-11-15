package base

type Base struct {
	id int
}

func NewBase(id int) *Base {
	return &Base{id: id}
}

func (b *Base) GetID() int {
	return b.id
}
