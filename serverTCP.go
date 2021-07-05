//Build With: go build -buildmode=c-archive -ldflags "-w -s" serverTCP.go
package main

import "C"
import (
  "fmt"
  "bufio"
  "io"
	"net"
  "sync"
)

var clients sync.Map   //A sync map to store all clients info so we can brodcast
                       //a message from one client to all the others

//export SayHello
func SayHello(){
  for i := 0; i<100; i++ {
    fmt.Println("Yo World What's Up?")
  }
}

//Handles communications between each client and server
func communications (conn net.Conn) {
  defer conn.Close()

  //So we can plop and edit data in a buffer when reading or writing to the
  //TCP connection
  connReader := bufio.NewReader(conn)

  for {
    b := make([]byte, 512)
    _, err := connReader.Read(b)
    if err == io.EOF {
      fmt.Println("Client Disconnected")
      clients.Delete(conn.RemoteAddr())
      break
    } else if err != nil {
      fmt.Println(err)
      break
    }

    fmt.Println(b)

    // Now brodcast message to all other users
    clients.Range(func(key, value interface{}) bool {
      if key != conn.RemoteAddr() {
        connWriter := bufio.NewWriter(value.(net.Conn))
        connWriter.Write(b)
        connWriter.Flush() // Flush writes any buffered data to the underlying io.Writer
                           // (in this case the connection to the client)
      }
      return true //tells range to keep goin
    })
  }
}

//export TCPListener
func TCPListener() {
  // Bind to TCP port 20080 on all interfaces.
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
	  fmt.Println("Unable to bind to port")
	}
	fmt.Println("Listening on 0.0.0.0:20080")
	for {
		// Wait for connection. Create net.Conn on connection established.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Unable to accept connection")
		}
    fmt.Println("Received connection from: " + conn.RemoteAddr().String())

    //add client to the sync.Map
    clients.Store(conn.RemoteAddr(), conn)
		// Handle the connection. Using goroutine for concurrency.
		go communications(conn)

	}

}

func main(){
  TCPListener()
}
