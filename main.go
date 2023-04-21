package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/enzoforreal/project-crypto-blockchain/blockchain"
	"github.com/gorilla/mux"
)

func main() {
	// Créer une nouvelle instance de blockchain
	chain := blockchain.NewBlockchain()

	// Configurer le routeur HTTP
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/blockchain", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(chain)
	}).Methods("GET")
	router.HandleFunc("/transaction", func(w http.ResponseWriter, r *http.Request) {
		// Créer une nouvelle transaction
		decoder := json.NewDecoder(r.Body)
		var t blockchain.Transaction
		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer r.Body.Close()

		// Ajouter la transaction à la blockchain
		chain.AddTransaction(t)
	}).Methods("POST")

	router.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {
		// Miner un nouveau bloc
		chain.MineBlock()
	}).Methods("POST")

	router.HandleFunc("/nodes/register", func(w http.ResponseWriter, r *http.Request) {
		// Enregistrer un nouveau noeud
		decoder := json.NewDecoder(r.Body)
		var nodes []string
		err := decoder.Decode(&nodes)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer r.Body.Close()

		// Ajouter les noeuds à la liste des noeuds
		for _, node := range nodes {
			chain.RegisterNode(node)
		}
	}).Methods("POST")

	router.HandleFunc("/nodes/resolve", func(w http.ResponseWriter, r *http.Request) {
		// Résoudre les conflits de chaîne
		chain.ResolveConflicts()
	}).Methods("GET")

	// Lancer le serveur HTTP
	log.Fatal(http.ListenAndServe(":8080", router))
}
