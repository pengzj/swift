package protocol

import (
	"testing"
	"bytes"
	"strings"
)

var (
	data []byte
	body []byte
	msgId int
	routeId int
	encode []byte
	packageType int
)

func TestMessageEncode(t *testing.T) {
	data = MessageEncode(0,100, []byte("hello go"))
	if len(data) != 10 {
		t.Errorf("encode data length is incorrect, expect 10, result %v", len(data))
	}
}

func TestMessageDecode(t *testing.T) {
	msgId, routeId, body = MessageDecode(data)
	if strings.Compare(string(body), "hello go") != 0 {
		t.Errorf("expect hello go, result %v", string(body))
	}

	if msgId != 0 {
		t.Errorf("expect msgId: 0, result %v", msgId)
	}

	if routeId != 100 {
		t.Errorf("expect routeId 100, result %v", routeId)
	}
}


func TestGetHeadLength(t *testing.T) {
	length := GetHeadLength()
	if length != 4 {
		t.Errorf("expect head length 4, result %v", length)
	}
}

func TestGetNumberBytes(t *testing.T) {
	byte127 := GetNumberBytes(127)
	byte225 := GetNumberBytes(225)

	if byte127 != 1 {
		t.Errorf("byte length of 127 expect 1, result %v", byte127)
	}

	if byte225 != 2 {
		t.Errorf("byte length of 225 expect 2, result %v", byte225)
	}
}

func TestEncode(t *testing.T) {
	encode = Encode(TYPE_HANDSHAKE_ACK, data)
	t.Log(encode)
}

func TestGetBodyLength(t *testing.T) {
	t.Log("body length: ", GetBodyLength(encode))
}


func TestDecode(t *testing.T) {
	packageType, body = Decode(encode)
	if packageType != TYPE_HANDSHAKE_ACK {
		t.Errorf("expect " + string(TYPE_HANDSHAKE_ACK) + ", result %v", packageType)
	}
	if bytes.Equal(data, body) == false {
		t.Errorf("expect %v, result %v", data, body)
	}
}

func BenchmarkMessageEncode(b *testing.B) {
	for n :=0; n < b.N; n++ {
		MessageEncode(10, 20, []byte("hello go"))
	}
}

func BenchmarkMessageDecode(b *testing.B) {
	data := MessageEncode(10, 20, []byte("hello go"))
	for n :=0; n < b.N; n++ {
		MessageDecode(data)
	}
}

func BenchmarkEncode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Encode(TYPE_DATA_REQUEST, MessageEncode(100, 30, []byte("hello data")))
	}
}

func BenchmarkDecode(b *testing.B) {
	data = Encode(TYPE_DATA_REQUEST, MessageEncode(100, 30, []byte("hello data")))
	for n := 0; n < b.N; n++ {
		Decode(data)
	}
}






