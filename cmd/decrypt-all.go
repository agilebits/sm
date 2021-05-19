package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"sm/secrets"

	"github.com/spf13/cobra"
)

const defaultWorkerCount = 25

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

		workerCount := workerCount()
		workCh := make(chan string, 2*workerCount)

		for i := 0; i < workerCount; i++ {
			go worker(workCh, &wg)
		}

		for _, line := range lines {
			workCh <- line
		}

		wg.Wait()
		close(workCh)
	},
}

func init() {
	RootCmd.AddCommand(decryptAllCmd)
}

func worker(workCh chan string, wg *sync.WaitGroup) {
	for {
		select {
		case file, ok := <-workCh:
			// channel closed
			if !ok {
				return
			}

			// read up the file from disk
			message, err := ioutil.ReadFile(fmt.Sprintf("%s.sm", file))
			if err != nil {
				log.Printf("failed to read: %s", err.Error())
				continue
			}

			success := false
			for i := 0; i < 3; i++ {
				if err = decryptSecretAndWrite(message, file); err != nil {
					log.Printf("Could not decrypt file %s: %s (Attempt %d/%d, next attempt in 1s)\n", file, err.Error(), i, 3)
					time.Sleep(time.Second)
					continue
				}

				success = true
				break
			}

			if !success {
				log.Fatalf("Could not decrypt file %s: %s", file, err)
			}

			wg.Done()
		}
	}
}

func workerCount() int {
	count := os.Getenv("WORKER_COUNT")
	if count == "" {
		return defaultWorkerCount
	}

	val, err := strconv.Atoi(count)
	if err != nil {
		log.Fatal("WORKER_COUNT is not a valid int")
	}

	return val
}
