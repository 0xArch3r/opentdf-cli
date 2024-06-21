/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package encrypt

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
	attr     string
)

// encryptCmd represents the encrypt command
var EncryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Generate an encrypted TDF Object",
	Long:  `Generate an encrypted TDF Object`,
	Run:   encryptHandler,
}

func init() {
	EncryptCmd.PersistentFlags().StringVar(&in_file, "input-file", "./test.tdf", "The encrypted TDF Object file to decrypt")
	EncryptCmd.PersistentFlags().StringVar(&out_file, "output-file", "", "The destination file to save the decrypted TDF Object. If ommitted, output will render to std-out")

	EncryptCmd.PersistentFlags().StringVar(&attr, "attribute", "", "Attribute to apply to the object")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func encryptHandler(cmd *cobra.Command, args []string) {
	if in_file == "" {
		slog.Error("input file is not set")
		os.Exit(1)
	}

	attrs := make([]string, 0)

	if attr != "" {
		attrs = append(attrs, attr)
	}

	s := handler.NewSDK(cmd)
	defer s.Close()

	_, ciphertext, err := s.CreateTDF(in_file, attrs...)
	if err != nil {
		slog.Error("unable to create TDF Object", "err", err)
		os.Exit(1)
	}

	if out_file == "" {
		fmt.Println(ciphertext)
	} else {
		o, err := os.Create(out_file)
		if err != nil {
			slog.Error("unable to write TDF to file", "err", err)
			os.Exit(1)
		}
		defer o.Close()
		_, err = o.Write([]byte(ciphertext))
		if err != nil {
			slog.Error("unable to write TDF to file", "err", err)
			os.Exit(1)
		}
	}
}
