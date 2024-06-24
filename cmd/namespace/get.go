package namespace

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/0xArch3r/opentdf-cli/client"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
	"github.com/spf13/cobra"
)

var ns_id string

var GetNamespaceCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a namespace",
	Run:   getNamespaceHandler,
}

func init() {
	NamespaceCmd.AddCommand(GetNamespaceCmd)
	GetNamespaceCmd.Flags().StringVar(&ns_id, "namespace", "", "The namespace id you want to retrieve")
	GetNamespaceCmd.MarkFlagRequired("namespace")
}

func getNamespaceHandler(cmd *cobra.Command, args []string) {
	s := client.NewSDK(cmd)
	defer s.Close()

	if ns_id == "" {
		slog.Error("the namespace id provided cannot be empty")
		os.Exit(1)
	}

	request := &namespaces.GetNamespaceRequest{
		Id: ns_id,
	}

	response, err := s.SDK().Namespaces.GetNamespace(context.Background(), request)
	if err != nil {
		slog.Error("unable to create namespace", "id", request.Id, "err", err)
		os.Exit(1)
	}

	j, err := json.MarshalIndent(response.Namespace, "", "    ")
	if err != nil {
		slog.Error("unable to marshal namespace response", "id", request.Id, "err", err)
		os.Exit(1)
	}

	fmt.Println(string(j))

}
