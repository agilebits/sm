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

		if out != "" {
			f, err := os.Create(out)
			if err != nil {
				log.Fatal(fmt.Sprintf("failed to open %s for writing", out))
			}
			defer f.Close()

			w := bufio.NewWriter(f)
			_, err = w.WriteString(string(result))
			if err != nil {
				log.Fatal(fmt.Sprintf("failed to write output to %s", out))
			}
			w.Flush()
			fmt.Println(fmt.Sprintf("output written to %s", out))
		} else {
			fmt.Println(string(result))
		}
	},
}

func init() {
	RootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringVarP(&out, "out", "o", "", "A file to write the output to")
}
