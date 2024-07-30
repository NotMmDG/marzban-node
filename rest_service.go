package restservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Service struct {
	connected    bool
	clientIP     string
	sessionID    uuid.UUID
	coreVersion  string
	config       *XRayConfig
	upgrader     websocket.Upgrader
}

type XRayConfig struct {
	// Add fields as necessary
}

// XRayCore and other structs/functions should be implemented based on your specific requirements.

func NewService() *Service {
	return &Service{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (s *Service) base(c echo.Context) error {
	return c.JSON(http.StatusOK, s.response())
}

func (s *Service) connect(c echo.Context) error {
	s.sessionID = uuid.New()
	s.clientIP = c.RealIP()
	s.connected = true
	log.Printf("%s connected, Session ID = \"%s\".", s.clientIP, s.sessionID.String())
	return c.JSON(http.StatusOK, s.response())
}

func (s *Service) disconnect(c echo.Context) error {
	if s.connected {
		log.Printf("%s disconnected, Session ID = \"%s\".", s.clientIP, s.sessionID.String())
	}

	s.sessionID = uuid.UUID{}
	s.clientIP = ""
	s.connected = false
	// Stop the core if necessary
	return c.JSON(http.StatusOK, s.response())
}

func (s *Service) ping(c echo.Context) error {
	var req struct {
		SessionID uuid.UUID `json:"session_id"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.SessionID != s.sessionID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Session ID mismatch."})
	}

	return c.NoContent(http.StatusOK)
}

func (s *Service) start(c echo.Context) error {
	var req struct {
		SessionID uuid.UUID `json:"session_id"`
		Config    string    `json:"config"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.SessionID != s.sessionID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Session ID mismatch."})
	}

	// Start the core with the provided config
	// XRayConfig and core logic should be implemented
	return c.JSON(http.StatusOK, s.response())
}

func (s *Service) response() map[string]interface{} {
	return map[string]interface{}{
		"connected":     s.connected,
		"core_version":  s.coreVersion,
		// Add other relevant fields
	}
}

func (s *Service) logs(c echo.Context) error {
	ws, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// WebSocket logic to handle logs
	// This part depends heavily on how your logging and core management works

	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	service := NewService()
	e.POST("/", service.base)
	e.POST("/connect", service.connect)
	e.POST("/disconnect", service.disconnect)
	e.POST("/ping", service.ping)
	e.POST("/start", service.start)
	e.GET("/logs", service.logs)

	log.Fatal(e.Start(":8080"))
}
