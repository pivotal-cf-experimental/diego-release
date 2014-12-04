package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kr/pty"
	"github.com/pkg/term"

	"github.com/cloudfoundry-incubator/garden/api"
	wclient "github.com/cloudfoundry-incubator/garden/client"
	wconnection "github.com/cloudfoundry-incubator/garden/client/connection"
)

var wardenNetwork = flag.String(
	"wardenNetwork",
	"tcp",
	"warden server network (e.g. unix, tcp)",
)

var wardenAddr = flag.String(
	"wardenAddr",
	"127.0.0.1:7777",
	"warden server address",
)

var create = flag.Bool(
	"create",
	false,
	"create a new container",
)

func main() {
	flag.Parse()

	handle := flag.Arg(0)

	client := wclient.New(wconnection.New(*wardenNetwork, *wardenAddr))

	var container api.Container
	var err error

	if *create {
		container, err = client.Create(api.ContainerSpec{})
	} else if handle == "" {
		containers, err := client.Containers(nil)
		if err != nil {
			log.Fatalln("failed to get containers")
		}

		if len(containers) == 0 {
			log.Fatalln("no containers")
		}

		container = containers[len(containers)-1]
	} else {
		fmt.Println("Hi")
		container, err = client.Lookup(handle)
		fmt.Printf("%+v --- %+v", container, err)
	}

	if err != nil {
		log.Fatalln("failed to lookup container:", err)
	}

	term, err := term.Open(os.Stdin.Name())
	if err != nil {
		log.Fatalln("failed to open terminal:", err)
	}

	err = term.SetRaw()
	if err != nil {
		log.Fatalln("failed to thaw terminal:", err)
	}

	rows, cols, err := pty.Getsize(os.Stdin)
	if err != nil {
		term.Restore()
		log.Fatalln("failed to get window size:", err)
	}

	process, err := container.Run(api.ProcessSpec{
		// Path:       "C:/containerizer/mycontainer/myfile.bat",
		Path:       "myfile.bat",
		Args:       []string{},
		Env:        []string{"TERM=" + os.Getenv("TERM")},
		TTY:        &api.TTYSpec{WindowSize: &api.WindowSize{Columns: cols, Rows: rows}},
		Privileged: true,
	}, api.ProcessIO{
		Stdin:  term,
		Stdout: term,
		Stderr: term,
	})
	if err != nil {
		term.Restore()
		log.Fatalln("failed to run:", err)
	}

	resized := make(chan os.Signal, 10)
	signal.Notify(resized, syscall.SIGWINCH)

	go func() {
		for {
			<-resized

			rows, cols, err := pty.Getsize(os.Stdin)
			if err == nil {
				process.SetTTY(api.TTYSpec{WindowSize: &api.WindowSize{Columns: cols, Rows: rows}})
			}
		}
	}()

	process.Wait()
	term.Restore()
}
