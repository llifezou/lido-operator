package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"path/filepath"
)

var (
	operatorId      int
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
	keyGenerator.PersistentFlags().StringVarP(&depositDataPath, "depositDataPath", "d", "", "deposit data file path")
	keyGenerator.PersistentFlags().StringVarP(&outputDir, "inputDataDir", "i", ".", "inputdata output dir")
	keyGenerator.PersistentFlags().IntVarP(&operatorId, "operatorId", "o", -1, "operator id")
}

func main() {
	rootCmd.AddCommand(keyGenerator)
	_ = rootCmd.Execute()
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

		batchDepositDatas, err := SpiteDepositData(depositDataPath)
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
			writeFunc(i+1, []byte("0x"+common.Bytes2Hex(inputData)))
		}
	},
}

func writeFunc(i int, d []byte) {
	fi, err := os.Create(filepath.Join(outputDir, fmt.Sprintf("%d-transactionData.txt", i)))
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
