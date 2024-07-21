package wshelper

import (
	"github.com/gofiber/websocket/v2"
	"github.com/sirupsen/logrus"
)

var wsClients = make(map[*websocket.Conn]bool)

func NotifyClients(eventType string, data any) {
	for client := range wsClients {
		msg := map[string]any{
			"type": eventType,
			"data": data,
		}
		err := client.WriteJSON(msg)
		logrus.Info("Send %s notify by ws to %s", msg, client.LocalAddr())
		if err != nil {
			logrus.Error("WebSocket error:", err)
			client.Close()
			delete(wsClients, client)
		}
	}
}

func WebSocketHandler(ctx *websocket.Conn) {
	wsClients[ctx] = true
	defer func() {
		delete(wsClients, ctx)
		ctx.Close()
	}()

	for {
		_, _, err := ctx.ReadMessage()
		if err != nil {
			logrus.Error("WebSocket error:", err)
			return
		}
	}
}
