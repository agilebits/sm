package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/josegonzalez/sm/secrets"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

  cat app-config.yml | sm encrypt --env aws --region us-east-1 --master arn:aws:kms:us-east-1:123123123123:key/d845cfa3-0719-4631-1d00-10ab63e40ddf > encrypted-app-config.sm
`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		message, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal("failed to read:", err)
		}

		env := viper.GetString("env")
		region := viper.GetString("region")
		masterKeyID := viper.GetString("master")
		out := viper.GetString("out")

		encryptSecret(env, region, masterKeyID, message, out)
	},
}

func encryptSecret(env string, region string, masterKeyID string, message []byte, out string) {
	envelope, err := secrets.EncryptEnvelope(env, region, masterKeyID, message)
	if err != nil {
		log.Fatal("failed to encrypt:", err)
	}

	buf, err := json.Marshal(envelope)
	if err != nil {
		log.Fatal("failed to Marshal:", err)
	}

	if out != "" {
		f, err := os.Create(out)
		if err != nil {
			log.Fatal(fmt.Sprintf("failed to open %s for writing", out))
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		_, err = w.WriteString(string(buf))
		if err != nil {
			log.Fatal(fmt.Sprintf("failed to write output to %s", out))
		}
		w.Flush()
		fmt.Println(fmt.Sprintf("output written to %s", out))

		manifest := "./.sm/manifest"
		unencryptedFile := strings.TrimSuffix(out, ".sm")
		if _, err := os.Stat(manifest); !os.IsNotExist(err) {
			err = secrets.EnsureInManifest(manifest, unencryptedFile)
			if err != nil {
				log.Fatal("failed to update manifest", err)
			}
		}

		if _, err := os.Stat("./.gitignore"); !os.IsNotExist(err) {
			err = secrets.EnsureInManifest("./.gitignore", unencryptedFile)
			if err != nil {
				log.Fatal("failed to update gitignore", err)
			}
		}
	} else {
		fmt.Println(string(buf))
	}
}

func init() {
	RootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringP("env", "e", "dev", "Environment type: 'dev' or 'aws")
	encryptCmd.Flags().StringP("region", "r", "", "AWS Region ('us-east-1')")
	encryptCmd.Flags().StringP("master", "m", "", "Master key identifier")
	encryptCmd.Flags().StringP("out", "o", "", "A file to write the output to")
	viper.BindPFlag("env", encryptCmd.Flags().Lookup("env"))
	viper.BindPFlag("region", encryptCmd.Flags().Lookup("region"))
	viper.BindPFlag("master", encryptCmd.Flags().Lookup("master"))
	viper.BindPFlag("out", encryptCmd.Flags().Lookup("out"))
}
