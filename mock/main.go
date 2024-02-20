package main

import (
	"context"
	"fmt"

	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/share"
	"github.com/cmwaters/blobusign/node"
	"github.com/cmwaters/blobusign/server"
)

func main() {
	fmt.Println("Running mock server")
	n := &mockNode{}
	if err := server.Start(n); err != nil {
		fmt.Printf("ERR: %s", err.Error())
	}
}

var id node.ID

func init() {
	mockNamespace, err := share.NewBlobNamespaceV0([]byte("mock"))
	if err != nil {
		panic(err)
	}
	mockBlob, err := blob.NewBlobV0(mockNamespace, []byte("mockData"))
	if err != nil {
		panic(err)
	}
	id = node.NewID(100, mockNamespace, mockBlob.Commitment)
}

type mockNode struct{}

func (m *mockNode) Publish(ctx context.Context, data []byte) (node.ID, error) {
	// Mock implementation
	// In a real scenario, this would interact with a node to publish data
	fmt.Println("Received request to publish data", data)
	return []byte("mockID"), nil
}

func (m *mockNode) Get(ctx context.Context, id node.ID) (node.SignedDocument, error) {
	// Mock implementation
	// In a real scenario, this would retrieve data from a node using the provided ID
	fmt.Println("Received request to get data for ID", id)
	return node.SignedDocument{
		Document: []byte("mockData"),
		Signatures: []node.Signature{
			{
				Signature: []byte("mockSig"),
				PubKey:    []byte("mockPubKey"),
			},
		},
	}, nil
}

func (m *mockNode) Sign(ctx context.Context, id node.ID) error {
	// Mock implementation
	// In a real scenario, this would sign the data associated with the provided ID
	fmt.Println("Received request to sign data for ID", id)
	return nil
}
