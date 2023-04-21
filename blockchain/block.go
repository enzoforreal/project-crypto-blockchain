package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// hashBlock : Calcule le hash SHA-256 d'un bloc
func hashBlock(block *Block) string {
	data := fmt.Sprintf("%d%s%v%s%d", block.Index, block.Timestamp.String(), block.Transactions, block.PrevHash, block.Nonce)
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

// validateBlock : Vérifie si un bloc est valide
func validateBlock(block *Block) bool {
	hash := hashBlock(block)
	if hash[:4] != "0000" {
		return false
	}
	return true
}
