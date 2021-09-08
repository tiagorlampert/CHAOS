package environment

import (
	"fmt"
	"strings"
)

func LoadConfiguration(serverAddress, serverPort, token string) *Configuration {
	return &Configuration{
		Connection: Connection{
			Token:             fmt.Sprint("jwt=", token),
			ContextDeadline:   5,
			ContentTypeHeader: "Content-Type",
			ContentTypeJSON:   "application/json",
			CookieHeader:      "Cookie",
		},
		Server: Server{
			Address: serverAddress,
			Port:    serverPort,
			URL:     newServerURL(serverAddress, serverPort),
			Endpoint: Endpoint{
				Health:   "health",
				Device:   "device",
				Command:  "command",
				Upload:   "upload",
				Download: "download",
			},
		},
		CommandHandler: CommandHandler{
			CommandFileExplorer: "explore",
			CommandDownload:     "download",
			CommandUpload:       "upload",
		},
	}
}

func newServerURL(serverAddress, serverPort string) string {
	if len(strings.TrimSpace(serverPort)) == 0 {
		return fmt.Sprintf("%s/", strings.TrimRight(serverAddress, "/"))
	}
	return fmt.Sprintf("http://%s:%s/", serverAddress, serverPort)
}
