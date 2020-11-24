package main

import (
	"os"
	"zeb-cli/collection"

	"github.com/spf13/cobra"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
		os.Exit(1)
	}
}

func run() error {
	root := cobra.Command{
		Use:   "zebedee",
		Short: "z",
		Long:  "TODO",
	}

	root.AddCommand(collection.GetCommands())
	return root.Execute()
}
