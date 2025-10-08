package engine

import (
	"ember/globals"
	"fmt"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var SERVER *http.Server
var PORT = 8080

func StartDevEngine(config *globals.GameConfig, runButton *widget.Button) error {
	var mux = http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		var res, err = GenerateSourceFromConfig(config)
		if err != nil {
			return
		}

		fmt.Fprint(w, res)
	})

	SERVER = &http.Server{Addr: fmt.Sprintf(":%v", PORT), Handler: mux}
	go func() error {
		if err := SERVER.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	}()

	fyne.Do(func() {
		fyne.CurrentApp().SendNotification(fyne.NewNotification("New Process", fmt.Sprintf("Game is running, visit http:/127.0.0.1:%v", PORT)))
		runButton.SetText("Stop")
	})

	return nil
}
func StopDevEngine(runButton *widget.Button) {
	if SERVER != nil {
		SERVER.Close()
	}
	fyne.Do(func() {
		runButton.SetText("Run Game")
	})
}

