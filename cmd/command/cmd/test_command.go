package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var customCmd = &cobra.Command{
	Use:   "custom",
	Short: "Run a custom command",
	Run: func(cmd *cobra.Command, args []string) {
		param1, _ := cmd.Flags().GetString("param1")
		param2, _ := cmd.Flags().GetString("param2")
		fmt.Printf("Running custom command with param1=%s and param2=%s\n", param1, param2)
	},
}

func init() {
	customCmd.Flags().String("param1", "", "First parameter")
	customCmd.Flags().String("param2", "", "Second parameter")
}
