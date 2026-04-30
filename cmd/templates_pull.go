package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/skel-tech/mdpress/internal/auth"
	"github.com/skel-tech/mdpress/internal/cloud"
	tmpl "github.com/skel-tech/mdpress/internal/template"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull <name>",
	Short: "Pull a template from the cloud",
	Long:  `Download a template from the mdpress cloud template library.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		force, _ := cmd.Flags().GetBool("force")
		out := cmd.OutOrStdout()

		client := cloud.NewClient()

		// Fetch template list to check if template exists and if it's free
		cloudTemplates, err := client.ListTemplates()
		if err != nil {
			var netErr *cloud.ErrNetworkFailure
			if errors.As(err, &netErr) {
				return fmt.Errorf("network failure: %w", netErr.Err)
			}
			return err
		}

		// Find the template metadata
		var templateMeta *cloud.CloudTemplate
		for i := range cloudTemplates {
			if cloudTemplates[i].Name == name {
				templateMeta = &cloudTemplates[i]
				break
			}
		}

		if templateMeta == nil {
			return fmt.Errorf("Template '%s' not found in cloud", name)
		}

		// Check if template requires Pro auth
		if !templateMeta.Free {
			if err := auth.RequirePro("templates pull"); err != nil {
				return err
			}
		}

		// Check if local template exists
		if tmpl.TemplateExistsLocal(name) && !force {
			fmt.Fprintf(out, "Template '%s' already exists locally. Overwrite? [y/N]: ", name)

			scanner := bufio.NewScanner(os.Stdin)
			if !scanner.Scan() {
				return fmt.Errorf("Aborted")
			}
			response := strings.TrimSpace(strings.ToLower(scanner.Text()))
			if response != "y" && response != "yes" {
				return fmt.Errorf("Aborted")
			}
		}

		// Fetch template content
		content, err := client.FetchTemplate(name)
		if err != nil {
			var notFoundErr *cloud.ErrTemplateNotFound
			if errors.As(err, &notFoundErr) {
				return fmt.Errorf("Template '%s' not found in cloud", name)
			}
			var netErr *cloud.ErrNetworkFailure
			if errors.As(err, &netErr) {
				return fmt.Errorf("network failure: %w", netErr.Err)
			}
			var unauthErr *cloud.ErrUnauthorized
			if errors.As(err, &unauthErr) {
				return auth.RequirePro("templates pull")
			}
			return err
		}

		// Ensure global template directory exists
		globalDir := tmpl.GlobalTemplateDir()
		if err := os.MkdirAll(globalDir, 0755); err != nil {
			return fmt.Errorf("Failed to save template: %w", err)
		}

		// Write template to file
		destPath := filepath.Join(globalDir, name+".yml")
		if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("Failed to save template: %w", err)
		}

		fmt.Fprintf(out, "Template '%s' saved to %s\n", name, destPath)
		return nil
	},
}

func init() {
	pullCmd.Flags().Bool("force", false, "Overwrite existing template if it exists")
	templatesCmd.AddCommand(pullCmd)
}
