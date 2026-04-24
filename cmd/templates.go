package cmd

import (
	"fmt"
	"strings"

	"github.com/skel-tech/mdpress/internal/auth"
	"github.com/skel-tech/mdpress/internal/cloud"
	tmpl "github.com/skel-tech/mdpress/internal/template"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage templates",
	Long:  `Manage presentation templates - list available templates or pull from the cloud.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default to running list subcommand
		_ = listCmd.RunE(listCmd, args)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Long:  `List all available templates from project, global, and cloud sources.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		localOnly, _ := cmd.Flags().GetBool("local-only")

		// Fetch local templates (project + global)
		templates, err := tmpl.ListTemplates()
		if err != nil {
			return fmt.Errorf("listing templates: %w", err)
		}

		// Fetch cloud templates unless --local-only is set
		var cloudErr error
		if !localOnly {
			cloudTemplates, fetchErr := fetchCloudTemplates()
			if fetchErr != nil {
				cloudErr = fetchErr
			} else {
				templates = append(templates, cloudTemplates...)
			}
		}

		if len(templates) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No templates found.")
			fmt.Fprintln(cmd.OutOrStdout(), "")
			fmt.Fprintln(cmd.OutOrStdout(), "Create templates in:")
			fmt.Fprintf(cmd.OutOrStdout(), "  project: ./templates/\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  global:  %s\n", tmpl.GlobalTemplateDir())
			return nil
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Available templates:")
		for _, t := range templates {
			desc := t.Description
			if desc == "" {
				desc = "(no description)"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  %-12s %-35s (%s)\n", t.Name, desc, t.Source)
		}

		// Show warning if cloud fetch failed
		if cloudErr != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "")
			fmt.Fprintf(cmd.OutOrStderr(), "Warning: could not fetch cloud templates: %v\n", cloudErr)
		}

		return nil
	},
}

// fetchCloudTemplates fetches templates from the cloud API and converts them to TemplateInfo.
func fetchCloudTemplates() ([]tmpl.TemplateInfo, error) {
	client := cloud.NewClient()
	cloudTemplates, err := client.ListTemplates()
	if err != nil {
		return nil, err
	}

	// Check if user is authenticated (has API key configured)
	isAuthenticated := hasAuthCredentials()

	var templates []tmpl.TemplateInfo
	for _, ct := range cloudTemplates {
		desc := ct.Description
		// Append "(requires login)" for non-free templates when unauthenticated
		if !ct.Free && !isAuthenticated {
			desc = strings.TrimSpace(desc + " (requires login)")
		}
		templates = append(templates, tmpl.TemplateInfo{
			Name:        ct.Name,
			Description: desc,
			Source:      tmpl.SourceCloud,
			Free:        ct.Free,
		})
	}

	return templates, nil
}

// hasAuthCredentials checks if the user has auth credentials configured.
func hasAuthCredentials() bool {
	creds, err := auth.Load()
	if err != nil || creds == nil {
		return false
	}
	return strings.TrimSpace(creds.APIKey) != ""
}

func init() {
	listCmd.Flags().Bool("local-only", false, "Only show local templates, skip cloud fetch")
	templatesCmd.AddCommand(listCmd)
	rootCmd.AddCommand(templatesCmd)
}
