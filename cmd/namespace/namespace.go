package namespace

import "github.com/spf13/cobra"

var NamespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "List, Set, and Get Namespace values for the platform",
	Long:  `List, Set, and Get Namespace values for the platform`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// EncryptCmd.PersistentFlags().StringVar(&in_file, "input-file", "./test.tdf", "The encrypted TDF Object file to decrypt")
	// EncryptCmd.PersistentFlags().StringVar(&out_file, "output-file", "", "The destination file to save the decrypted TDF Object. If ommitted, output will render to std-out")

	// EncryptCmd.PersistentFlags().StringVar(&attr, "attribute", "", "Attribute to apply to the object")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
