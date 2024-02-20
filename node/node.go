package node

import (
	"context"
	"crypto/sha256"
	"os"
	"path/filepath"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	"github.com/celestiaorg/celestia-node/logs"
	"github.com/celestiaorg/celestia-node/nodebuilder"
	"github.com/celestiaorg/celestia-node/nodebuilder/node"
	"github.com/celestiaorg/celestia-node/nodebuilder/p2p"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	appns "github.com/celestiaorg/celestia-app/pkg/namespace"
	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/share"
	"github.com/ipfs/go-log/v2"
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
	logs.SetAllLoggers(log.LevelInfo)
	keysPath := filepath.Join(nodePath, "keys")
	encConf := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	signer, err := keyring.New(app.Name, keyring.BackendTest, keysPath, os.Stdin, encConf.Codec)
	if err != nil {
		return nil, err
	}

	err = nodebuilder.Init(*nodebuilder.DefaultConfig(nodeType), nodePath, nodeType)
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

func (n *Node) Start(ctx context.Context) error {
	return n.celnode.Start(ctx)
}

func (n *Node) Stop(ctx context.Context) error {
	return n.celnode.Stop(ctx)
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
	namespace := id.Namespace()
	earliestHeight := id.Height()

	latestHeader, err := n.celnode.HeaderServ.NetworkHead(ctx)
	if err != nil {
		return nil, err
	}

	for height := latestHeader.Height(); height >= earliestHeight; height-- {
		_, err := n.celnode.BlobServ.GetAll(ctx, height, []share.Namespace{namespace})
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (n *Node) getDocument(ctx context.Context, id ID) ([]byte, error) {
	blob, err := n.celnode.BlobServ.Get(ctx, id.Height(), id.Namespace(), id.Committment())
	if err != nil {
		return nil, err
	}
	return blob.Data, nil
}

func (n *Node) Sign(ctx context.Context, id ID) error {
	return nil
}
