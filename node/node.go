package node

import (
	"context"
	"crypto/sha256"
	"os"
	"path/filepath"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	"github.com/celestiaorg/celestia-node/nodebuilder"
	"github.com/celestiaorg/celestia-node/nodebuilder/node"
	"github.com/celestiaorg/celestia-node/nodebuilder/p2p"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	"github.com/celestiaorg/celestia-node/share"
	appns "github.com/celestiaorg/celestia-app/pkg/namespace"
	"github.com/celestiaorg/celestia-node/blob"
)

const (
	nodePath    = "~/.blobusign"
	nodeType    = node.Light
	nodeNetwork = p2p.Mocha
)

type Node struct {
	celnode *nodebuilder.Node
	signer  keyring.Keyring
}

func NewNode() (*Node, error) {
	keysPath := filepath.Join(nodePath, "keys")
	encConf := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	signer, err := keyring.New(app.Name, keyring.BackendTest, keysPath, os.Stdin, encConf.Codec)
	if err != nil {
		return nil, err
	}

	store, err := nodebuilder.OpenStore(nodePath, signer)
	if err != nil {
		return nil, err
	}

	node, err := nodebuilder.New(nodeType, nodeNetwork, store)
	if err != nil {
		return nil, err
	}

	return &Node{celnode: node, signer: signer}, nil
}

func (n *Node) Publish(ctx context.Context, data []byte) (ID, error) {
	hash := sha256.Sum256(data)
	ns, err := share.NewBlobNamespaceV0(hash[:appns.NamespaceVersionZeroIDSize])
	if err != nil {
		return nil, err
	}

	b, err := blob.NewBlobV0(ns, data)
	if err != nil {
		return nil, err
	}


	height, err := n.celnode.BlobServ.Submit(ctx, []*blob.Blob{b}, blob.DefaultGasPrice())
	if err != nil {
		return nil, err
	}

	return NewID(height, ns, b.Commitment), nil
}

func (n *Node) Get(ctx context.Context, id ID) ([]byte, error) {
	return nil, nil
}

func (n *Node) Sign(ctx context.Context, id ID) error {
	return nil
}
