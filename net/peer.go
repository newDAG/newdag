package net

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"sync"

	"github.com/newdag/store"
)

type Peer struct {
	NetAddr   string
	PubKeyHex string
}

func (p *Peer) PubKeyBytes() ([]byte, error) {
	return hex.DecodeString(p.PubKeyHex[2:])
}

// PeerStore provides an interface for persistent storage and
// retrieval of peers.
type PeerStore interface {
	// Peers returns the list of known peers.
	Peers() ([]Peer, error)

	// SetPeers sets the list of known peers. This is invoked when a peer is
	// added or removed.
	SetPeers([]Peer) error
}

// StaticPeers is used to provide a static list of peers.
type StaticPeers struct {
	StaticPeers []Peer
	l           sync.Mutex
}

// Peers implements the PeerStore interface.
func (s *StaticPeers) Peers() ([]Peer, error) {
	s.l.Lock()
	peers := s.StaticPeers
	s.l.Unlock()
	return peers, nil
}

// SetPeers implements the PeerStore interface.
func (s *StaticPeers) SetPeers(p []Peer) error {
	s.l.Lock()
	s.StaticPeers = p
	s.l.Unlock()
	return nil
}

// JSONPeers is used to provide peer persistence on disk in the form
// of a JSON file. This allows human operators to manipulate the file.
type JSONPeers struct {
	l sync.Mutex
}

// NewJSONPeers creates a new JSONPeers store.
func NewJSONPeers() *JSONPeers {
	store := &JSONPeers{}
	return store
}

// Peers implements the PeerStore interface.
func (j *JSONPeers) Peers() ([]Peer, error) {
	j.l.Lock()
	defer j.l.Unlock()

	buf, err := store.ReadFile("peer")
	if err != nil || len(buf) == 0 {
		return nil, err
	}

	// Decode the peers
	var peerSet []Peer
	dec := json.NewDecoder(bytes.NewReader(buf))
	if err := dec.Decode(&peerSet); err != nil {
		return nil, err
	}

	return peerSet, nil
}

// SetPeers implements the PeerStore interface.
func (j *JSONPeers) SetPeers(peers []Peer) error {
	j.l.Lock()
	defer j.l.Unlock()

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(peers); err != nil {
		return err
	}

	return store.WriteFile("peer", buf.Bytes()) // Write out as JSON

}

// ExcludePeer is used to exclude a single peer from a list of peers.
func ExcludePeer(peers []Peer, peer string) (int, []Peer) {
	index := -1
	otherPeers := make([]Peer, 0, len(peers))
	for i, p := range peers {
		if p.NetAddr != peer {
			otherPeers = append(otherPeers, p)
		} else {
			index = i
		}
	}
	return index, otherPeers
}

//Sorting

// ByPubKey implements sort.Interface for []Peer based on
// the PubKeyHex field.
type ByPubKey []Peer

func (a ByPubKey) Len() int      { return len(a) }
func (a ByPubKey) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPubKey) Less(i, j int) bool {
	ai := a[i].PubKeyHex
	aj := a[j].PubKeyHex
	return ai < aj
}
