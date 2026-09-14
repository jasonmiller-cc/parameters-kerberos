// Command server is the entry point for the parameters-kerberos HTTP service.
package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	corelog "github.com/jasonmiller-cc/parameters-core/pkg/log"
	coreserver "github.com/jasonmiller-cc/parameters-core/pkg/server"
	"github.com/jasonmiller-cc/parameters-kerberos/internal/api"
	"github.com/jasonmiller-cc/parameters-kerberos/internal/config"
	"github.com/jasonmiller-cc/parameters-kerberos/internal/service"
)

func main() {
	cfgPath := flag.String("config", "", "path to config.yaml (optional)")
	flag.Parse()

	log := corelog.Service("parameters-kerberos")

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	svc := service.New(
		cfg.Kerberos.Realm,
		cfg.Kerberos.KAdminServer,
		cfg.Kerberos.KAdminPrincipal,
		cfg.Kerberos.KeytabPath,
		cfg.Kerberos.KDCHost,
	)

	mux := http.NewServeMux()
	handler := api.New(svc)
	handler.Register(mux)

	srv := coreserver.New(cfg.Server, mux, log)
	if err := srv.Run(context.Background()); err != nil {
		log.Error("server error", "error", err)
		os.Exit(1)
	}
}
