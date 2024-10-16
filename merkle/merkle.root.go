package merkle

import (
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/types"
)

type IRoot[TI index.Value, TH types.HashType] interface {
	Hash() TH
	ValidateProof(proof *Proof[TI, TH]) bool
}

type root[TI index.Value, TH types.HashType] struct {
	hf   types.Hasher[TH]
	hash TH
}

func newRoot[TI index.Value, TH types.HashType](hash TH, hf types.Hasher[TH]) IRoot[TI, TH] {
	return &root[TI, TH]{hf: hf, hash: hash}
}

func (r *root[TI, TH]) Hash() TH {
	return r.hash
}

func (r *root[TI, TH]) ValidateProof(proof *Proof[TI, TH]) bool {
	if len(proof.Hashes) == 0 {
		return false
	}
	hashesToProof := proof.RightPeaks

	currentIndex := index.LeafIndex[TI](proof.Target)
	currentHash := proof.Hashes[0]
	for _, siblingHash := range proof.Hashes[1:] {
		upper := currentIndex.Up()
		var currentNode INode[TH]
		if currentIndex.IsRight() {
			currentNode = Node[TH](siblingHash, currentHash)
		} else {
			currentNode = Node[TH](currentHash, siblingHash)
		}

		if err := buildNodeHash(r.hf, currentNode, func(nodeHash TH) error {
			currentHash = nodeHash
			return nil
		}); err != nil {
			return false
		}
		currentIndex = upper
	}
	hashesToProof = append(hashesToProof, currentHash)

	hashesToProof = append(hashesToProof, proof.LeftPeaks...)
	hashBytes := make([][]byte, 0, len(hashesToProof))
	for _, h := range hashesToProof {
		data, err := types.HashBytes[TH](h)
		if err != nil {
			return false
		}
		hashBytes = append(hashBytes, data)
	}

	calculatedHash := r.hf(hashBytes...)
	return calculatedHash == r.hash
}
