package aurafx

import (
	"context"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof" // Enables pprof endpoint.
	"strings"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var debugfx = fx.Invoke(NewDebug)

type DebugConfig struct {
	Enabled bool   `json:"enabled" default:"true"`
	Addr    string `json:"addr" default:":6060"`
}

type DebugParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	OnErrorCh chan error

	Log         *zap.SugaredLogger
	DebugConfig *DebugConfig
}

func NewDebug(p DebugParams) error {
	addr, err := net.ResolveTCPAddr("tcp", p.DebugConfig.Addr)
	if err != nil {
		return fmt.Errorf("could not resolve TCP address: %v", err)
	}

	pprofHostPort := addr.String()
	parts := strings.Split(pprofHostPort, ":")
	if len(parts) == 2 && parts[0] == "" {
		pprofHostPort = fmt.Sprintf("localhost:%s", parts[1])
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			if p.DebugConfig.Enabled {
				go func() {
					p.Log.Infof("starting pprof HTTP server at: http://%s/debug/pprof/", pprofHostPort)

					if err := http.ListenAndServe(pprofHostPort, nil); err != nil && err != http.ErrServerClosed {
						p.Log.Errorf("unable to start debug server: %v", err)
						p.OnErrorCh <- err
					}
				}()
			}
			return nil
		},
	})

	return nil
}
