package attributes

import "github.com/spf13/cobra"

var ns string

var AttributesCmd = &cobra.Command{
	Use:   "attributes",
	Short: "List, Set, and Get Attributes values for the platform",
	Long:  `List, Set, and Get Attributes values for the platform`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	AttributesCmd.PersistentFlags().StringVar(&ns, "namespace", "", "The namespace to list attributes for")
	AttributesCmd.AddCommand(ListAttributesCmd)
}
