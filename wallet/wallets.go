package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/newdag/store"
)

// Wallets stores a collection of wallets
type Wallets struct {
	Wallets map[string]*Wallet
}

// GetWallets creates Wallets and fills it from a file if it exists
func GetWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile()

	return &wallets, err
}

// CreateWallet adds a Wallet to Wallets
func (ws *Wallets) CreateWallet() *Wallet {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet

	return wallet
}

// GetAddresses returns an array of addresses stored in the wallet file
func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetWallet returns a Wallet by its address
func (ws Wallets) GetWallet(address string) *Wallet {
	return ws.Wallets[address]
}

// LoadFromFile loads wallets from the file
func (ws *Wallets) LoadFromFile() error {
	fileContent, err := store.ReadFile("wallet")
	if err != nil {
		return err
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}

// SaveToFile saves wallets to a file
func (ws Wallets) SaveToFile(datadir string) {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	store.WriteFile("wallet", content.Bytes())
}

func (ws Wallets) listAddresses() {
	wallets, err := GetWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()
	for _, address := range addresses {
		fmt.Println(address)
	}
}
