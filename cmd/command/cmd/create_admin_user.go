package cmd

import (
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

var createAdminUserCmd = &cobra.Command{
	Use:   "create_admin_user",
	Short: "Create an admin user",
	Run: func(cmd *cobra.Command, args []string) {
		firstName, _ := cmd.Flags().GetString("first_name")
		lastName, _ := cmd.Flags().GetString("last_name")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		if firstName == "" || lastName == "" || email == "" || password == "" {
			log.Fatal("All parameters (first_name, last_name, email, password) are required")
		}

		fmt.Printf("Creating admin user with First Name: %s, Last Name: %s, Email: %s\n", firstName, lastName, email)

		msg, err := container.AdminUserService.Create(firstName, lastName, email, password)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(msg)
	},
}

func init() {
	rootCmd.AddCommand(createAdminUserCmd)

	// Define flags for the command
	createAdminUserCmd.Flags().String("first_name", "", "First name of the admin user")
	createAdminUserCmd.Flags().String("last_name", "", "Last name of the admin user")
	createAdminUserCmd.Flags().String("email", "", "Email of the admin user")
	createAdminUserCmd.Flags().String("password", "", "Password of the admin user")

	// Mark flags as required
	createAdminUserCmd.MarkFlagRequired("first_name")
	createAdminUserCmd.MarkFlagRequired("last_name")
	createAdminUserCmd.MarkFlagRequired("email")
	createAdminUserCmd.MarkFlagRequired("password")
}
