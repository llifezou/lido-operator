package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"path/filepath"
)

var (
	operatorId      int
	splitCount      int
	depositDataPath string
	outputDir       string
)

var rootCmd = &cobra.Command{
	Use:   "lido-operator",
	Short: "lido-operator",
	Long:  `lido operator tool`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&depositDataPath, "depositDataPath", "d", "", "deposit data file path")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "inputDataDir", "i", ".", "inputdata output dir")
	rootCmd.PersistentFlags().IntVarP(&operatorId, "operatorId", "o", -1, "operator id")
	rootCmd.PersistentFlags().IntVarP(&splitCount, "splitCount", "s", 100, "deposit data split count")
}

func main() {
	rootCmd.AddCommand(splitDepositDataCmd, keyGenerator)
	_ = rootCmd.Execute()
}

var splitDepositDataCmd = &cobra.Command{
	Use:     "split-depositdata",
	Short:   "split-depositdata",
	Example: "./lido-operator split-depositdata",
	Run: func(cmd *cobra.Command, args []string) {
		batchDepositDatas, err := SplitDepositData(depositDataPath, splitCount)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for i, depositDatas := range batchDepositDatas {
			d, err := json.Marshal(depositDatas)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			writeFunc(i+1, "deposit-data.json", d)
		}
	},
}

var keyGenerator = &cobra.Command{
	Use:     "key-generator",
	Short:   "key-generator",
	Example: "./lido-operator key-generator",
	Run: func(cmd *cobra.Command, args []string) {
		if operatorId == -1 || depositDataPath == "" {
			fmt.Println("--depositDataPath and --operatorId are required")
			return
		}

		nodeOperatorId := big.NewInt(int64(operatorId))

		batchDepositDatas, err := SplitDepositData(depositDataPath, splitCount)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for i, depositDatas := range batchDepositDatas {
			publicKeys := []byte{}
			signatures := []byte{}
			for _, depositData := range depositDatas {
				pubkey := common.Hex2Bytes(depositData.Pubkey)
				signature := common.Hex2Bytes(depositData.Signature)

				publicKeys = append(publicKeys, pubkey...)
				signatures = append(signatures, signature...)
			}

			keysCount := big.NewInt(int64(len(depositDatas)))
			inputData, err := AddSigningKeysOperatorBH(nodeOperatorId, keysCount, publicKeys, signatures)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			writeFunc(i+1, "transactionData.txt", []byte("0x"+common.Bytes2Hex(inputData)))
		}
	},
}

func writeFunc(i int, name string, d []byte) {
	fi, err := os.Create(filepath.Join(outputDir, fmt.Sprintf("%d-%s", i, name)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = fi.Write(d)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fi.Close()
}
