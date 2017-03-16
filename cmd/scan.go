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
	"github.com/iain17/logger"
	"os"
	"github.com/iain17/patcher/scanner"
)

var (
	// File path we want to scan in
	filePath string
	// Sigature we'll use to find the memory address.
	signature string
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for a signature",
	Long: `Find a address based on a given signature in a binary file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			logger.Errorf("file '%s' you gave up to scan does not exist.", filePath)
			return
		}
		if len(signature) == 0 {
			logger.Error("Please provide a signature you'd like to scan for. E.G. 'E8 ? ? ? ?' will give you the first call in i368")
			return
		}

		res, err := scanner.Scan(signature, filePath)
		if err != nil {
			logger.Error("Please provide a signature you'd like to scan for. E.G. 'E8 ? ? ? ?' will give you the first call in i368")
		} else {
			logger.Infof("sig(%s) found on %#08x", signature, res)
		}

	},
}

func init() {
	RootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&filePath, "file", "f", "", "File path we want to scan in")
	scanCmd.Flags().StringVarP(&signature, "signature", "s", "", "Signature you want to scan for. E.G. 'E8 ? ? ? ?' will give you the first call in i368")
}
