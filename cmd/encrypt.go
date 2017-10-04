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

  cat app-config.yml | sm encrypt --env aws --region us-east-1 --master arn:aws:kms:us-east-1:123123123123:key/d845cfa3-0719-4631-1d00-10ab63e40ddf	> encrypted-app-config.sm
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

	encryptCmd.Flags().StringP("env", "e", "dev", "Environment type: 'dev' or 'aws")
	encryptCmd.Flags().StringP("region", "r", "", "AWS Region ('us-east-1')")
	encryptCmd.Flags().StringP("master", "m", "", "Master key identifier")
	viper.BindPFlag("env", encryptCmd.Flags().Lookup("env"))
	viper.BindPFlag("region", encryptCmd.Flags().Lookup("region"))
	viper.BindPFlag("master", encryptCmd.Flags().Lookup("master"))
}
