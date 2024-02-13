package wallet

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	Mnemonic   string
	PrivateKey string
	PublicKey  string
	Address    common.Address
}

func GetRandomWallet() Wallet {
	// Generate random entropy (higher entropy means higher security)
	entropy, _ := bip39.NewEntropy(128) // 128 bits of entropy

	// Generate mnemonic
	mnemonic, _ := bip39.NewMnemonic(entropy)

	// Convert mnemonic to seed
	seed := bip39.NewSeed(mnemonic, "")

	// Derive private key from seed
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0") //最后一位是同一个助记词的地址id，从0开始，相同助记词可以生产无限个地址
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}
	address := account.Address
	privateKey, _ := wallet.PrivateKeyHex(account)
	publicKey, _ := wallet.PublicKeyHex(account)
	return Wallet{
		PrivateKey: privateKey,
		Address:    address,
		Mnemonic:   mnemonic,
		PublicKey:  publicKey,
	}
}
