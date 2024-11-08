package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	ollama "github.com/ollama/ollama/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Global Flags
var ollamaHost string
var baseModel string
var modelName string
var contextWindow int
var verbose bool

// Clients
var ollamaClient *ollama.Client

func init() {
	rootCmd.PersistentFlags().StringVar(&ollamaHost, "ollama-host", "http://127.0.0.1:11434", "Ollama API Endpoint To Utilize")
	rootCmd.PersistentFlags().StringVarP(&baseModel, "base-model", "m", "qwen2.5-coder:14b", "Base Ollama Model")
	rootCmd.PersistentFlags().IntVar(&contextWindow, "ctx", 16384, "Context Window")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Set Debug Output")
	// TODO: make the model name based off the base model & remove flag
	rootCmd.PersistentFlags().StringVar(&modelName, "model-name", "git-auto-message", "Model Name")
}

func LoadClients() error {
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	ollamaUrl, err := url.Parse(ollamaHost)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("URL Parse failed")
		return err
	}
	ollamaClient = ollama.NewClient(ollamaUrl, http.DefaultClient)

	return nil
}

var rootCmd = &cobra.Command{
	Use: "git-auto-message",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("--- GitAutoMessage 🦙🦾 ---")
	},
}

func Execute() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
