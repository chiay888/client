package sigchain3

import (
	keybase1 "github.com/keybase/client/go/protocol/keybase1"
)

type UID [16]byte
type LinkType int
type ChainType int
type SigVersion int
type LinkID []byte
type Seqno = keybase1.Seqno
type Time = keybase1.Time
type SigIgnoreIfUnsupported bool
type KID []byte
type TeamID = keybase1.TeamID
type PerTeamKeyGeneration = keybase1.PerTeamKeyGeneration

// These values are picked so they don't conflict with Sigchain V1 and V2 link types
const (
	LinkTypeNone              LinkType = 0
	LinkTypeUserSecretSummary LinkType = 65
	LinkTypeUserPassiveFollow LinkType = 66
	LinkTypeTeamPerTeamKey    LinkType = 81
)

// The values are picked so they don't conflict with Sigchain V1 and V2 SeqType's
const (
	ChainTypeUserPrivateOffTree ChainType = 16
	ChainTypeTeamPrivateOffTree ChainType = 17
)

// OuterLink V3 is the third version of Keybase sigchain signatures, it roughly approximates
// the outer link v2s that we have previously used.
type OuterLink struct {
	_struct             bool                   `codec:",toarray"`
	Version             SigVersion             `codec:"version"` // comment should be 3
	Seqno               Seqno                  `codec:"seqno"`
	Prev                LinkID                 `codec:"prev"`
	Curr                LinkID                 `codec:"curr"`
	LinkType            LinkType               `codec:"type"`
	ChainType           ChainType              `codec:"chaintype"`
	IgnoreIfUnsupported SigIgnoreIfUnsupported `codec:"ignore_if_unsupported"`
	// New field for V3; if this link is encrypted, specify the format, nonce and PUK
	EncParams *EncryptionParameters `codec:"encryption_parameters"`
}

type InnerLink struct {
	Signer     Signer      `codec:"s"` // Info on the signer, including UID, KID and eldest
	TeamID     *TeamID     `codec:"t"` // for teams, the TeamID, and null otherwise
	Ctime      Time        `codec:"c"` // Seconds since 1970 UTC.
	MerkleRoot *MerkleRoot `codec:"m"` // Optional snapshot of merkle root at time of sig
	ClientInfo *ClientInfo `codec:"i"` // Optional client type making sig
	Body       interface{} `codec:"b"` // The actual body, which varies based on the type in the outer link
}

type Signer struct {
	UID         UID            `codec:"u"`
	EldestSeqno keybase1.Seqno `codec:"e"`
	KID         KID            `codec:"k"`
}

type PassiveFollowBody struct {
	Follows map[UID]Seqno `codec:"f"`
}

type SecretSummaryBody struct {
	Follows map[UID]Seqno `codec:"f"`
}

type PerTeamKeyBody struct {
	Generation    PerTeamKeyGeneration `codec:"g"`
	SigningKID    KID                  `codec:"s"`
	EncryptionKID KID                  `codec:"e"`
	ReverseSig    []byte               `codec:"r"`
}

type MerkleRoot struct {
	Hash  []byte `codec:"h"` // HashMeta of the MerkleRoot
	Seqno Seqno  `codec:"s"`
	Ctime Time   `codec:"c"`
}

type ClientInfo struct {
	Desc    string `codec:"d"`
	Version string `codec:"v"`
}

// If the inner link is encrypted, we specify the encryption parameters
// with this offloaded structure. So far, we don't know of any such encrypted
// payloads, but we'll allow it.
type EncryptionParameters struct {
	Version int    `codec:"v"`
	KID     KID    `codec:"k"`
	Nonce   []byte `codec:"n"`
}

type Tail struct {
	_struct   bool      `codec:",toarray"`
	ChainType ChainType `codec:"seqtype"`
	Seqno     Seqno     `codec:"seqno"`
	Hash      LinkID    `codec:"hash"` // hash of the outer link
}