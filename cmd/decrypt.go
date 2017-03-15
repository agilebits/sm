// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/AgileBits/sm/secrets"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt content using key management system",
	Long: `This command will decrypt content that was encrypted using encrypt command. 
	
It requires access to the same key management system (KMS) that was used for encryption.

For example:

  cat encrypted-app-config.sm | sm decrypt > app-config.yml

`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		message, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal("failed to read:", err)
		}

		envelope := &secrets.Envelope{}
		if err := json.Unmarshal(message, &envelope); err != nil {
			log.Fatal("failed to Unmarshal:", err)
		}

		result, err := secrets.DecryptEnvelope(envelope)
		if err != nil {
			log.Fatal("failed to DecryptEnvelope:", err)
		}

		fmt.Println(string(result))
	},
}

func init() {
	RootCmd.AddCommand(decryptCmd)
}
