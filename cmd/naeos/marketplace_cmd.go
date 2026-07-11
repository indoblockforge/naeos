package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/NAEOS-foundation/naeos/internal/marketplace"
)

func newMarketplaceCommand() *cobra.Command {
	var cacheDir string

	cmd := &cobra.Command{
		Use:   "marketplace",
		Short: "Browse and install templates, profiles, and plugins",
		Long: `NAEOS Marketplace for templates, profiles, and plugins.

Example:
  naeos marketplace search "web api"
  naeos marketplace install web-api-template
  naeos marketplace profile list
  naeos marketplace plugin list`,
	}

	// Template commands
	searchCmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for templates",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := marketplace.NewClient(cacheDir)
			query := strings.Join(args, " ")
			results, err := client.Search(marketplace.SearchFilter{Query: query, Limit: 10})
			if err != nil {
				return err
			}
			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No results found")
				return nil
			}
			for _, r := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%-25s %-10s %s\n", r.Name, r.Version, r.Description)
			}
			return nil
		},
	}

	installCmd := &cobra.Command{
		Use:   "install [name]",
		Short: "Install a template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := marketplace.NewClient(cacheDir)
			if err := client.Install(args[0], "."); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Installed template %s\n", args[0])
			return nil
		},
	}

	// Profile commands
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage marketplace profiles",
	}

	profileListCmd := &cobra.Command{
		Use:   "list",
		Short: "List available profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			pm := marketplace.NewProfileMarketplace(cacheDir)
			profiles, err := pm.List()
			if err != nil {
				return err
			}
			for _, p := range profiles {
				fmt.Fprintf(cmd.OutOrStdout(), "%-25s %-10s %-15s %s\n", p.Name, p.Version, p.Industry, p.Description)
			}
			return nil
		},
	}

	profileSearchCmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search profiles",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pm := marketplace.NewProfileMarketplace(cacheDir)
			query := strings.Join(args, " ")
			results, err := pm.Search(query, nil)
			if err != nil {
				return err
			}
			for _, p := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%-25s %-10s %-15s %s\n", p.Name, p.Version, p.Industry, p.Description)
			}
			return nil
		},
	}

	profileDownloadCmd := &cobra.Command{
		Use:   "download [name]",
		Short: "Download a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pm := marketplace.NewProfileMarketplace(cacheDir)
			if err := pm.Download(args[0], "."); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Downloaded profile %s\n", args[0])
			return nil
		},
	}

	profilePublishCmd := &cobra.Command{
		Use:   "publish [file]",
		Short: "Publish a profile from JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pm := marketplace.NewProfileMarketplace(cacheDir)
			if err := pm.Upload(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Published profile from %s\n", args[0])
			return nil
		},
	}

	profileCmd.AddCommand(profileListCmd, profileSearchCmd, profileDownloadCmd, profilePublishCmd)

	// Plugin commands
	pluginCmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage marketplace plugins",
	}

	pluginListCmd := &cobra.Command{
		Use:   "list",
		Short: "List available plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			pm := marketplace.NewPluginMarketplace(cacheDir, filepath.Join(".", ".naeos", "plugins"))
			plugins, err := pm.List()
			if err != nil {
				return err
			}
			for _, p := range plugins {
				status := ""
				if p.Installed {
					status = " [installed]"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%-25s %-10s %-15s %s%s\n", p.Name, p.Version, p.Type, p.Description, status)
			}
			return nil
		},
	}

	pluginInstallCmd := &cobra.Command{
		Use:   "install [name]",
		Short: "Install a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pm := marketplace.NewPluginMarketplace(cacheDir, filepath.Join(".", ".naeos", "plugins"))
			if err := pm.Install(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Installed plugin %s\n", args[0])
			return nil
		},
	}

	pluginUninstallCmd := &cobra.Command{
		Use:   "uninstall [name]",
		Short: "Uninstall a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pm := marketplace.NewPluginMarketplace(cacheDir, filepath.Join(".", ".naeos", "plugins"))
			if err := pm.Uninstall(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Uninstalled plugin %s\n", args[0])
			return nil
		},
	}

	pluginSearchCmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search plugins",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pm := marketplace.NewPluginMarketplace(cacheDir, filepath.Join(".", ".naeos", "plugins"))
			query := strings.Join(args, " ")
			results, err := pm.Search(query, nil)
			if err != nil {
				return err
			}
			for _, p := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%-25s %-10s %-15s %s\n", p.Name, p.Version, p.Type, p.Description)
			}
			return nil
		},
	}

	pluginCmd.AddCommand(pluginListCmd, pluginInstallCmd, pluginUninstallCmd, pluginSearchCmd)

	cmd.AddCommand(searchCmd, installCmd, profileCmd, pluginCmd)
	cmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", filepath.Join(".", ".naeos", "cache"), "cache directory")
	return cmd
}
