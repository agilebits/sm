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

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt content using key management system",
	Long: ` 
	
Encrypt command is used to encrypt the contents of the standard input and write
encrypted "envelope" into the standard output. 

The envelope is a JSON file that contains encrypted data along with the 
additional information that is needed to decrypt it back if the access to the
key management system is available.

For example:

  cat app-config.yml | sm encrypt --env aws --region us-east-1 --master arn:aws:kms:us-east-1:123123123123:key/d845cfa3-0719-4631-1d00-10ab63e40ddf	> encrypted-app-config.sm
`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		message, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal("failed to read:", err)
		}

		envelope, err := secrets.EncryptEnvelope(env, region, masterKeyID, message)
		if err != nil {
			log.Fatal("failed to encrypt:", err)
		}

		buf, err := json.Marshal(envelope)
		if err != nil {
			log.Fatal("failed to Marshal:", err)
		}

		fmt.Println(string(buf))
	},
}

func init() {
	RootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringVarP(&env, "env", "e", "dev", "Environment type: 'dev' or 'aws")
	encryptCmd.Flags().StringVarP(&region, "region", "r", "", "AWS Region ('us-east-1')")
	encryptCmd.Flags().StringVarP(&masterKeyID, "master", "m", "", "Master key identifier")
}
