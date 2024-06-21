/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/0xArch3r/opentdf-cli/cmd/attributes"
	"github.com/0xArch3r/opentdf-cli/cmd/decrypt"
	"github.com/0xArch3r/opentdf-cli/cmd/encrypt"
	"github.com/0xArch3r/opentdf-cli/cmd/namespace"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opentdf-cli",
	Short: "Custom CLI to interact with Open TDF Platform",
	Long:  `Custom CLI to interact with Open TDF Platform`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.opentdf-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringP("server", "s", "https://localhost:8080", "Open TDF host to make the requests to")
	rootCmd.PersistentFlags().StringP("client-id", "c", "opentdf-sdk", "Client ID to authenticate against")
	rootCmd.PersistentFlags().StringP("client-secret", "p", "secret", "Client Secret for given Client ID")
	rootCmd.PersistentFlags().StringP("user-cert", "u", "", "User Certificate to exchange")
	rootCmd.PersistentFlags().StringP("user-key", "k", "", "User Key for certificate")
	rootCmd.PersistentFlags().StringP("ca-cert", "r", "", "Root CA certificate for user cert")
	rootCmd.PersistentFlags().BoolP("ignore-insecure", "i", false, "Ignore any TLS insecure issues")
	rootCmd.PersistentFlags().BoolP("plaintext", "t", false, "Use plaintext protocols, disabling TLS")

	rootCmd.AddCommand(encrypt.EncryptCmd)
	rootCmd.AddCommand(decrypt.DecryptCmd)
	rootCmd.AddCommand(namespace.NamespaceCmd)
	rootCmd.AddCommand(attributes.AttributesCmd)

}
