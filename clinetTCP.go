//Build With: go build -buildmode=c-archive -ldflags "-w -s" clinetTCP.go

package main

import (
    "bufio"
    "io"
    "fmt"
    "net"
    "os"
)

func read(conn net.Conn) {
  reader := bufio.NewReader(conn)
  b := make([]byte, 512)

  for {
    _, err := reader.Read(b)
    if err == io.EOF {
      fmt.Println("Disconnected From Server")
      panic(err)
    } else if err != nil {
      fmt.Println(err)
      panic(err)
    }
    fmt.Println("Echo From Sever: " + string(b))
  }

}

func write(conn net.Conn) {
  writer := bufio.NewWriter(conn)
  stdin := bufio.NewReader(os.Stdin)

  for {
    userInput, _ := stdin.ReadString('\n')
    writer.Write([]byte(userInput))
    writer.Flush()
  }
}

func main() {
    fmt.Println("Where would you like to connect? e.g. 127.0.0.1:6000")
    var ip string
    fmt.Scanln(&ip)

    conn, err := net.Dial("tcp", ip)
    if err != nil {
        fmt.Println(err)
        //os.Exit(1)
    }
    fmt.Println(conn.LocalAddr().String())
    fmt.Println("Type and press enter to send messages!")

    go read(conn)
    go write(conn)

    select{} // Block forever
}
