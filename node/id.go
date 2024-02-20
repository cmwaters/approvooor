package node

type ID []byte

// TODO: add constructor

func (id ID) Namespace() []byte {
	panic("todo")
}

func (id ID) Height() uint64 {
	panic("todo")
}

func (id ID) Committment() []byte {
	panic("todo")
}

func Parse(id []byte) (ID, error) {
	panic("todo")
}

func NewID() (ID, error) {
	panic("todo")
}
