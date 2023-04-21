package blockchain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// BlockChain : La structure de la blockchain
type BlockChain struct {
	Chain               []*Block        // La liste des blocs de la blockchain
	CurrentTransactions []Transaction   // La liste des transactions en attente
	Nodes               map[string]bool // La liste des noeuds du réseau
}

// Block : La structure d'un bloc
type Block struct {
	Index        int           // L'index du bloc dans la chaîne
	Timestamp    time.Time     // La date et l'heure de création du bloc
	Transactions []Transaction // La liste des transactions du bloc
	PrevHash     string        // Le hash du bloc précédent
	Hash         string        // Le hash du bloc actuel
	Nonce        int           // Le nonce utilisé pour le minage
}

// NewBlockchain : Crée une nouvelle instance de blockchain
func NewBlockchain() *BlockChain {
	chain := make([]*Block, 0)
	transactions := make([]Transaction, 0)
	nodes := make(map[string]bool)
	nodes["localhost"] = true

	chain = append(chain, createGenesisBlock())

	return &BlockChain{Chain: chain, CurrentTransactions: transactions, Nodes: nodes}
}

// createGenesisBlock : Crée le bloc initial de la blockchain
func createGenesisBlock() *Block {
	genesis := &Block{Index: 0, Timestamp: time.Now(), Transactions: []Transaction{}, PrevHash: "", Nonce: 0}
	genesis.Hash = hashBlock(genesis)
	return genesis
}

// AddTransaction : Ajoute une transaction à la liste des transactions en attente
func (bc *BlockChain) AddTransaction(t Transaction) int {
	bc.CurrentTransactions = append(bc.CurrentTransactions, t)
	return len(bc.Chain) + 1
}

// MineBlock : Mine un nouveau bloc à partir des transactions en attente
func (bc *BlockChain) MineBlock() *Block {
	lastBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := &Block{Index: lastBlock.Index + 1, Timestamp: time.Now(), Transactions: bc.CurrentTransactions, PrevHash: lastBlock.Hash, Nonce: 0}

	// Trouver le nonce pour le minage
	for !validateBlock(newBlock) {
		newBlock.Nonce++
		newBlock.Hash = hashBlock(newBlock)
	}

	bc.Chain = append(bc.Chain, newBlock)
	bc.CurrentTransactions = []Transaction{}

	return newBlock
}

// RegisterNode : Ajoute un nouveau noeud à la liste des noeuds
func (bc *BlockChain) RegisterNode(address string) {
	bc.Nodes[address] = true
}

// isValidChain : Vérifie si une chaîne de blocs est valide
func isValidChain(chain []Block) bool {
	prevHash := ""
	for _, block := range chain {
		if block.Hash != hashBlock(&block) {
			return false
		}

		if block.PrevHash != prevHash {
			return false
		}
		prevHash = block.Hash
	}

	return true
}

// ResolveConflicts : Résout les conflits de chaîne en remplaçant la chaîne courante par la chaîne la plus longue et sécurisée du réseau
func (bc *BlockChain) ResolveConflicts() bool {
	var newChain []*Block
	maxLength := len(bc.Chain)
	maxNonce := -1

	// Récupérer la chaîne la plus longue et sécurisée du réseau
	for node := range bc.Nodes {
		resp, err := http.Get(fmt.Sprintf("http://%s/blockchain", node))
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var chain BlockChain
			err := json.NewDecoder(resp.Body).Decode(&chain)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Vérifier si la chaîne est plus longue et valide, ou si elle a la même longueur mais une preuve de travail plus grande
			if len(chain.Chain) > maxLength || (len(chain.Chain) == maxLength && chain.Chain[len(chain.Chain)-1].Nonce > maxNonce) {
				chainCopy := make([]Block, len(chain.Chain))
				for i, b := range chain.Chain {
					chainCopy[i] = *b
				}
				if isValidChain(chainCopy) {
					maxLength = len(chain.Chain)
					maxNonce = chain.Chain[len(chain.Chain)-1].Nonce
					newChain = chain.Chain
				}
			}
		}
	}

	// Remplacer la chaîne courante par la nouvelle chaîne si elle est plus longue et sécurisée
	if newChain != nil {
		bc.Chain = newChain
		return true
	}

	return false
}
