package tcp

import "net"

func handleIncoming(conn net.TCPConn)  {

}

func handleConn(conn net.TCPConn)  {


	for {
		data := make([]byte, 1024)

		_, _ = conn.Read(data)
	}
}
