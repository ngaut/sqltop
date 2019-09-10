package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	ui "github.com/gizak/termui/v3"
	flag "github.com/spf13/pflag"

	_ "github.com/go-sql-driver/mysql"
)

const version = "0.1"

var (
	host  = flag.StringP("host", "h", "127.0.0.1", "host")
	pwd   = flag.StringP("password", "p", "", "pwd")
	user  = flag.StringP("user", "u", "root", "user")
	port  = flag.IntP("port", "P", 3306, "port")
	count = flag.IntP("count", "n", 50, "Number of process to show")
)

func InitDB() error {
	globalDS = newDataSource(*user, *pwd, *host, *port)
	if err := globalDS.Connect(); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	if err := InitDB(); err != nil {
		cleanExit(err)
	}
	if err := ui.Init(); err != nil {
		cleanExit(err)
	}
	defer func() {
		ui.Close()
		DB().Close()
	}()

	refreshUI()
}

func cleanExit(err error) {
	ui.Close()
	exec.Command("clear").Run()
	if err != nil {
		log.Print(err)
	}
	os.Exit(0)
}

// refreshUI periodically refreshes the screen.
func refreshUI() {
	controllers := []UIController{
		newProcessListController(),
		newOverviewController(),
		//newHotSpotsController(),
	}
	redraw := make(chan struct{})
	go func() {
		for {
			for _, c := range controllers {
				c.UpdateData()
			}

			redraw <- struct{}{}
			// update every 2 seconds
			time.Sleep(2 * time.Second)
		}
	}()

	evt := ui.PollEvents()
	for {
		select {
		case e := <-evt:
			if e.Type == ui.KeyboardEvent && (e.ID == "q" || e.ID == "<C-c>") {
				cleanExit(nil)
			}
			if e.ID == "<Resize>" {
				for _, c := range controllers {
					c.OnResize(e.Payload.(ui.Resize))
				}
			}

		case <-redraw:
			for _, c := range controllers {
				c.Render()
			}
		}
	}
}
