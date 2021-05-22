package aurafx

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var adminfx = fx.Invoke(NewAdmin)

type AdminConfig struct {
	Enabled bool   `json:"enabled" default:"true"`
	Addr    string `json:"addr" default:":8081"`
}

type AdminHandlerResult struct {
	fx.Out

	AdminHandlers map[string]http.Handler `group:"admin_api_handler"`
}

type AdminParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	OnErrorCh chan error

	Log         *zap.SugaredLogger
	AdminConfig *AdminConfig

	AdminHandlers []map[string]http.Handler `group:"admin_api_handler"`
}

func NewAdmin(p AdminParams) error {
	addr, err := net.ResolveTCPAddr("tcp", p.AdminConfig.Addr)
	if err != nil {
		return fmt.Errorf("could not resolve TCP address: %v", err)
	}

	hostPort := addr.String()
	parts := strings.Split(hostPort, ":")
	if len(parts) == 2 && parts[0] == "" {
		hostPort = fmt.Sprintf("localhost:%s", parts[1])
	}

	adminMux := http.NewServeMux()

	// register admin handlers
	for _, handlers := range p.AdminHandlers {
		for k, v := range handlers {
			p.Log.Debugf("register admin handler: %s", k)
			adminMux.Handle(k, v)
		}
	}

	server := &http.Server{
		Addr:    hostPort,
		Handler: adminMux,
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			if p.AdminConfig.Enabled {
				go func() {
					p.Log.Infof("starting admin HTTP server at: http://%s/", hostPort)

					if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						p.Log.Errorf("unable to start admin server: %v", err)
						p.OnErrorCh <- err
					}
				}()
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Log.Infof("stopping admin server")
			return server.Shutdown(ctx)
		},
	})

	return nil
}
