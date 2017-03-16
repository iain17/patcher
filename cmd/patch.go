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
	"github.com/iain17/patcher/patcher"
	"os"
	"github.com/iain17/logger"
)

var (
	outPath string
	patchPath string
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Using a patch file. Patch an executable by searching for signatures and overwriting a sequence of bytes.",
	Long: `Provide a json file with the patches you'd like to make. Add the signature, this tool will add the address and make a patched version'`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			logger.Errorf("file '%s' you gave to read does not exist.", filePath)
			return
		}

		if _, err := os.Stat(patchPath); os.IsNotExist(err) {
			logger.Errorf("file '%s' you gave up as the patch file, does not exist.", patchPath)
			return
		}

		if outPath == "" {
			logger.Error("Please provide a out path. Where the patched file will be saved.")
			return
		}

		patcher, err := patcher.New(patchPath, filePath, outPath)
		if err != nil {
			logger.Error(err)
			return
		}
		err = patcher.Find()
		if err != nil {
			logger.Error(err)
			return
		}
		err = patcher.Patch()
		if err != nil {
			logger.Error(err)
			return
		}
		err = patcher.Save()
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Info("Patching finished.")
	},
}

func init() {
	RootCmd.AddCommand(patchCmd)

	patchCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path of file we want to patch")
	patchCmd.Flags().StringVarP(&outPath, "output", "o", "", "Path of patched file.")
	patchCmd.Flags().StringVarP(&patchPath, "patch", "p", "", "Path of patch file.")

}
