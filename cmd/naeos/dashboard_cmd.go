package main

import (
	"github.com/spf13/cobra"

	"github.com/NAEOS-foundation/naeos/internal/api"
	"github.com/NAEOS-foundation/naeos/internal/dashboard"
	ws "github.com/NAEOS-foundation/naeos/internal/websocket"
)

var (
	dashPort string
)

func newDashboardCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Start NAEOS web dashboard",
		Long:  `Start the NAEOS web dashboard for monitoring and managing projects.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dash, err := dashboard.New()
			if err != nil {
				return err
			}

			wsServer := ws.NewServer()
			wsServer.SetAllowedOrigins([]string{"*"})
			go wsServer.Run()

			mux := api.NewServer(dashPort, &api.AuthConfig{Enabled: false})
			mux.SetWebSocketServer(wsServer)
			mux.Router.HandleFunc("/", dash.ServeHTTP)
			mux.Router.HandleFunc("/ws", wsServer.HandleWebSocket)

			go func() {
				broadcaster := ws.NewEventBroadcaster(wsServer)
				al := dashboard.NewActivityLog(500)
				al.SetLogCallback(func(entry dashboard.LogEntry) {
					level := string(entry.Level)
					if level == "" {
						level = "info"
					}
					broadcaster.LogMessage(level, entry.Message)
				})
				_ = al
				_ = broadcaster
			}()

			return mux.Start()
		},
	}

	cmd.Flags().StringVarP(&dashPort, "port", "p", "3000", "Dashboard port")

	return cmd
}
