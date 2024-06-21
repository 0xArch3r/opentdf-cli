package namespace

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/0xArch3r/opentdf-cli/handler"

	"github.com/spf13/cobra"
)

var listNamespaceCmd = &cobra.Command{
	Use:   "list",
	Short: "List all namespaces",
	Run:   listNamespaceHandler,
}

func init() {
	NamespaceCmd.AddCommand(listNamespaceCmd)
}

func listNamespaceHandler(cmd *cobra.Command, args []string) {
	s := handler.NewSDK(cmd)
	defer s.Close()

	response, err := s.ListNamespaces()
	if err != nil {
		slog.Error("unable to list namespaces", "err", err)
		os.Exit(1)
	}

	j, err := json.MarshalIndent(response.Namespaces, "", "    ")
	if err != nil {
		slog.Error("unable to marshal namespaces response", "err", err)
		os.Exit(1)
	}

	fmt.Println(string(j))

}
