package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"sync"

	"github.com/newdag/store"
)

type PemKey struct {
	l sync.Mutex
}

func NewPemKey() *PemKey {
	pemKey := &PemKey{}
	return pemKey
}

func (k *PemKey) ReadKey() (*ecdsa.PrivateKey, error) {
	k.l.Lock()
	defer k.l.Unlock()

	buf, err := store.ReadFile("peer_key")
	if err != nil {
		return nil, err
	}

	return k.ReadKeyFromBuf(buf)
}

func (k *PemKey) ReadKeyFromBuf(buf []byte) (*ecdsa.PrivateKey, error) {
	if len(buf) == 0 {
		return nil, nil
	}

	block, _ := pem.Decode(buf)
	if block == nil {
		return nil, fmt.Errorf("Error decoding PEM block from data")
	}
	return x509.ParseECPrivateKey(block.Bytes)
}

func (k *PemKey) WriteKey(key *ecdsa.PrivateKey) error {
	k.l.Lock()
	defer k.l.Unlock()

	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}
	pemBlock := &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	data := pem.EncodeToMemory(pemBlock)
	return store.WriteFile("peer_key", data)
}

type PemDump struct {
	PublicKey  string
	PrivateKey string
}

func GeneratePemKey() (*PemDump, error) {
	key, err := GenerateECDSAKey()
	if err != nil {
		return nil, err
	}

	pub := fmt.Sprintf("0x%X", FromECDSAPub(&key.PublicKey))

	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, err
	}
	pemBlock := &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	data := pem.EncodeToMemory(pemBlock)

	pemDump := PemDump{
		PublicKey:  pub,
		PrivateKey: string(data),
	}

	return &pemDump, err
}
