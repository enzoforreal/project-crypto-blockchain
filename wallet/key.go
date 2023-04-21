package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
)

// GenerateKeyPair : Génère une paire de clés publique/privée
func GenerateKeyPair() (*ecdsa.PrivateKey, []byte, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	return privKey, pubKey, nil
}

// SavePrivateKey : Stocke la clé privée dans un fichier
func SavePrivateKey(filename string, Key *ecdsa.PrivateKey) error {
	keyBytes, err := x509.MarshalECPrivateKey(Key)
	if err != nil {
		return err
	}

	keyPem := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	return ioutil.WriteFile(filename, keyPem, 0600)
}

// LoadPrivateKey : Charge la clé privée à partir d'un fichier
func LoadPrivateKey(filename string) (*ecdsa.PrivateKey, error) {
	keyPem, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyPem)
	if block == nil {
		return nil, errors.New("Clé privée invalide")
	}

	Key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return Key, nil
}

// SignTransaction : Signe une transaction avec la clé privée
func SignTransaction(key *ecdsa.PrivateKey, txHash []byte) ([]byte, []byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, key, txHash)
	if err != nil {
		return nil, nil, err
	}
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	return rBytes, sBytes, nil
}

// VerifyTransaction : Vérifie la signature de la transaction avec la clé publique
func VerifyTransaction(pubKey []byte, txHash []byte, rBytes []byte, sBytes []byte) bool {
	xBytes := pubKey[:len(pubKey)/2]
	yBytes := pubKey[len(pubKey)/2:]
	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)
	pub := ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)
	return ecdsa.Verify(&pub, txHash, r, s)
}
