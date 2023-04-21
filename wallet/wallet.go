package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/enzoforreal/project-crypto-blockchain/blockchain"
)

// Wallet : La structure du portefeuille
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey // La clé privée du portefeuille
	PublicKey  []byte            // La clé publique du portefeuille
	Address    string            // L'adresse du portefeuille (dérivée de la clé publique)
}

// NewWallet : Crée un nouveau portefeuille
func NewWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(ecdsaelliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	address := getAddress(publicKey)

	return &Wallet{PrivateKey: privateKey, PublicKey: publicKey, Address: address}, nil
}

// getAddress : Génère l'adresse à partir de la clé publique
func getAddress(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	return fmt.Sprintf("%x", hash)
}

// SignTransaction : Signe une transaction avec la clé privée du portefeuille
func (wallet *Wallet) SignTransaction(tx *blockchain.Transaction) error {
	if !wallet.HasSufficientFunds(tx) {
		return errors.New("Fonds insuffisants pour la transaction")
	}

	txInput := blockchain.Input{Address: wallet.Address, Amount: tx.Inputs[0].Amount}
	tx.Inputs = []blockchain.Input{txInput}
	txData := blockchain.SerializeTransaction(tx)
	r, s, err := ecdsa.Sign(rand.Reader, wallet.PrivateKey, txData)
	if err != nil {
		return err
	}
	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signature = signature

	return nil
}

// HasSufficientFunds : Vérifie si le portefeuille a suffisamment de fonds pour la transaction
func (wallet *Wallet) HasSufficientFunds(tx *blockchain.Transaction) bool {
	inputAmount := 0
	for _, input := range tx.Inputs {
		if input.Address == wallet.Address {
			inputAmount += input.Amount
		}
	}
	outputAmount := 0
	for _, output := range tx.Outputs {
		if output.Address == wallet.Address {
			outputAmount += output.Amount
		}
	}
	return inputAmount >= outputAmount
}

// String : Affiche les informations du portefeuille
func (wallet *Wallet) String() string {
	return fmt.Sprintf("Portefeuille - Adresse : %s, Clé publique : %s", wallet.Address, hex.EncodeToString(wallet.PublicKey))
}
