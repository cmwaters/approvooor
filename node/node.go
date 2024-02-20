package node

import "context"

type Node struct {
}

func (n *Node) Publish(ctx context.Context, data []byte) (ID, error) {
	return nil, nil
}

func (n *Node) Get(ctx context.Context, id ID) ([]byte, error) {
	return nil, nil
}

func (n *Node) Sign(ctx context.Context, id ID) error {
	return nil
}
