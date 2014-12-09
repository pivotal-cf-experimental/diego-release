package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// url := container.containerizerURL.String() + "/api/containers/" + container.Handle() + "/files?destination=" + dstPath
	url := "http://35cd5c7c.ngrok.com:80//api/containers/10514793-1bb3-4106-bd2c-840332126eb5-746d6f4e542847a097650919e7675ae1/files?destination=/app"

	fmt.Printf("Hi 1\n")

	// tarStream, err := os.Open("./main.go.tgz") // For read access.
	// if err != nil {
	// 	log.Fatal(err)
	// }

	tarStream, writer := io.Pipe()
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("Sending Hi\n")
			io.WriteString(writer, "Hi")
			time.Sleep(100 * time.Millisecond)
		}
		writer.Close()
	}()

	fmt.Printf("Hi 2\n")

	req, err := http.NewRequest("PUT", url, tarStream)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hi 3\n")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hi 4\n")
}
