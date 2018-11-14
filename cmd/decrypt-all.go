package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/josegonzalez/sm/secrets"
	"github.com/spf13/cobra"
)

// decryptAllCmd represents the decrypt-all command
var decryptAllCmd = &cobra.Command{
	Use:   "decrypt-all",
	Short: "Decrypt all files in manifest using key management system",
	Long: `This command will decrypt all files in the .sm/manifest.

It requires access to the same key management system (KMS) that was used for encryption.

For example:

  sm decrypt-all

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
				message, err := ioutil.ReadFile(fmt.Sprintf("%s.sm", line))
				if err == nil {
					decryptSecretAndWrite(message, line)
				} else {
					log.Fatal("failed to read:", err)
				}
			}(line)
		}

		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(decryptAllCmd)
}
