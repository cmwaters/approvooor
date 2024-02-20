package node

import (
	"context"
	"crypto/sha256"
	"encoding/json"
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

	startHash   = "0E005E02A1EE6F9350E2B74EF388925F65B8EC1E6D11E3DD92DDD20C82A860F8"
	startHeight = 1213740
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

	if !nodebuilder.IsInit(nodePath) {
		cfg := nodebuilder.DefaultConfig(nodeType)
		cfg.Header.TrustedHash = startHash
		cfg.DASer.SampleFrom = startHeight

		err = nodebuilder.Init(*cfg, nodePath, nodeType)
		if err != nil {
			return nil, err
		}
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

func (n *Node) Get(ctx context.Context, id ID) (SignedDocument, error) {
	signedDoc := SignedDocument{
		Document:   nil,
		Signatures: make([]Signature, 0),
	}
	namespace := id.Namespace()
	earliestHeight := id.Height()

	latestHeader, err := n.celnode.HeaderServ.NetworkHead(ctx)
	if err != nil {
		return signedDoc, err
	}

	for height := latestHeader.Height(); height >= earliestHeight; height-- {
		blobs, err := n.celnode.BlobServ.GetAll(ctx, height, []share.Namespace{namespace})
		if err != nil {
			return signedDoc, err
		}
		for _, blob := range blobs {
			var sigData Signature
			err := json.Unmarshal(blob.Data, &sigData)
			if err != nil {
				// assume it's not a signature
				signedDoc.Document = blob.Data
				continue
			}
			signedDoc.Signatures = append(signedDoc.Signatures, sigData)
		}

	}

	return signedDoc, nil
}

func (n *Node) getDocument(ctx context.Context, id ID) ([]byte, error) {
	blob, err := n.celnode.BlobServ.Get(ctx, id.Height(), id.Namespace(), id.Committment())
	if err != nil {
		return nil, err
	}
	return blob.Data, nil
}

func (n *Node) Sign(ctx context.Context, id ID) error {
	keys, err := n.signer.List()
	if err != nil {
		return err
	}

	file, err := n.getDocument(ctx, id)
	if err != nil {
		return err
	}

	signature, pubkey, err := n.signer.Sign(keys[0].Name, file)
	if err != nil {
		return err
	}

	data, err := json.Marshal(&Signature{
		signature, pubkey.Bytes(),
	})
	if err != nil {
		return err
	}

	_, err = n.Publish(ctx, data) // Maybe we wanna use that sigID?
	if err != nil {
		return err
	}

	return nil
}

type Signature struct {
	Signature []byte
	PubKey    []byte
}

type SignedDocument struct {
	Document   []byte
	Signatures []Signature
}
