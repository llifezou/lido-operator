package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type DepositData struct {
	Pubkey                string `json:"pubkey"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
	Amount                int64  `json:"amount"`
	Signature             string `json:"signature"`
	DepositMessageRoot    string `json:"deposit_message_root"`
	DepositDataRoot       string `json:"deposit_data_root"`
	ForkVersion           string `json:"fork_version"`
	NetworkName           string `json:"network_name"`
	DepositCliVersion     string `json:"deposit_cli_version"`
}

func SplitDepositData(depositDataPath string, splitCount int) ([][]DepositData, error) {
	data, err := os.ReadFile(depositDataPath)
	if err != nil {
		return nil, err
	}

	var depositDatas []DepositData
	err = json.Unmarshal(data, &depositDatas)
	if err != nil {
		return nil, err
	}

	depositDatasLen := len(depositDatas)

	if len(depositDatas) <= splitCount {
		return [][]DepositData{depositDatas}, nil
	}

	count := len(depositDatas) / splitCount
	if len(depositDatas)%splitCount != 0 {
		count++
	}

	var batchDepositDatas = [][]DepositData{}
	splitLen := 0
	for i := 0; i < count; i++ {
		start := i * splitCount
		end := start + splitCount
		if i == count-1 {
			end = len(depositDatas)
		}
		tem := depositDatas[start:end]

		splitLen += len(tem)

		batchDepositDatas = append(batchDepositDatas, tem)
	}

	if splitLen != depositDatasLen {
		fmt.Println("Length mismatch after splitting")
		os.Exit(1)
	}

	return batchDepositDatas, nil
}
