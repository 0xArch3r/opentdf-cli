package namespace

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/0xArch3r/opentdf-cli/client"
	"github.com/spf13/cobra"
)

var ns string

var CreateNamespaceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a namespace",
	Run:   createNamespaceHandler,
}

func init() {
	NamespaceCmd.AddCommand(CreateNamespaceCmd)
	CreateNamespaceCmd.Flags().StringVar(&ns, "namespace", "", "The namespace you want to create")
	CreateNamespaceCmd.MarkFlagRequired("namespace")
}

func createNamespaceHandler(cmd *cobra.Command, args []string) {
	s := client.NewSDK(cmd)
	defer s.Close()

	if ns == "" {
		slog.Error("the namespace provided cannot be empty")
		os.Exit(1)
	}

	response, err := s.CreateNamespace(ns)
	if err != nil {
		slog.Error("unable to create namespace", "name", ns, "err", err)
		os.Exit(1)
	}

	j, err := json.MarshalIndent(response.Namespace, "", "    ")
	if err != nil {
		slog.Error("unable to marshal namespace response", "name", ns, "err", err)
		os.Exit(1)
	}

	fmt.Println(string(j))

}
