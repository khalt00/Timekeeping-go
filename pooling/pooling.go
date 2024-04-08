package pooling

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"os"
	"time"
	"timekeeping/lib/config"
	"timekeeping/model"

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

func Pooling(config config.Config) {
	client, err := ethclient.Dial(config.AlchemyUrl)
	// client, err := ethclient.Dial("wss://127.0.0.1:8545/")

	if err != nil {
		log.Fatal(err)
	}

	// Load ABI from JSON file
	abiBytes, err := os.ReadFile("./abi/abi.json")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(config.SmartContractAddress)

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
		case logData := <-ch:
			if logData.Topics[0].Hex() == attendanceEvent.Hex() {
				data, err := contractAbi.Unpack("AttendanceEvent", logData.Data)
				if err != nil {
					log.Println("Failed to decode event:", err)
					continue
				}
				// var record model.AttendanceRecord
				err = insertAttendanceToDatabase(data)
				if err != nil {
					log.Println("Failed to insert to database:", err, data[0])
					continue
				}
				log.Println("inserted attendance", data[0])
			}
			if logData.Topics[0].Hex() == updateAttendanceEvent.Hex() {
				data, err := contractAbi.Unpack("UpdateAttendanceEvent", logData.Data)
				if err != nil {
					fmt.Println("Failed to decode event:", err)
					continue
				}
				fmt.Println("UpdateAttendanceEvent:", data)
			}
		}
	}

}

func insertAttendanceToDatabase(data []interface{}) error {
	// var record model.AttendanceRecord
	recordType := ""
	if data[3].(uint8) == 0 {
		recordType = "CHECKIN"
	} else {
		recordType = "CHECKOUT"
	}

	eventTimestamp := data[1].(*big.Int).Int64()
	eventDate := time.Unix(eventTimestamp, 0).Format("02/01/2006")
	eventHours := time.Unix(eventTimestamp, 0).Format("15:04")
	record := model.AttendanceRecord{
		EmployeeID:     data[0].(common.Address).Hex(),
		EventTimestamp: data[1].(*big.Int).String(),
		Details:        data[2].(string),
		RecordType:     recordType,
		EventDate:      eventDate,
		EventHours:     eventHours,
	}
	err := model.InsertAttendanceRecord(record)
	if err != nil {
		log.Println("model.InsertAttendanceRecord(record)", err)
	}

	return nil
}
