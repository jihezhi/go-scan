package main

import (
	"context"
	"go-scan/address"
	"go-scan/token" // for demo
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// LogTransfer ..
type LogTransfer struct {
	From         common.Address
	To           common.Address
	Tokens       *big.Int
	TokenAddress common.Address
	TokenSymbol  string
	BlockNumber  int
	BlockTime    string
	Time         int64
}

func main() {
	// connect to network
	client, err := ethclient.Dial("wss://wsapi.fantom.network/")
	if err != nil {
		log.Fatal(err)
	}

	// get the ERC-20 abi
	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))
	if err != nil {
		log.Fatal(err)
	}
	contractAbi = contractAbi

	// get the last block number
	lastBlockNumber := getLastBlock(*client)

	// Encode Transfer method to Hash
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	// Get the ERC-20 transaction
	// step 1: From firebase get the started block number
	// step 2: From the started block number to last block number get the Transaction
	// step 3: Store the Transaction to firestore
	for i := 30784105; i < int(lastBlockNumber); i = i + 10000 {
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(i)),
			ToBlock:   big.NewInt(int64(i + 10000)),
			Addresses: address.Returncontractaddress(),
		}
		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Fatal("FilterLogs Error:", err)
		}
		getERC20Transaction(logs, contractAbi, logTransferSigHash, *client)
	}
}

func getLastBlock(client ethclient.Client) int64 {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return header.Number.Int64()
}

func getBlockTime(client ethclient.Client, blocknumber int) int64 {
	block,err := client.BlockByNumber(context.Background(), blocknumber.(*big.Int))
	if err != nil {
		log.Fatal("")
	}
	return client.BlockByNumber(context.Background(), blocknumber.(*big.Int)).Time().Uint64()
}

func getERC20Transaction(logs []types.Log, contractAbi abi.ABI, logTransferSigHash common.Hash, client ethclient.Client) {
	for _, vLog := range logs {
		if logTransferSigHash.Hex() == vLog.Topics[0].Hex() {
			m := make(map[string]interface{})
			var transferEvent LogTransfer
			err := contractAbi.UnpackIntoMap(m, "Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			transferEvent.TokenAddress = vLog.Address
			transferEvent.BlockNumber = int(vLog.BlockNumber)
			transferEvent.Tokens = m["tokens"].(*big.Int)
			transferEvent.Time = getBlockTime(client, transferEvent.BlockNumber)
			StoreToFirestore(transferEvent)
		}
	}
}

// get the end block number from firestore
func getEndblock() {

}
