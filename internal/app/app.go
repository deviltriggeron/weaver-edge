package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"weaver-edge/config"
	"weaver-edge/internal/application"
	"weaver-edge/internal/infrastructure/connectors/telegram"
	"weaver-edge/internal/infrastructure/core"
	"weaver-edge/internal/infrastructure/db"
	"weaver-edge/internal/infrastructure/db/repository"
	"weaver-edge/internal/ports"
	router "weaver-edge/internal/transport/http"
	"weaver-edge/internal/transport/http/handler"
)

// TODO: db closer
type connectorClients struct {
	tg ports.EventSource
}

type repositories struct {
	db      *sql.DB
	context ports.ContextRepository
	audit   ports.AuditRepository
	policy  ports.PolicyRepository
}

type services struct {
	contextSrv ports.ContextService
	auditSrv   ports.AuditService
	policySrv  ports.PolicyEngine
}

type servers struct {
	http *http.Server
}

func initConnectors(cfg config.Config) connectorClients {
	client := telegram.NewTelegramClient(cfg)
	tg := telegram.NewTelegramConnector(client)
	return connectorClients{
		tg: tg,
	}
}

func initRepositories(cfg config.Config) (repositories, error) {
	db, err := db.InitDB(cfg)
	if err != nil {
		return repositories{}, err
	}

	contextRepo := repository.NewContextRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	policyRepo := repository.NewPolicyRepo(db)

	return repositories{
		db:      db,
		context: contextRepo,
		audit:   auditRepo,
		policy:  policyRepo,
	}, nil
}

func initServices(ctx context.Context, repos repositories) (services, error) {
	contextService := application.NewContextService(repos.context)
	auditService := application.NewAuditService(repos.audit)

	policyServices, err := application.NewPolicyService(ctx, repos.policy)
	if err != nil {
		return services{}, err
	}

	return services{
		contextSrv: contextService,
		auditSrv:   auditService,
		policySrv:  policyServices,
	}, nil
}

func initUsecase(
	connectorClients connectorClients,
	services services,
	core ports.CoreClient,
) application.Application {
	return application.NewApplication(
		connectorClients.tg,
		services.policySrv,
		core,
		services.contextSrv,
		services.auditSrv,
	)
}

func initCoreClient(cfg config.Config) ports.CoreClient {
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}
	core := core.NewCoreClient(httpClient, cfg)

	return core
}

func initServers(cfg config.Config, handler http.Handler) servers {
	log.Println(cfg.Addr.HttpAddr)
	srv := &http.Server{
		Addr:    cfg.Addr.HttpAddr,
		Handler: handler,
	}

	return servers{
		http: srv,
	}
}

func runHTTPServer(ctx context.Context, srv *http.Server) error {
	errCh := make(chan error, 1)

	go func() {
		log.Printf("api server started on %s", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("shutting down api server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			_ = srv.Close()
			return fmt.Errorf("shutdown http server: %w", err)
		}

		return nil

	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("http server: %w", err)
	}
}

func Run(cfg config.Config) error {
	ctx := context.Background()

	connectors := initConnectors(cfg)

	repositories, err := initRepositories(cfg)
	if err != nil {
		return fmt.Errorf("initialize repositories: %w", err)
	}
	defer repositories.db.Close()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	services, err := initServices(ctx, repositories)
	if err != nil {
		return fmt.Errorf("initialize services: %w", err)
	}
	core := initCoreClient(cfg)

	application := initUsecase(connectors, services, core)

	handler := handler.NewHandler(application)
	if cfg.Addr.ContextToken == "" {
		return fmt.Errorf("context token empty")
	}
	router := router.NewRouter(&handler, cfg.Addr.ContextToken)

	srv := initServers(cfg, router)

	if err := runHTTPServer(ctx, srv.http); err != nil {
		return fmt.Errorf("run application server: %w", err)
	}

	return nil
}
