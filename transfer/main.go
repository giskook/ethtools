package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"log"
	"math/big"
)

const (
	// URL = "https://endpoints.omniatech.io/v1/op/mainnet/public"
	URL = "https://rpc.linea.build"
	// URL = "https://optimism-mainnet.infura.io"
	//URL        = "https://okbtestrpc.okbchain.org"
	PrivateKey = ""
	// ChainID    = 10
	ChainID = 59144
)

func main() {
	chainID := big.NewInt(ChainID)
	client, err := ethclient.Dial(URL)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("suggest %v\n", gasPrice)

	privateKeyHex := PrivateKey

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	toAddr := common.HexToAddress("54DF7eF5CB805Fea7B59c580Ab4663680fEBbF98")
	tx := types.NewTx(&types.LegacyTx{
		Nonce: nonce,
		//GasPrice: big.NewInt(1600000000),
		GasPrice: gasPrice,
		Gas:      30000,
		To:       &toAddr,
		Value:    big.NewInt(1),
	})
	// signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	// txLegacy := types.NewTransaction(nonce, toAddr, value, 30000, big.NewInt(100000000), nil)
	//	signedTxLegacy, err := types.SignTx(txLegacy, types.NewEIP155Signer(chainID), privateKey)
	//	CheckErr(err)
	ts := types.Transactions{signedTx}
	b := new(bytes.Buffer)
	ts.EncodeIndex(0, b)
	rawTxBytes := b.Bytes()
	txToSend := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &txToSend)

	err = client.SendTransaction(context.Background(), txToSend)
	if err != nil {
		log.Println("send transaction error", err)
	}
	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
}
