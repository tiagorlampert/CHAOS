package client

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/tiagorlampert/CHAOS/internal/utils"
	"github.com/tiagorlampert/CHAOS/internal/utils/system"
	"sync"
)

// Author: Viru
// Email: viru@gmail.com

type Service struct {
	connections map[string]*websocket.Conn
	mu          sync.Mutex
}

func (s *Service) AddConnection(clientID string, conn *websocket.Conn) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[clientID] = conn
	return nil
}

func (s *Service) GetConnection(clientID string) (*websocket.Conn, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn, exists := s.connections[clientID]
	return conn, exists
}

func (s *Service) RemoveConnection(clientID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if conn, exists := s.connections[clientID]; exists {
		err := conn.Close()
		delete(s.connections, clientID)
		return err
	}
	return errors.New("connection not found")
}

func (s *Service) SendCommand(ctx context.Context, input SendCommandInput) (SendCommandOutput, error) {
	conn, exists := s.GetConnection(input.ClientID)
	if !exists {
		return SendCommandOutput{}, errors.New("client not connected")
	}
	defer func() {
		if err := conn.Close(); err != nil {
			// Log error if needed
		}
	}()

	if err := conn.WriteMessage(websocket.TextMessage, []byte(input.Command)); err != nil {
		s.RemoveConnection(input.ClientID)
		return SendCommandOutput{}, err
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		s.RemoveConnection(input.ClientID)
		return SendCommandOutput{}, err
	}

	return SendCommandOutput{Response: string(message)}, nil
}
