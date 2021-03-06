package tcp

import (
	"net"
	"github.com/pengzj/swift/hub"
	"bytes"
	"github.com/pengzj/swift/protocol"
	"time"
	"github.com/pengzj/swift/logger"
)

var (
	heartbeatInterval = 10 * time.Second
)

type TcpSocket struct {
	CloseChan chan bool
}

func (socket *TcpSocket) Start(host ,port string)  {
	tcpAddr, err := net.ResolveTCPAddr("tcp", host + ":" + port)
	if err != nil {
		logger.Fatal(err)
	}
	listener, err := net.Listen("tcp", tcpAddr.String())
	if err != nil {
		logger.Fatal(err)
	}
	defer listener.Close()

	for {
		select {
		case <- socket.CloseChan:
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				logger.Fatal(err)
			}

			session := &hub.Session{
				Id: hub.UniqueId(),
				Conn: conn,
				Send: make(chan []byte, 10),
			}

			hub.GetHub().Register <- session

			go readPump(session)
			go writePump(session)
		}

	}
}

func readPump(session *hub.Session)  {
	defer session.Close()
	conn := session.Conn.(net.Conn)
	var buffer bytes.Buffer
	var headerLength = protocol.GetHeadLength()
	var currentTotalLength int
	var length int


	for {
		conn.SetReadDeadline(time.Now().Add(heartbeatInterval))

		tmp := make([]byte, 512)
		n, err := conn.Read(tmp)
		if err != nil {
			logger.Error(err)
			return
		}

		buf := make([]byte, n)
		copy(buf, tmp[0:n])

		buffer.Write(buf)


		//do with packet splicing
		for {
			currentTotalLength = len(buffer.Bytes())
			length = headerLength +  protocol.GetBodyLength(buffer.Bytes())
			message := make([]byte, length)

			if length > currentTotalLength {
				break
			}

			_, err := buffer.Read(message)
			if err != nil {
				logger.Fatal("read data error: ", err)
			}

			session.HandleData(message)

			leftLength := currentTotalLength - length
			if leftLength > 0 {
				leftData := make([]byte, leftLength)
				_, err = buffer.Read(leftData)
				if err != nil {
					logger.Fatal("package data error: %v", err)
				}
				buffer.Reset()
				buffer.Write(leftData)
			} else {
				buffer.Reset()
				break
			}
		}
	}
}

func writePump(session *hub.Session)  {
	defer session.Close()
	conn := session.Conn.(net.Conn)
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()
	var buffer bytes.Buffer
	for {
		select {
		case message := <-session.Send:
			buffer.Write(message)
			n := len(session.Send)

			for i := 0; i < n; i++ {
				buffer.Write(<-session.Send)
			}

			conn.SetWriteDeadline(time.Now().Add(heartbeatInterval))
			_, err := conn.Write(buffer.Bytes())
			if err != nil {
   				return
			}
			buffer.Reset()
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(heartbeatInterval))
			_, err := conn.Write(protocol.Encode(protocol.TYPE_HEARTBEAT, []byte{}))
			if err != nil {
				return
			}
		}
	}
}


func (socket *TcpSocket) Close()  {
	socket.CloseChan <- true
}