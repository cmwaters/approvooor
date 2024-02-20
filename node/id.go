package node

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/celestiaorg/celestia-node/blob"
	ns "github.com/celestiaorg/celestia-node/share"
)

type ID []byte

const (
	CommitmentSize = 32
	HeightSize     = 8
	IDSize         = CommitmentSize + ns.NamespaceSize + HeightSize
)

// TODO: add constructor
func NewID(height uint64, namespace ns.Namespace, committment blob.Commitment) ID {
	heightBytes := make([]byte, HeightSize)
	binary.BigEndian.PutUint64(heightBytes, height)
	buf := bytes.NewBuffer(heightBytes)
	if _, err := buf.Write(namespace); err != nil {
		panic(err)
	}
	if _, err := buf.Write(committment); err != nil {
		panic(err)
	}
	return ID(buf.Bytes())
}

func Parse(id []byte) (ID, error) {
	if err := validateIDSize(id); err != nil {
		return nil, err
	}
	return ID(id), nil
}

func (id ID) Namespace() ns.Namespace {
	if err := validateIDSize(id); err != nil {
		panic(err)
	}
	return ns.Namespace(id[HeightSize : HeightSize+ns.NamespaceSize])
}

func (id ID) Height() uint64 {
	if err := validateIDSize(id); err != nil {
		panic(err)
	}
	return binary.BigEndian.Uint64(id[:HeightSize])
}

func (id ID) Committment() blob.Commitment {
	if err := validateIDSize(id); err != nil {
		panic(err)
	}
	return blob.Commitment(id[HeightSize+ns.NamespaceSize:])
}

func validateIDSize(id ID) error {
	if len(id) != HeightSize+ns.NamespaceSize+CommitmentSize {
		return fmt.Errorf("invalid ID length: expected %d, got %d", IDSize, len(id))
	}
	return nil
}
