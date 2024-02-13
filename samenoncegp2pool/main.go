package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/giskook/ethtools/utils/wallet"
)

const (
	// 	URL           = "https://rpc.mxc.com"
	URL           = "http://44.232.123.15:8123"
	PrivateKey    = ""
	TransferCount = 10
	gasPrice      = 1000000000
)

func main() {
	client, err := ethclient.Dial(URL)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	log.Println("time start", time.Now())

	f := func(nonce uint64, toAddr common.Address) {
		tx := types.NewTx(&types.LegacyTx{
			Nonce: nonce,
			// GasPrice: big.NewInt(1000000000),
			GasPrice: big.NewInt(gasPrice),
			// GasPrice: gasPrice,
			Gas:   30000,
			To:    &toAddr,
			Value: big.NewInt(1),
		})

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Fatal(err)
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			log.Println(err, signedTx.Hash().Hex())
			return
		}

		log.Printf("nonce %d, tx sent: %s\n", nonce, signedTx.Hash().Hex()) // tx sent: 0xa56316b637a94c4cc0331c73ef26389d6c097506d581073f927275e7a6ece0bc
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	// nonce, err := client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

	addr1 := wallet.GetRandomWallet().Address
	addr2 := wallet.GetRandomWallet().Address
	//	for i := 0; i <= TransferCount; i++ {
	f(nonce, addr1)
	f(nonce, addr2)
	f(nonce, addr1)
	//	}

	time.Sleep(10 * time.Second)
	log.Println("time end", time.Now())
}
