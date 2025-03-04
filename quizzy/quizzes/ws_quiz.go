package quizzes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"log"
)

func configureWs(router *gin.RouterGroup) {
	server := socketio.NewServer(nil)
	server.OnConnect("/", onConnect)
	server.OnEvent("/", "host", onHostEvent)

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket.IO server error: %v", err)
		}
	}()

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
}

func onConnect(s socketio.Conn) error {
	fmt.Println(" Client connecté:", s.ID())
	return nil
}

type hostEventPayload struct {
	ExecutionId string `json:"executionId"`
}
type hostEventResponse struct {
	Quiz Quiz `json:"quiz"`
}

func onHostEvent(s socketio.Conn, msg string) string {
	fmt.Println(" Reçu event 'host' avec message:", msg)

	var payload hostEventPayload
	if err := json.Unmarshal([]byte(msg), &payload); err != nil {
		fmt.Println(" Erreur Unmarshal:", err)
		return ""
	}

	fmt.Println("ExecutionId reçu:", payload.ExecutionId)

	// Simuler une réponse
	response := hostEventResponse{Quiz: Quiz{}}
	res, _ := json.Marshal(response)
	s.Emit("hostDetails", string(res))

	return string(res)
}
