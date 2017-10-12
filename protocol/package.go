package protocol

import (
	"encoding/binary"
)

const (
	PACKAGE_TYPE_HANDSHAKE = 0x01
	PACKAGE_TYPE_ACK = 0x02
	PACKAGE_TYPE_HEARTBEAT = 0x03
	PACKAGE_TYPE_DATA = 0x04
	PACKAGE_TYPE_KICK = 0x05
)

const (
	MESSAGE_TYPE_REQUEST = 0x01
	MESSAGE_TYPE_NOTIFY
	MESSAGE_TYPE_RESPONSE
	MESSAGE_TYPE_PUSH
)

/**
 * Package protocol encode.
 *    1         2             2        2
 * +------+-------------+-------------------------------+
 * | type | body length | messageId | routeId  |  message
 * +------+-------------+--------------------------------+
 *
 * Head: 3bytes
 *   0: package type,
 *      1 - handshake,
 *      2 - handshake ack,
 *      3 - heartbeat,
 *      4 - data
 *      5 - kick
 *   1 - 2: big-endian body length
 * Body: messageId + routeId + message
 *
 */


func GetPackageLength(data []byte) int {
	length := int(binary.BigEndian.Uint16(data[1:3]))
	return 3 + length
}

func Encode(packageType int, body []byte) []byte {
	var data []byte = make([]byte, (3+len(body)))
	copy(data, []byte{byte(packageType)})
	binary.BigEndian.PutUint16(data[1:], uint16(len(body)))
	copy(data[3:], body)

	return data
}

func Decode(data []byte) (packageType int, body []byte) {
	packageType = int(data[0])
	length := GetPackageLength(data)
	body = data[3:length]
	return
}

func MessageEncode(messageId int, routeId int, data[]byte) []byte {
	var body []byte = make([]byte, 4+ len(data))
	binary.BigEndian.PutUint16(body, uint16(messageId))
	binary.BigEndian.PutUint16(body[2:], uint16(routeId))
	copy(body[4:], data)
	return body
}

func MessageDecode(data []byte) (messageId int, routeId int, body []byte) {
	messageId = int(binary.BigEndian.Uint16(data[0:2]))
	routeId = int(binary.BigEndian.Uint16(data[2:4]))
	body = data[4:]
	return
}
