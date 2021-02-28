package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Hsn723/certspotter-client/api"
	"github.com/Hsn723/ct-exporter/server"
	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
	"github.com/spf13/cobra"
)

const (
	tokenEnv           = "CERTSPOTTER_TOKEN"
	defaultEndpoint    = "https://api.certspotter.com/v1/issuances"
	defaultPositionDir = "/var/log/ct-exporter/"
	defaultListenPort  = 9809
)

var (
	rootCmd = &cobra.Command{
		Use:   "ct-exporter",
		Short: "ct-exporter exports Prometheus metrics for certificate issuances",
		RunE:  runRoot,
	}

	port        uint16
	positionDir string

	// CurrentVersion stores the current version number.
	CurrentVersion string
)

func init() {
	rootCmd.Flags().Uint16VarP(&port, "port", "p", defaultListenPort, "listen port")
	rootCmd.Flags().StringVarP(&positionDir, "position-dir", "d", defaultPositionDir, "path to position file directory")
}

func runRoot(cmd *cobra.Command, args []string) error {
	_ = log.Info("starting ct-exporter", map[string]interface{}{
		"version": CurrentVersion,
	})
	if err := os.MkdirAll(positionDir, 0755); err != nil {
		return err
	}

	token := os.Getenv(tokenEnv)
	if token == "" {
		return fmt.Errorf("must provide %s environment variable", tokenEnv)
	}

	csp := api.CertspotterClient{
		Endpoint: defaultEndpoint,
		Token:    token,
	}

	env := well.NewEnvironment(context.Background())
	server := server.CTExporter{
		Addr:        fmt.Sprintf("0.0.0.0:%d", port),
		Client:      csp,
		Env:         env,
		PositionDir: positionDir,
	}
	if err := server.Start(); err != nil {
		return err
	}

	if err := well.Wait(); err != nil && !well.IsSignaled(err) {
		return err
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.ErrorExit(err)
	}
}
