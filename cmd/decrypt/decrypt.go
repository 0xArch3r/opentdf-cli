/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package decrypt

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/0xArch3r/opentdf-cli/handler"

	"github.com/spf13/cobra"
)

var (
	in_file  string
	out_file string
)

// encryptCmd represents the encrypt command
var DecryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "decrypt a TDF Object",
	Long:  `decrypt a TDF Object`,
	Run:   handleDecrypt,
}

func init() {
	DecryptCmd.PersistentFlags().StringVar(&in_file, "input-file", "./test.tdf", "The encrypted TDF Object file to decrypt")
	DecryptCmd.PersistentFlags().StringVar(&out_file, "output-file", "", "The destination file to save the decrypted TDF Object. If ommitted, output will render to std-out")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func handleDecrypt(cmd *cobra.Command, args []string) {
	if in_file == "" {
		slog.Error("input file is not set")
		os.Exit(1)
	}

	s := handler.NewSDK(cmd)
	defer s.Close()

	_, plaintext, err := s.LoadTDF(in_file)
	if err != nil {
		slog.Error("unable to decrypt TDF", "err", err)
		os.Exit(1)
	}

	if out_file == "" {
		fmt.Println(plaintext)
	} else {
		o, err := os.Create(out_file)
		if err != nil {
			slog.Error("unable to write plaintext TDF to file", "err", err)
			os.Exit(1)
		}
		defer o.Close()
		_, err = o.Write([]byte(plaintext))
		if err != nil {
			slog.Error("unable to write TDF to file", "err", err)
			os.Exit(1)
		}
	}
}
