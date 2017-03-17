// Copyright © 2017 Iain Munro <iain@imunro.nl>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"github.com/iain17/logger"
	"strconv"
	"github.com/iain17/patcher/scanner"
	"strings"
	"fmt"
)

var (
	address string
	sigLength int32
)

// sigGenCmd represents the sigGen command
var sigGenCmd = &cobra.Command{
	Use:   "sigGen",
	Short: "Generate a signature based on a address",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			logger.Errorf("file '%s' you gave up to scan does not exist.", filePath)
			return
		}
		//Convert address hex string to int64
		val, err := strconv.ParseInt(strings.Replace(address, "0x", "", -1), 16,64)
		if err != nil {
			logger.Error(err)
			return
		}
		if val == int64(0) {
			logger.Error("Please provide the address you'd like the signature of")
		}
		length, _ := cmd.Flags().GetInt("length")
		res, err := scanner.GenSig(val, filePath, length)
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Infof("%#08x has a sig of:", val)
		for _, byte := range res {
			fmt.Printf("%X ", byte)
		}
		println("")
		println("Be sure to '?' out any references which might be variable.")
	},
}

func init() {
	RootCmd.AddCommand(sigGenCmd)
	sigGenCmd.Flags().IntP("length", "l", 16, "Length of signature. The amount of bytes.")
	sigGenCmd.Flags().StringVarP(&filePath, "file", "f", "", "File path we want to scan in")
	sigGenCmd.Flags().StringVarP(&address, "address", "a", "0x0", "The address you'd like a signature of")
}
