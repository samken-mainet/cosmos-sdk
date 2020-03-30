package types

import (
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

// ToTmTypes casts a proto ValidatorSet to tendendermint type.
func (valset ValidatorSet) ToTmTypes() *tmtypes.ValidatorSet {
	vs := tmtypes.ValidatorSet{
		Validators: valset.Validators,
		Proposer:   valset.Proposer,
	}
	_ = vs.TotalVotingPower()
	return &vs
}

// ToTmTypes casts a proto ValidatorSet to tendendermint type.
func ValSetFromTmTypes(valset tmtypes.ValidatorSet) ValidatorSet {
	return ValidatorSet{
		Validators:       valset.Validators,
		Proposer:         valset.Proposer,
		TotalVotingPower: valset.TotalVotingPower(),
	}
}

// ToTmTypes casts a proto SignedHeader to tendendermint type.
func (sh SignedHeader) ToTmTypes() *tmtypes.SignedHeader {
	tmHeader := &tmtypes.Header{
		Version: version.Consensus{
			Block: version.Protocol(sh.Header.Version.Block),
			App:   version.Protocol(sh.Header.Version.App),
		},
		ChainID:            sh.Header.ChainID,
		Height:             sh.Header.Height,
		Time:               sh.Header.Time,
		LastBlockID:        tmtypes.TM2PB.BlockID(sh.Header.LastBlockId),
		LastCommitHash:     sh.Header.LastCommitHash,
		DataHash:           sh.Header.DataHash,
		ValidatorsHash:     sh.Header.ValidatorsHash,
		NextValidatorsHash: sh.Header.NextValidatorsHash,
		ConsensusHash:      sh.Header.ConsensusHash,
		AppHash:            sh.Header.AppHash,
		LastResultsHash:    sh.Header.LastResultsHash,
		EvidenceHash:       sh.Header.EvidenceHash,
		ProposerAddress:    sh.Header.ProposerAddress,
	}

	return &tmtypes.SignedHeader{
		Header: tmHeader,
		Commit: sh.Commit.ToTmTypes(),
	}
}

// ToTmTypes casts a proto ValidatorSet to tendendermint type.
func SignedHeaderFromTmTypes(sh *tmtypes.SignedHeader) *SignedHeader {
	abciHeader := tmtypes.TM2PB.Header(sh.Header)
	return &SignedHeader{
		Header: &abciHeader,
		Commit: CommitFromTmTypes(sh.Commit),
	}
}

func (bid BlockID) ToTmTypes() tmtypes.BlockID {
	return tmtypes.BlockID{
		Hash:        bid.Hash,
		PartsHeader: bid.PartsHeader.ToTmTypes(),
	}
}

// ToTmTypes casts a proto ValidatorSet to tendendermint type.
func BlockIDFromTmTypes(bid tmtypes.BlockID) *BlockID {
	return &BlockID{
		Hash:        bid.Hash,
		PartsHeader: PartSetHeaderFromTmTypes(bid.PartsHeader),
	}
}

func (ph *PartSetHeader) ToTmTypes() tmtypes.PartSetHeader {
	return tmtypes.PartSetHeader{
		Total: int(ph.Total),
		Hash:  ph.Hash,
	}
}

// ToTmTypes casts a proto ValidatorSet to tendendermint type.
func PartSetHeaderFromTmTypes(ph tmtypes.PartSetHeader) *PartSetHeader {
	return &PartSetHeader{
		Total: int32(ph.Total),
		Hash:  ph.Hash,
	}
}

// ToTmTypes casts a proto ToTmTypes to tendendermint type.
func (c Commit) ToTmTypes() *tmtypes.Commit {
	tmCommit := &tmtypes.Commit{
		Height:     c.Height,
		Round:      int(c.Round),
		BlockID:    c.BlockID.ToTmTypes(),
		Signatures: c.Signatures.ToTmTypes(),
	}
	_ = tmCommit.Hash()
	_ = tmCommit.BitArray()
	return tmCommit
}

// ToTmTypes casts a proto ValidatorSet to tendendermint type.
func CommitFromTmTypes(c *tmtypes.Commit) Commit {
	commitSigs := make([]*CommitSig, len(c.Signatures))

	for i := range c.Signatures {
		cs := CommitSig{
			BlockIDFlag:      []byte{byte(c.Signatures[i].BlockIDFlag)},
			ValidatorAddress: c.Signatures[i].ValidatorAddress,
			Timestamp:        c.Signatures[i].Timestamp,
			Signature:        c.Signatures[i].Signature,
		}
		commitSigs[i] = &cs
	}

	return Commit{
		Height:     c.Height,
		Round:      int32(c.Round),
		BlockID:    BlockIDFromTmTypes(c.BlockID),
		Signatures: commitSigs,
		hash:       c.Hash(),
		bitArray:   c.BitArray().Bytes(),
	}
}