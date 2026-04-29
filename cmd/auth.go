package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/skel-tech/mdpress/internal/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  "Log in, log out, and check your mdpress Pro authentication status.",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to your mdpress Pro account",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()

		fmt.Fprintln(out, "Log in to mdpress Pro")
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, "Get your API key from: https://mdpress.app/settings/api-keys")
		fmt.Fprintln(out, "")
		fmt.Fprint(out, "Paste your API key: ")

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return fmt.Errorf("no input received")
		}
		key := strings.TrimSpace(scanner.Text())

		if key == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		if !strings.HasPrefix(key, "mdp_") {
			return fmt.Errorf("invalid API key format (expected mdp_...)")
		}

		// Prompt for email (optional, for whoami display)
		fmt.Fprint(out, "Email (optional): ")
		var email string
		if scanner.Scan() {
			email = strings.TrimSpace(scanner.Text())
		}

		creds := &auth.Credentials{
			APIKey: key,
			Email:  email,
		}

		if err := auth.Save(creds); err != nil {
			return fmt.Errorf("saving credentials: %w", err)
		}

		if email != "" {
			fmt.Fprintf(out, "\nLogged in as %s\n", email)
		} else {
			fmt.Fprintln(out, "\nLogged in successfully")
		}

		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of your mdpress Pro account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := auth.Clear(); err != nil {
			return fmt.Errorf("clearing credentials: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Logged out successfully")
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show the currently authenticated user",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()

		creds, err := auth.Load()
		if err != nil {
			return fmt.Errorf("reading credentials: %w", err)
		}

		if creds == nil || creds.APIKey == "" {
			fmt.Fprintln(out, "Not logged in. Run: mdpress auth login")
			return nil
		}

		if creds.Email != "" {
			fmt.Fprintf(out, "Email:   %s\n", creds.Email)
		}
		fmt.Fprintf(out, "API Key: %s\n", maskKey(creds.APIKey))
		return nil
	},
}

// maskKey masks an API key, showing only the prefix and last 4 characters.
func maskKey(key string) string {
	if len(key) <= 8 {
		return "mdp_****"
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(authCmd)
}
