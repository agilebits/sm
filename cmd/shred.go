package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/agilebits/sm/secrets"
	"github.com/spf13/cobra"
)

// shredCmd represents the shred command
var shredCmd = &cobra.Command{
	Use:   "shred",
	Short: "Shred all files listed in the manifest",
	Long: `This command will remove all files in the .sm/manifest.

For example:

  sm shred

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
				if err := os.Remove(line); err != nil {
					log.Println("failed to remove file:", err)
				} else {
					fmt.Println(fmt.Sprintf("file shredded: %s", line))
				}
			}(line)
		}

		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(shredCmd)
}
