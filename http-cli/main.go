package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

var contract *gateway.Contract
var gw *gateway.Gateway

func main() {
	if err := initContract(); err != nil {
		panic(err)
	}
	defer gw.Close()

	http.HandleFunc("/query", queryHandle)
	http.HandleFunc("/invoke", invokeHandle)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

type TransactionRequest struct {
	FuncName     string                 `json:"funcName"`
	Args         []string               `json:"args"`
	TransientMap map[string]interface{} `json:"transientMap"`
}

type Response struct {
	Status string `json:"status"`
}

type SuccessResponse struct {
	Response
	Data interface{} `json:"data"`
}
type FailedResponse struct {
	Response
	Err string `json:"err"`
}

func invokeHandle(rw http.ResponseWriter, r *http.Request) {
	tr, t, err := createTransaction(r)
	if err != nil {
		writerFailed(rw, fmt.Errorf("parse params error: %v", err))
		return
	}
	res, err := t.Submit(tr.Args...)
	if err != nil {
		writerFailed(rw, err)
		return
	}
	writerSuccess(rw, res)
}

func queryHandle(rw http.ResponseWriter, r *http.Request) {
	tr, t, err := createTransaction(r)
	if err != nil {
		writerFailed(rw, fmt.Errorf("parse params error: %v", err))
		return
	}
	res, err := t.Evaluate(tr.Args...)
	if err != nil {
		writerFailed(rw, err)
		return
	}
	writerSuccess(rw, res)
}

func writerSuccess(rw http.ResponseWriter, res []byte) {
	var m interface{}
	if err := json.Unmarshal(res, &m); err != nil {
		log.Printf("Unmarshal err: %v", err)
	}
	response := SuccessResponse{
		Response: Response{Status: "success"},
		Data:     m,
	}
	bytes, _ := json.Marshal(response)
	rw.Write(bytes)
}

func writerFailed(rw http.ResponseWriter, err error) {
	response := FailedResponse{
		Response: Response{Status: "failed"},
		Err:      err.Error(),
	}
	bytes, _ := json.Marshal(response)
	rw.Write(bytes)
}

func createTransaction(r *http.Request) (*TransactionRequest, *gateway.Transaction, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		return nil, nil, err
	}
	log.Printf("recived request body: %s", string(body))
	tr := new(TransactionRequest)
	if err := json.Unmarshal(body, tr); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return nil, nil, err
	}
	tm := make(map[string][]byte)
	for k, v := range tr.TransientMap {
		tm[k] = []byte(fmt.Sprintf("%v", v))
	}
	t, err := contract.CreateTransaction(tr.FuncName, gateway.WithTransient(tm), gateway.WithEndorsingPeers("peer0.org1.example.com"))
	if err != nil {
		log.Printf("Failed to create transaction : %v", err)
	}
	return tr, t, err
}

func initContract() error {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
		return err
	}
	if !wallet.Exists("appUser") {
		if err = populateWallet(wallet); err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
			return err
		}
	}
	ccpPath := filepath.Join(
		"..",
		"sampleconfig",
		"msp",
		"connection-profile.yaml",
	)

	gw, err = gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
		return err
	}

	network, err := gw.GetNetwork("ch1")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
		return err
	}

	contract = network.GetContract("mycc")
	return nil
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"sampleconfig",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "peer.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("SampleOrg", string(cert), string(key))

	return wallet.Put("appUser", identity)
}
