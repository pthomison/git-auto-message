package cmd

import (
	"fmt"
	"os"
	"regexp"

	ollama "github.com/ollama/ollama/api"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var systemMessage string = `You are a git commit message assistant. You are to analyze git diff messages and produce a short commit message that explains those changes. You should follow these rules:
- Do not explain yourself. Your output will directly be used for commit messages, so do not add additional information
- Respond in 15 words or less
- Respond in a single line
- Attempt to note every change of substance in the git diff
- You can use emojis, but don't have to
- Do not use quotes!
- Always annotate your messages with [🦾 automessage]
`

func init() {
	rootCmd.AddCommand(loadModelCmd)
}

var loadModelCmd = &cobra.Command{
	Use: "load-model",
	Run: loadModel,
}

func loadModel(cmd *cobra.Command, args []string) {
	LoadClients()

	modelfile := ollama.CreateRequest{
		Model: modelName,
		Parameters: map[string]interface{}{
			"num_ctx": contextWindow,
		},
		System: systemMessage,
		From:   baseModel,
	}

	bar := progressbar.DefaultBytes(
		100,
		"Pulling Base Model",
	)

	updateBar := func(p ollama.ProgressResponse) error {
		max := p.Total
		if bar.GetMax64() != max {
			bar.ChangeMax64(max)
		}
		return bar.Set64(p.Completed)
	}

	err := ollamaClient.Create(cmd.Context(), &modelfile, func(p ollama.ProgressResponse) error {
		if matched, _ := regexp.MatchString("pulling", p.Status); matched {
			updateBar(p)
		} else {
			fmt.Println(p.Status)
		}
		return nil
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Create Model Request failed")
		os.Exit(1)
	}
}
