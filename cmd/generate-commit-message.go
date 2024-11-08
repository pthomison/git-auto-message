package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	ollama "github.com/ollama/ollama/api"
)

var (
	autoBuild bool
)

func init() {
	generateCommitMessageCmd.PersistentFlags().BoolVarP(&autoBuild, "autobuild", "b", false, "Build the model if not present")

	rootCmd.AddCommand(generateCommitMessageCmd)
}

var generateCommitMessageCmd = &cobra.Command{
	Use: "generate-commit-message",
	Run: generateCommitMesssage,
}

func generateCommitMesssage(cmd *cobra.Command, args []string) {
	LoadClients()

	if autoBuild {
		listResp, err := ollamaClient.List(cmd.Context())
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("ollama list")
			os.Exit(1)
		}

		buildNeeded := true
		for _, model := range listResp.Models {
			if model.Model == fmt.Sprintf("%v:latest", modelName) {
				buildNeeded = false
			}
		}

		if buildNeeded {
			loadModel(cmd, args)
		}
	}

	out, err := exec.Command("git", "diff", "--cached", ".").Output()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("git diff failed")
		os.Exit(1)
	}

	// TODO: add flag for keepalive
	gResp, err := generate(cmd.Context(), ollamaClient, string(out), time.Second*-1)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Generate Request failed")
		os.Exit(1)
	}

	fmt.Println(gResp.Response)

	if verbose {
		gResp.Summary()
	}
}

func generate(ctx context.Context, client *ollama.Client, prompt string, keepAlive time.Duration) (*ollama.GenerateResponse, error) {
	f := false
	block := make(chan error)
	var resp ollama.GenerateResponse
	go func() {
		err := client.Generate(ctx, &ollama.GenerateRequest{
			Model:     modelName,
			Prompt:    prompt,
			Stream:    &f,
			KeepAlive: &ollama.Duration{Duration: keepAlive},
			// Suffix: "test",
		}, func(gr ollama.GenerateResponse) error {
			resp = gr
			return nil
		})
		block <- err
	}()

	err := <-block

	return &resp, err
}
