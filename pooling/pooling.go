package pooling

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	// "github.com/miguelmota/go-web3-example/greeter"
	"log"
)

type AttendanceEvent struct {
	EmployeeID     common.Address // address type in solidity corresponds to common.Address in go-ethereum
	Date           uint64         // uint256 type in solidity corresponds to uint64 in Go
	Details        string
	AttendanceType uint8 // Assuming Type is an enum or uint8 in solidity
}

func Pooling() {
	client, err := ethclient.Dial("wss://polygon-mumbai.g.alchemy.com/v2/-ONXHsPzU78tail5a8J3CvSbM0OFiiOP")

	if err != nil {
		log.Fatal(err)
	}

	myContractAddress := "0x71338c7DfBDfb243F1f58Df6535ac3C0EBC64362"
	// priv := "6a62e93e53d6c166b85235a0dd56d640b2c17567bd7ebd4e6a972baad441534f"

	// key, err := crypto.HexToECDSA(priv)

	// Load ABI from JSON file
	abiBytes, err := os.ReadFile("./abi/abi.json")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(myContractAddress)
	// greeterClient, err := greeter.NewGreeter(contractAddress, client)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// not sure why I have to set this when using testrpc
	// var nonce int64 = 0
	// auth.Nonce = big.NewInt(nonce)

	// conn := dbsvc.GetPostgresConnection()

	contractAbi, err := abi.JSON(bytes.NewReader(abiBytes))
	if err != nil {
		log.Fatal(err)
	}

	attendanceEvent := contractAbi.Events["AttendanceEvent"].ID
	updateAttendanceEvent := contractAbi.Events["UpdateAttendanceEvent"].ID

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics: [][]common.Hash{
			{attendanceEvent, updateAttendanceEvent},
		},
	}

	var ch = make(chan types.Log, 2)
	ctx := context.Background()

	sub, err := client.SubscribeFilterLogs(ctx, query, ch)

	if err != nil {
		log.Println("Subscribe:", err)
		return
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case log := <-ch:
			if log.Topics[0].Hex() == attendanceEvent.Hex() {
				data, err := contractAbi.Unpack("AttendanceEvent", log.Data)
				if err != nil {
					fmt.Println("Failed to decode event:", err)
					continue
				}
				fmt.Println("AttendanceEvent:", data[0])
			}
			if log.Topics[0].Hex() == updateAttendanceEvent.Hex() {
				data, err := contractAbi.Unpack("UpdateAttendanceEvent", log.Data)
				if err != nil {
					fmt.Println("Failed to decode event:", err)
					continue
				}
				fmt.Println("UpdateAttendanceEvent:", data)
			}
		}
	}

}
