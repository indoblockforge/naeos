package main

import (
	"fmt"
	"net/http"
	"github.com/NAEOS-foundation/naeos/internal/websocket"
	"github.com/spf13/cobra"
)

var (
	wsPort string
)

func newWebSocketCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ws",
		Short: "Start WebSocket server for real-time updates",
		Long:  `Start WebSocket server for real-time dashboard updates and event streaming.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			server := websocket.NewServer()
			go server.Run()

			http.HandleFunc("/ws", server.HandleWebSocket)
			http.HandleFunc("/ws/health", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"status":"healthy","clients":%d}`, server.ClientCount())
			})

			fmt.Printf("WebSocket server starting on ws://localhost%s/ws\n", wsPort)
			return http.ListenAndServe(wsPort, nil)
		},
	}

	cmd.Flags().StringVarP(&wsPort, "port", "p", ":8081", "WebSocket server port")

	return cmd
}
