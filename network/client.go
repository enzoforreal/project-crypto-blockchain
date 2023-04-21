package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/enzoforreal/project-crypto-blockchain/blockchain"
)

// RegisterNode : Enregistre un nouveau noeud sur le réseau
func RegisterNode(address string, nodes []string) error {
	reqData := map[string][]string{
		"nodes": nodes,
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/nodes/register", address), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// ResolveConflicts : Résout les conflits de chaîne en remplaçant la chaîne actuelle par la plus longue
func ResolveConflicts(nodes []string) error {
	maxLength := 0
	var newChain []*blockchain.Block

	for _, node := range nodes {
		resp, err := http.Get(fmt.Sprintf("http://%s/blockchain", node))
		if err == nil && resp.StatusCode == http.StatusOK {
			var chain blockchain.BlockChain
			err = json.NewDecoder(resp.Body).Decode(&chain)
			if err == nil {
				length := len(chain.Chain)
				if length > maxLength && blockchain.ValidateChain(&chain) {
					maxLength = length
					newChain = chain.Chain
				}
			}
		}
	}

	if newChain != nil {
		blockchain.Chain = newChain
		return nil
	}

	return fmt.Errorf("Aucune chaîne plus longue trouvée")
}
