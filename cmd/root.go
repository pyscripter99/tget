/*
Copyright © 2023 Ryder Retzlaff <ryder@retzlaff.family>
*/
package cmd

import (
	"os"

	"github.com/pyscripter99/tget/internals"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tget <url> [output]",
	Short: "A golang implementation of wget",
	Long:  `A golang implementation of wget.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if url is provided
		if len(args) < 1 {
			cmd.Help()
			return
		}

		// Check if output is provided
		if len(args) < 2 {
			args = append(args, "")
		}

		// Get url & download file
		url := args[0]
		output := args[1]
		err := internals.Download(url, output, cmd.Flag("user-agent").Value.String())
		if err != nil {
			internals.Fatal(err)
		}
	},
}

func Execute() {

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("user-agent", "U", "", "set user agent")
}
