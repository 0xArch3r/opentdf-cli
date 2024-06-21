package attributes

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/0xArch3r/opentdf-cli/handler"

	"github.com/spf13/cobra"
)

var outfile string

var ListAttributesCmd = &cobra.Command{
	Use:   "list",
	Short: "list attributes for a given namespace",
	Long:  `list attributes for a given namespace`,
	Run:   listAttributesHandler,
}

func init() {
	AttributesCmd.Flags().StringVar(&outfile, "output-file", "", "Optional: file to save the output to")
}

func listAttributesHandler(cmd *cobra.Command, args []string) {
	s := handler.NewSDK(cmd)

	attributes, err := s.ListAttributes(ns)
	if err != nil {
		slog.Error("unable to list attributes", "err", err)
		os.Exit(1)
	}

	js, err := json.MarshalIndent(attributes.Attributes, "", "    ")
	if err != nil {
		slog.Error("unable to marshal attributes into json", "err", err)
		os.Exit(1)
	}

	if outfile == "" {
		fmt.Println(string(js))
		return
	}

	o, err := os.Create(outfile)
	if err != nil {
		slog.Error("unable to write attributes to file", "err", err)
		os.Exit(1)
	}
	defer o.Close()

	_, err = o.Write(js)
	if err != nil {
		slog.Error("unable to write attributes to file", "err", err)
		os.Exit(1)
	}
}
