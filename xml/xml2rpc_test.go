// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"reflect"
	"testing"
	"time"
)

type SubStructXML2RPC struct {
	Foo  int
	Bar  string
	Data []int
}

type StructXML2RPC struct {
	Int    int
	Float  float64
	Str    string
	Bool   bool
	Sub    SubStructXML2RPC
	Time   time.Time
	Base64 []byte
}

func TestXML2RPC(t *testing.T) {
	req := new(StructXML2RPC)
	err := xml2RPC("<methodCall><methodName>Some.Method</methodName><params><param><value><i4>123</i4></value></param><param><value><double>3.145926</double></value></param><param><value><string>Hello, World!</string></value></param><param><value><boolean>0</boolean></value></param><param><value><struct><member><name>Foo</name><value><int>42</int></value></member><member><name>Bar</name><value><string>I'm Bar</string></value></member><member><name>Data</name><value><array><data><value><int>1</int></value><value><int>2</int></value><value><int>3</int></value></data></array></value></member></struct></value></param><param><value><dateTime.iso8601>20120717T14:08:55</dateTime.iso8601></value></param><param><value><base64>eW91IGNhbid0IHJlYWQgdGhpcyE=</base64></value></param></params></methodCall>", req)
	if err != nil {
		t.Error("XML2RPC conversion failed", err)
	}
	expectedReq := &StructXML2RPC{123, 3.145926, "Hello, World!", false, SubStructXML2RPC{42, "I'm Bar", []int{1, 2, 3}}, time.Date(2012, time.July, 17, 14, 8, 55, 0, time.Local), []byte("you can't read this!")}
	if !reflect.DeepEqual(req, expectedReq) {
		t.Error("XML2RPC conversion failed")
		t.Error("Expected", expectedReq)
		t.Error("Got", req)
	}
}

type StructSpecialCharsXML2RPC struct {
	String1 string
}

func TestXML2RPCSpecialChars(t *testing.T) {
	req := new(StructSpecialCharsXML2RPC)
	err := xml2RPC("<methodResponse><params><param><value><string> &amp; &quot; &lt; &gt; </string></value></param></params></methodResponse>", req)
	if err != nil {
		t.Error("XML2RPC conversion failed", err)
	}
	expectedReq := &StructSpecialCharsXML2RPC{" & \" < > "}
	if !reflect.DeepEqual(req, expectedReq) {
		t.Error("XML2RPC conversion failed")
		t.Error("Expected", expectedReq)
		t.Error("Got", req)
	}
}

type StructNilXML2RPC struct {
	Ptr *int
}

func TestXML2RPCNil(t *testing.T) {
	req := new(StructNilXML2RPC)
	err := xml2RPC("<methodResponse><params><param><value><nil/></value></param></params></methodResponse>", req)
	if err != nil {
		t.Error("XML2RPC conversion failed", err)
	}
	expectedReq := &StructNilXML2RPC{nil}
	if !reflect.DeepEqual(req, expectedReq) {
		t.Error("XML2RPC conversion failed")
		t.Error("Expected", expectedReq)
		t.Error("Got", req)
	}
}

type StructXML2RPCSubArgs struct {
	String1 string
	String2 string
	ID      int
}

type StructXML2RPCHelloArgs struct {
	Args StructXML2RPCSubArgs
}

func TestXML2RPCLowercasedMethods(t *testing.T) {
	req := new(StructXML2RPCHelloArgs)
	err := xml2RPC(`
	<methodCall>
	<params>
		<param>
			<value>
				<struct>
					<member>
						<name>string1</name>
						<value>
							<string>I'm a first string</string>
						</value>
					</member>
					<member>
						<name>string2</name>
						<value>
							<string>I'm a second string</string>
						</value>
					</member>
					<member>
						<name>id</name>
						<value>
							<int>1</int>
						</value>
					</member>
				</struct>
			</value>
		</param>
	</params>
</methodCall>`, req)
	if err != nil {
		t.Error("XML2RPC conversion failed", err)
	}
	args := StructXML2RPCSubArgs{"I'm a first string", "I'm a second string", 1}
	expectedReq := &StructXML2RPCHelloArgs{args}
	if !reflect.DeepEqual(req, expectedReq) {
		t.Error("XML2RPC conversion failed")
		t.Error("Expected", expectedReq)
		t.Error("Got", req)
	}
}

func TestXML2PRCFaultCall(t *testing.T) {
	req := new(StructXML2RPCHelloArgs)
	data := `<?xmlversion="1.0"?><methodResponse><fault><value><struct><member><name>faultCode</name><value><int>116</int></value></member><member><name>faultString</name><value><string>Error
Requiredattribute'user'notfound:
[{'User',"gggg"},{'Host',"sss.com"},{'Password',"ssddfsdf"}]
</string></value></member></struct></value></fault></methodResponse>`

	errstr := `Error
Requiredattribute'user'notfound:
[{'User',"gggg"},{'Host',"sss.com"},{'Password',"ssddfsdf"}]
`

	err := xml2RPC(data, req)

	fault, ok := err.(Fault)
	if !ok {
		t.Errorf("error should be of concrete type Fault, but got %v", err)
	} else {
		if fault.Code != 116 {
			t.Errorf("expected fault.Code to be %d, but got %d", 116, fault.Code)
		}
		if fault.String != errstr {
			t.Errorf("fault.String should be:\n\n%s\n\nbut got:\n\n%s\n", errstr, fault.String)
		}
	}
}

// ProcessInfo Get info about a process named name
type ProcessInfo struct {
	Name          string `json:"name" xml:"name"`
	Group         string `json:"group" xml:"group"`
	Description   string `json:"description" xml:"description"`
	Start         int    `json:"start" xml:"start"`
	Stop          int    `json:"stop" xml:"stop"`
	Now           int    `json:"now" xml:"now"`
	State         int    `json:"state" xml:"state"`
	Statename     string `json:"statename" xml:"statename"`
	Spawnerr      string `json:"spawnerr" xml:"spawnerr"`
	Exitstatus    int    `json:"exitstatus" xml:"exitstatus"`
	Logfile       string `json:"logfile" xml:"logfile"`
	StdoutLogfile string `json:"stdout_logfile" xml:"stdout_logfile"`
	StderrLogfile string `json:"stderr_logfile" xml:"stderr_logfile"`
	Pid           int    `json:"pid" xml:"pid"`
}

// ProcessInfoResult ProcessInfoResult
type ProcessInfoResult struct {
	Process ProcessInfo
}

func TestXML2PRCUnderscore(t *testing.T) {
	req := new(ProcessInfoResult)
	data := `<?xml version='1.0'?>
	<methodResponse>
		<params>
			<param>
				<value>
					<struct>
						<member>
							<name>description</name>
							<value>
								<string>pid 6169, uptime 26 days, 16:14:19</string>
							</value>
						</member>
						<member>
							<name>pid</name>
							<value>
								<int>6169</int>
							</value>
						</member>
						<member>
							<name>stderr_logfile</name>
							<value>
								<string>/data/logs/supervisor/api-gateway_stderr.log</string>
							</value>
						</member>
						<member>
							<name>stop</name>
							<value>
								<int>0</int>
							</value>
						</member>
						<member>
							<name>logfile</name>
							<value>
								<string>/data/logs/supervisor/api-gateway_stdout.log</string>
							</value>
						</member>
						<member>
							<name>exitstatus</name>
							<value>
								<int>0</int>
							</value>
						</member>
						<member>
							<name>spawnerr</name>
							<value>
								<string></string>
							</value>
						</member>
						<member>
							<name>now</name>
							<value>
								<int>1566880590</int>
							</value>
						</member>
						<member>
							<name>group</name>
							<value>
								<string>api-gateway</string>
							</value>
						</member>
						<member>
							<name>name</name>
							<value>
								<string>api-gateway</string>
							</value>
						</member>
						<member>
							<name>statename</name>
							<value>
								<string>RUNNING</string>
							</value>
						</member>
						<member>
							<name>start</name>
							<value>
								<int>1564575731</int>
							</value>
						</member>
						<member>
							<name>state</name>
							<value>
								<int>20</int>
							</value>
						</member>
						<member>
							<name>stdout_logfile</name>
							<value>
								<string>/data/logs/supervisor/api-gateway_stdout.log</string>
							</value>
						</member>
					</struct>
				</value>
			</param>
		</params>
	</methodResponse>`

	err := xml2RPC(data, req)

	if err != nil {
		t.Error(err)
	}

	s := "/data/logs/supervisor/api-gateway_stdout.log"
	if req.Process.StdoutLogfile != s {
		t.Error("req.Process.StdoutLogfile should be \"", s, "\", not ", req.Process.StdoutLogfile)
	}
	s = "/data/logs/supervisor/api-gateway_stderr.log"
	if req.Process.StderrLogfile != s {
		t.Error("req.Process.StderrLogfile should be \"", s, "\", not ", req.Process.StderrLogfile)
	}

	if req.Process.Spawnerr != "" {
		t.Error("req.Process.Spawnerr should be \"\", not ", req.Process.Spawnerr)
	}
}

func TestXML2PRCISO88591(t *testing.T) {
	req := new(StructXML2RPCHelloArgs)
	data := `<?xml version="1.0" encoding="ISO-8859-1"?><methodResponse><fault><value><struct><member><name>faultCode</name><value><int>116</int></value></member><member><name>faultString</name><value><string>Error
Requiredattribute'user'notfound:
[{'User',"` + "\xd6\xf1\xe4" + `"},{'Host',"sss.com"},{'Password',"ssddfsdf"}]
</string></value></member></struct></value></fault></methodResponse>`

	errstr := `Error
Requiredattribute'user'notfound:
[{'User',"Öñä"},{'Host',"sss.com"},{'Password',"ssddfsdf"}]
`

	err := xml2RPC(data, req)

	fault, ok := err.(Fault)
	if !ok {
		t.Errorf("error should be of concrete type Fault, but got %v", err)
	} else {
		if fault.Code != 116 {
			t.Errorf("expected fault.Code to be %d, but got %d", 116, fault.Code)
		}
		if fault.String != errstr {
			t.Errorf("fault.String should be:\n\n%s\n\nbut got:\n\n%s\n", errstr, fault.String)
		}
	}
}
