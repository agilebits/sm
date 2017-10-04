package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/agilebits/sm/secrets"
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
