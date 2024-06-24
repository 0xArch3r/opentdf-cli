package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/opentdf/platform/sdk"
	"github.com/spf13/cobra"
)

type Sdk struct {
	sdk *sdk.SDK
}

func NewSDK(cmd *cobra.Command) *Sdk {
	sdk_opts := make([]sdk.Option, 0)

	host := cmd.Flag("server").Value.String()
	if host == "" {
		slog.Error("opentdf server is not set")
		os.Exit(1)
	}

	client_id := cmd.Flag("client-id").Value.String()
	if client_id == "" {
		slog.Error("client id is not set")
		os.Exit(1)
	}

	client_secret := cmd.Flag("client-secret").Value.String()
	if client_secret == "" {
		slog.Error("client secret is not set")
		os.Exit(1)
	}
	withClientCreds := sdk.WithClientCredentials(client_id, client_secret, []string{"email"})
	sdk_opts = append(sdk_opts, withClientCreds)

	ignore_insecure, err := cmd.Flags().GetBool("ignore-insecure")
	if err != nil {
		slog.Error("unable to get ignore insecure flag", "err", err)
		os.Exit(1)
	}
	if ignore_insecure {
		sdk_opts = append(sdk_opts, sdk.WithInsecureSkipVerifyConn())
	}

	plaintext, err := cmd.Flags().GetBool("plaintext")
	if err != nil {
		slog.Error("unable to get plaintext flag", "err", err)
		os.Exit(1)
	}
	if plaintext {
		sdk_opts = append(sdk_opts, sdk.WithInsecurePlaintextConn())
	}

	user_cert := cmd.Flag("user-cert").Value.String()
	user_key := cmd.Flag("user-key").Value.String()
	root_ca := cmd.Flag("ca-cert").Value.String()

	if user_cert == "" && user_key != "" {
		slog.Error("user key provided but no cert was provided")
		os.Exit(1)
	}
	if user_cert != "" && user_key == "" {
		slog.Error("user cert provided but no key was provided")
		os.Exit(1)
	}

	if user_cert != "" && user_key != "" {
		rootCAs, _ := x509.SystemCertPool()
		if root_ca != "" {
			f, err := os.ReadFile(root_ca)
			if err != nil {
				slog.Error("unable to open ca cert file", "err", err)
				os.Exit(1)
			}
			rootCAs.AppendCertsFromPEM(f)
		}

		cert, err := tls.LoadX509KeyPair(user_cert, user_key)
		if err != nil {
			slog.Error("failed to load X509 key pair: %v", "err", err)
		}

		sdk_opts = append(sdk_opts, sdk.WithTLSCredentials(&tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      rootCAs,
		},
			[]string{client_id}),
			sdk.WithPlatformConfiguration(sdk.PlatformConfiguration{
				"platform_issuer": "https://localhost:8443/auth/realms/opentdf",
			}),
			sdk.WithTokenEndpoint("https://localhost:8443/auth/realms/opentdf/protocol/openid-connect/token"))
	}

	s, err := sdk.New(host, sdk_opts...)
	if err != nil {
		slog.Error("unable to generate new opentdf sdk client", "err", err)
		os.Exit(1)
	}

	return &Sdk{
		sdk: s,
	}
}

func (s *Sdk) Close() {
	s.sdk.Close()
}

func (s *Sdk) CreateTDF(in string, attrs ...string) (*sdk.TDFObject, string, error) {
	f, err := os.Open(in)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	opts := []sdk.TDFOption{sdk.WithKasInformation(sdk.KASInfo{
		URL:       "https://localhost:8080",
		PublicKey: "",
	})}

	for _, o := range attrs {
		opts = append(opts, sdk.WithDataAttributes(o))
	}

	ciphertext := &strings.Builder{}
	tdf_obj, err := s.sdk.CreateTDF(
		ciphertext,
		f,
		opts...,
	)
	if err != nil {
		return nil, "", err
	}

	return tdf_obj, ciphertext.String(), nil
}

func (s *Sdk) LoadTDF(in string) (*sdk.Reader, string, error) {
	ciphertext, err := os.ReadFile(in)
	if err != nil {
		return nil, "", err
	}

	r, err := s.sdk.LoadTDF(bytes.NewReader(ciphertext))
	if err != nil {
		return nil, "", err
	}

	pt := make([]byte, len(ciphertext))

	_, err = r.Read(pt)
	if err != nil && !errors.Is(err, io.EOF) {
		slog.Error("unable to read ciphertext into plaintext", "err", err)
		os.Exit(1)
	}

	return r, string(pt), nil
}

func (s *Sdk) SDK() *sdk.SDK {
	return s.sdk
}
