package aurafx

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"

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
	AdminConfig AdminConfig

	AdminHandlers []map[string]http.Handler `group:"admin_api_handler"`
}

func NewAdmin(p AdminParams) error {
	addr, err := net.ResolveTCPAddr("tcp", p.AdminConfig.Addr)
	if err != nil {
		return fmt.Errorf("could not resolve TCP address: %v", err)
	}

	hostPort := addr.String()

	adminMux := http.NewServeMux()

	// register admin handlers
	for _, handlers := range p.AdminHandlers {
		for k, v := range handlers {
			p.Log.Debugf("register admin handler: %s", k)
			adminMux.Handle(k, v)
		}
	}

	// index responds with an HTML page listing the available admin handlers
	adminMux.HandleFunc("/", adminIndexHandler(p.Log, p.AdminHandlers))

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

type adminEntry struct {
	Name string
	Href string
}

func adminIndexHandler(
	log *zap.SugaredLogger,
	adminHandlers []map[string]http.Handler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		var entries []adminEntry
		for _, handlers := range adminHandlers {
			for k := range handlers {
				entries = append(entries, adminEntry{
					Name: k,
					Href: k,
				})
			}
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name < entries[j].Name
		})

		if err := adminIndexTmplExecute(w, entries); err != nil {
			log.Error(err)
		}
	}
}

func adminIndexTmplExecute(w io.Writer, entries []adminEntry) error {
	var b bytes.Buffer
	b.WriteString(`<html>
<head>
<title>admin</title>
</head>
<body>
admin
<br/>
<p>
<ul>
`)

	for _, entry := range entries {
		link := &url.URL{Path: entry.Href}
		fmt.Fprintf(&b, "<li><a href='%s'>%s</a></li>\n", link, html.EscapeString(entry.Name))
	}

	b.WriteString(`</ul>
</p>
</body>
</html>`)

	_, err := w.Write(b.Bytes())
	return err
}
