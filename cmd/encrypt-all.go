package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/josegonzalez/sm/secrets"
	"github.com/spf13/cobra"
)

// encryptAllCmd represents the encrypt-all command
var encryptAllCmd = &cobra.Command{
	Use:   "encrypt-all",
	Short: "Encrypt all files in manifest using key management system",
	Long: `This command will re-encrypt all changed files in the .sm/manifest.

It requires access to the same key management system (KMS) that was used for the initial
encryption. The key will be extracted from the existing encrypted file. If an unencrypted file
does not exist, this will be logged and skipped.

For example:

  sm encrypt-all

`,
	Run: func(cmd *cobra.Command, args []string) {
		lines, err := secrets.ReadManifest("./.sm/manifest")
		if err != nil {
			log.Fatal("error reading manifest:", err)
		}

		var wg sync.WaitGroup
		wg.Add(len(lines))

		for _, line := range lines {
			go func(line string) {
				defer wg.Done()
				existing, err := ioutil.ReadFile(line)
				if err != nil {
					fmt.Println(fmt.Sprintf("Skipping missing unencrypted file %s", line))
					return
				}

				encryptedFile := fmt.Sprintf("%s.sm", line)
				encryptedMessage, err := ioutil.ReadFile(encryptedFile)
				if err != nil {
					log.Fatal("failed to read:", err)
					return
				}

				// decrypt the encrypted file
				decryptedContent, err := decryptSecret(encryptedMessage)
				if err != nil {
					log.Fatal("failed to decrypt:", err)
					return
				}

				if string(existing) == string(decryptedContent) {
					return
				}

				fmt.Println(fmt.Sprintf("Attempting to re-encrypt %s", line))
				envelope := &secrets.Envelope{}
				if err := json.Unmarshal(encryptedMessage, &envelope); err != nil {
					log.Fatal("failed to read envelope:", err)
					return
				}

				encryptSecret(envelope.Env, envelope.Region, envelope.MasterKeyID, existing, encryptedFile)
			}(line)
		}

		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(encryptAllCmd)
}
