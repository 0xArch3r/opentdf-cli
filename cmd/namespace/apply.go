package namespace

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/0xArch3r/opentdf-cli/handler"

	"github.com/spf13/cobra"
)

var ns_file string

var applyNamespaceCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a complete namespace configuration through a json file",
	Run:   applyNamespaceHandler,
}

func init() {
	NamespaceCmd.AddCommand(applyNamespaceCmd)
	applyNamespaceCmd.Flags().StringVar(&ns_file, "input", "./namespace.json", "The namespace json you want to apply")
}

func applyNamespaceHandler(cmd *cobra.Command, args []string) {
	s := handler.NewSDK(cmd)
	defer s.Close()

	if ns_file == "" {
		slog.Error("the input file path provided cannot be empty")
		os.Exit(1)
	}

	f, err := os.Open(ns_file)
	if err != nil {
		slog.Error("unable to open input file", "file", ns_file, "err", err)
		os.Exit(1)
	}

	err = s.ApplyNamespace(f)
	if err != nil {
		slog.Error("unable to apply namespace", "err", err)
		os.Exit(1)
	}

	fmt.Println("[+] Successfully applied namespace json")
}
