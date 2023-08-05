package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"log"
)

func main() {
	gen()
	// getPrivateKey()
}

func gen() {
	// Generate random entropy (higher entropy means higher security)
	entropy, _ := bip39.NewEntropy(128) // 128 bits of entropy

	// Generate mnemonic
	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println("Mnemonic:", mnemonic)

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
	address := account.Address.Hex()
	privateKey, _ := wallet.PrivateKeyHex(account)
	publicKey, _ := wallet.PublicKeyHex(account)
	fmt.Printf("Private Key: %s\n", privateKey)
	fmt.Printf("Address: %s\n", address)
	fmt.Printf("Public Key: %s\n", publicKey)

	// Generate Keystore file
	ks := keystore.NewKeyStore(".", keystore.StandardScryptN, keystore.StandardScryptP)
	privK, _ := wallet.PrivateKey(account)

	accountKeystore, _ := ks.ImportECDSA(privK, "testonly")
	ks.Export(accountKeystore, "", "")
}

func getPrivateKey() {
	keystoreJSON := []byte(`{"address":"bf08ad4296878143f577ad60c35fc92c0a52d974","crypto":{"cipher":"aes-128-ctr","ciphertext":"9f4a043625329c183ec98d89e8b7cb670d4ec2e52e804b9d4d72328b6ffbd076","cipherparams":{"iv":"00e1e1437119cfbf0305df1db4b76ceb"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"b1cc13cec0b7df5cbd3ea8b00f55d514b7d6df56d6112fe25c5feb5f7410fc9b"},"mac":"379077256a8d4b177c4e6602f2e2b0d2700cd6dc498dcc8f751efec6f21e2ccb"},"id":"ffbfd36e-51f5-409e-a53d-06cca15b9ca7","version":3}`)

	password := "testonly"

	key, err := keystore.DecryptKey(keystoreJSON, password)
	if err != nil {
		log.Fatal(err)
	}

	privateKey := key.PrivateKey
	fmt.Printf("Private Key: %x\n", privateKey.D.Bytes())
}
