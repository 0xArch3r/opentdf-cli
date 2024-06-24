package namespace

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/0xArch3r/opentdf-cli/client"
	"github.com/spf13/cobra"
)

var out_file string
var namespace string

var exportNamespaceCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the current namespace configuration",
	Run:   exportNamespaceHandler,
}

func init() {
	NamespaceCmd.AddCommand(exportNamespaceCmd)
	exportNamespaceCmd.Flags().StringVar(&namespace, "namespace", "", "Namespace that you would like to export")
	exportNamespaceCmd.MarkFlagRequired("namespace")
	exportNamespaceCmd.Flags().StringVar(&out_file, "output", "", "Optional: File to save exported configuration to")
}

func exportNamespaceHandler(cmd *cobra.Command, args []string) {
	s := client.NewSDK(cmd)
	defer s.Close()

	ns, err := s.ExportNamespace(namespace)
	if err != nil {
		slog.Error("unable to export namespace", "err", err)
		os.Exit(1)
	}

	js, err := json.MarshalIndent(ns, "", "    ")
	if err != nil {
		slog.Error("unable to marshal namespace struct", "err", err)
		os.Exit(1)
	}

	if out_file != "" {
		o, err := os.Create(out_file)
		if err != nil {
			slog.Error("unable to write TDF to file", "err", err)
			os.Exit(1)
		}
		defer o.Close()

		_, err = o.Write(js)
		if err != nil {
			slog.Error("unable to write namespace config to file", "err", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(string(js))
	}
}
