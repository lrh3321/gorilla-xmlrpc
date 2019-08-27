// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"testing"
	"time"
)

type SubStructRPC2XML struct {
	Foo  int
	Bar  string
	Data []int
}

type StructRPC2XML struct {
	Int    int
	Float  float64
	Str    string
	Bool   bool
	Sub    SubStructRPC2XML
	Time   time.Time
	Base64 []byte
}

func TestRPC2XML(t *testing.T) {
	req := &StructRPC2XML{123, 3.145926, "Hello, World!", false, SubStructRPC2XML{42, "I'm Bar", []int{1, 2, 3}}, time.Date(2012, time.July, 17, 14, 8, 55, 0, time.Local), []byte("you can't read this!")}
	xml, err := rpcRequest2XML("Some.Method", req)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}
	expected := "<methodCall><methodName>Some.Method</methodName><params><param><value><int>123</int></value></param><param><value><double>3.145926</double></value></param><param><value><string>Hello, World!</string></value></param><param><value><boolean>0</boolean></value></param><param><value><struct><member><name>Foo</name><value><int>42</int></value></member><member><name>Bar</name><value><string>I'm Bar</string></value></member><member><name>Data</name><value><array><data><value><int>1</int></value><value><int>2</int></value><value><int>3</int></value></data></array></value></member></struct></value></param><param><value><dateTime.iso8601>20120717T14:08:55</dateTime.iso8601></value></param><param><value><base64>eW91IGNhbid0IHJlYWQgdGhpcyE=</base64></value></param></params></methodCall>"
	if xml != expected {
		t.Error("RPC2XML conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}
}

type StructSpecialCharsRPC2XML struct {
	String1 string
}

func TestRPC2XMLSpecialChars(t *testing.T) {
	req := &StructSpecialCharsRPC2XML{" & \" < > "}
	xml, err := rpcResponse2XML(req)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}
	expected := "<methodResponse><params><param><value><string> &amp; &quot; &lt; &gt; </string></value></param></params></methodResponse>"
	if xml != expected {
		t.Error("RPC2XML Special chars conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}
}

type StructNilRPC2XML struct {
	Ptr *int
}

func TestRPC2XMLNil(t *testing.T) {
	req := &StructNilRPC2XML{nil}
	xml, err := rpcResponse2XML(req)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}
	expected := "<methodResponse><params><param><value><nil/></value></param></params></methodResponse>"
	if xml != expected {
		t.Error("RPC2XML Special chars conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}
}
func TestRPC2XMLStruct(t *testing.T) {
	req := StructNilRPC2XML{nil}
	xml, err := rpcResponse2XML(req)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}
	expected := "<methodResponse><params><param><value><nil/></value></param></params></methodResponse>"
	if xml != expected {
		t.Error("RPC2XML Special chars conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}
}

func TestRPC2XMLMuiti(t *testing.T) {
	arg1 := "hello"
	arg2 := true
	xml, err := rpcResponse2XML(&arg1, &arg2)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}
	expected := "<methodResponse><params><param><value><string>hello</string></value></param><param><value><boolean>1</boolean></value></param></params></methodResponse>"
	if xml != expected {
		t.Error("RPC2XML Special chars conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}

}

func TestRPC2XMLMuitiStruct(t *testing.T) {
	arg1 := "hello"
	arg2 := true

	xml, err := rpcResponse2XML(arg1, arg2)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}
	expected := "<methodResponse><params><param><value><string>hello</string></value></param><param><value><boolean>1</boolean></value></param></params></methodResponse>"
	if xml != expected {
		t.Error("RPC2XML Special chars conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}
}
