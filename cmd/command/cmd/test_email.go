package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testEmailCmd = &cobra.Command{
	Use:   "test_email",
	Short: "Send a test email",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Sending test email\n")
	},
}
