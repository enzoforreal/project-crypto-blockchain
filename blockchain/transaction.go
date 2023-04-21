package blockchain

import (
	"math/big"
)

// Transaction : La structure d'une transaction// Transaction : La structure d'une transaction
type Transaction struct {
	Sender    string   // L'adresse du portefeuille de l'expéditeur
	Recipient string   // L'adresse du portefeuille du destinataire
	Amount    int      // Le montant de la transaction
	Inputs    []Input  // Les entrées de transaction
	Outputs   []Output // Les sorties de transaction
	Signature []byte   // La signature numérique de la transaction
}


// Output : La structure d'une sortie de transaction
type Output struct {
	Address string // L'adresse du portefeuille du destinataire
	Amount  int    // Le montant de la transaction

}

type Signature struct {
	R, S *big.Int // La signature numérique de la transaction
}


// Input : La structure d'une entrée de transaction
type Input struct {
	Txid      []byte // L'identifiant de transaction (hash)
	Vout      int    // La sortie de transaction à dépenser
	ScriptSig string // La signature script qui déverrouille les fonds
}

// NewTransaction : Crée une nouvelle instance de transaction
func NewTransaction(sender, recipient string, amount int) *Transaction {
	return &Transaction{Sender: sender, Recipient: recipient, Amount: amount}
}
