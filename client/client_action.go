package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/newdag/ledger"
	nstore "github.com/newdag/store"
	wal "github.com/newdag/wallet"
)

func Init_ClientWallet() error {

	//	m_lastblock := bc.GetLastBlock()

	_, err := wal.GetWallets()
	return err
}

//call by client/http_server.go
func GetTransactions(index int) string { //index==-1 , return all trans
	bc := ledger.GetBlockChain(0)
	defer bc.Close()

	sTrans := ""

	bci := ledger.GetBlockchainIterator(bc)
	for {
		block := bci.Next()

		tm := time.Unix(block.Timestamp, 0)

		sTrans += fmt.Sprintf("Index:%d, time:%s, hash:%x\n", block.Index, tm.Format("2006-01-02 03:04:05 PM"), block.Hash)
		//sTrans += fmt.Sprintf("Prev. block:          %x\n", block.PrevBlockHash)
		sTrans += fmt.Sprintf("Block Data as string: %s\n", block.Data)
		//sTrans += fmt.Sprintf("Block MinerSignature: %x\n", block.MinerSignature)
		//sTrans += fmt.Sprintf("Block MinnerPubkey: %x\n", block.MinerPubkey)
		for _, tx := range block.Transactions {
			sTrans += fmt.Sprintf("%s", tx) //tx.String()
		}
		sTrans += fmt.Sprintf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return sTrans
}

//call by http_server.go  -> client err := m_proxy.SubmitTx(*tran)
func client_GetBalance(sAddr string) uint64 {
	return ledger.GetBalance(nil, sAddr)
}

//call by http_server.go  -> client err := m_proxy.SubmitTx(*tran)
func ConvertToTrans(from, to, sAmount, data string) (*ledger.Transaction, error) {
	//uint64
	amount, err := strconv.Atoi(sAmount) //strconv.ParseUint(sAmount, 10, 64) //??uint64 的问题待解

	if amount <= 0 || err != nil {
		return nil, fmt.Errorf("ERROR: amount is not valid or covert to int error!")
	}
	bc := ledger.GetBlockChain(0)
	defer bc.Close()
	if !wal.ValidateAddress(from) {
		return nil, fmt.Errorf("ERROR: Sender address is not valid")
	}
	if !wal.ValidateAddress(to) {
		return nil, fmt.Errorf("ERROR: Recipient address is not valid")
	}

	wallets, e := wal.GetWallets()
	if e != nil {
		return nil, e
	}
	wallet := wallets.GetWallet(from)
	if wallet == nil {
		return nil, fmt.Errorf("ERROR: Sender address is not in wallets")
	}

	var e1 error
	UTXOSet := ledger.UTXOSet{bc}
	tran, e1 := ledger.NewUTXOTransaction(wallet, to, amount, data, &UTXOSet)
	//tran := NewCoinbaseTX(to, data)
	if e1 != nil {
		fmt.Println("Error: NewUTXOTransaction return nil")
		return nil, e1
	}
	return tran, nil
}

//call by client:  client/main.go:  select { case commit := <-proxy.CommitCh():
func Write_ClientBlock(block ledger.Block, logger *nstore.Logger) error {

	bc := ledger.GetBlockChain(1)
	defer bc.Close()

	bc.AddBlock(&block)

	return nil
}
