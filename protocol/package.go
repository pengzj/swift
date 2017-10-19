package protocol

import (
	"encoding/binary"
	"math"
)

const (
	TYPE_HANDSHAKE = 0x01
	TYPE_HANDSHAKE_ACK = 0x02
	TYPE_HEARTBEAT = 0x03
	TYPE_KICK = 0x04
	//real data
	TYPE_DATA_REQUEST = 0x05
	TYPE_DATA_RESPONSE = 0x06
	TYPE_DATA_NOTIFY = 0x07
	TYPE_DATA_PUSH = 0x08
)

/**
 * Package protocol encode.
 *    1             3                         n
 * +------+------------------------+-------------------------------+
 * | type |       body length      |           body
 * +------+------------------------+--------------------------------+
 *
 * Head: 3bytes
 *   0: package type,
 *      1 - handshake
 *      2 - handshake ack,
 *      3 - heartbeat,
 *      4 - kick
 *      5 - request
 *      6 - response
 *      7 - notify
 *      8 - push
 *   1 - 3: big-endian body length, support data up to 16MB
 * Body: messageId + routeId + message
 * +----------------+----------------+------------------+
 * |    messageId   |  routeId       |     message      |
 * +----------------+----------------+-------------------+
 *
 * messageId (1 ~ 10 bytes)
 * routeId (1 ~ 10 bytes)
 * data
 */


func GetHeadLength() int  {
	return 4
}

func GetBodyLength(data []byte) int {
	lenBytes := make([]byte, 4)
	lenBytes[0] = 0x0
	copy(lenBytes[1:], data[1:4])

	length := int(binary.BigEndian.Uint32(lenBytes))
	return length
}

func Encode(packageType int, body []byte) []byte {
	buffer := make([]byte, 4+len(body))

	//write package type
	copy(buffer, []byte{byte(packageType)})

	//write data length
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(body)))
	// only save the last 3 bytes
	copy(buffer[1:], lenBytes[1:])

	//write data
	copy(buffer[4:], body)


	return buffer
}

func Decode(data []byte) (packageType int, body []byte) {
	packageType = int(data[0])
	length := GetBodyLength(data)
	body = data[4:(4+length)]

	return
}

func MessageEncode(messageId int, routeId int, data[]byte) []byte {
	idLen := GetNumberBytes(messageId)
	routeLen := GetNumberBytes(routeId)
	var body []byte = make([]byte, idLen + routeLen + len(data))

	encodeNum(messageId, body, 0)
	encodeNum(routeId, body, idLen)

	copy(body[(idLen+routeLen):], data)
	return body
}

func MessageDecode(data []byte) (messageId int, routeId int, body []byte) {
	var offset = 0
	var count int

	//read messageId
	for {
		count++
		offset++
		if data[(offset-1)] < 0x80 {
			break
		}

	}
	var tp, idx, mid int
	for i := 0; i < count; i++ {
		idx = offset-(count-i)
		tp = int(data[idx])

		mid = int((tp & 0x7f) )* int(math.Pow(float64(2), float64(7 * (count-i-1))))
		messageId = messageId + mid
	}

	//read routeId
	count = 0
	for {
		count++
		offset++
		if data[(offset-1)] < 0x80 {
			break
		}


	}
	for i := 0; i < count; i++ {
		idx = offset-(count-i)
		tp = int(data[idx])
		mid = int(tp & 0x7f) * int(math.Pow(float64(2), float64(7 * (count-i-1))))
		routeId = routeId + mid
	}

	body = data[offset:]
	return
}

func GetNumberBytes(num int) int {
	var length = 0
	for {
		length += 1
		num >>= 7
		if num <= 0 {
			break
		}
	}
	return length
}

func encodeNum(num int, buffer []byte, offset int)  {
	var left , right int
	var temp = make([]int, binary.MaxVarintLen64)
	var count = 0
	for {
		left = int(num%128)
		right =  int(math.Floor(float64(num/128)))

		temp[count] = left | 0x80
		count++

		num = right
		if num == 0 {
			break
		}
	}
	temp[0] = 0x7f & temp[0]


	for i := 0; i < count; i++ {
		buffer[offset + (count -i-1)] = byte(temp[i])
	}
}
