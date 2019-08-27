// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

type stringWriter interface {
	io.Writer
	io.StringWriter
}

func rpcRequest2XML(method string, rpc ...interface{}) (string, error) {
	buffer := new(strings.Builder)
	buffer.WriteString("<methodCall><methodName>")

	buffer.WriteString(method)
	buffer.WriteString("</methodName>")
	err := rpcParams2XML(buffer, rpc...)
	buffer.WriteString("</methodCall>")
	return buffer.String(), err
}

func rpcResponse2XML(rpc ...interface{}) (string, error) {
	buffer := new(strings.Builder)
	buffer.WriteString("<methodResponse>")
	err := rpcParams2XML(buffer, rpc...)
	buffer.WriteString("</methodResponse>")
	return buffer.String(), err
}

func rpcParams2XML(buffer stringWriter, rpc ...interface{}) error {
	var err error
	buffer.WriteString("<params>")

	var elem reflect.Value
	for _, r := range rpc {
		val := reflect.ValueOf(r)
		switch val.Kind() {
		case reflect.Interface, reflect.Ptr:
			elem = val.Elem()
		default:
			elem = val
		}
		if elem.Kind() != reflect.Struct {
			buffer.WriteString("<param>")
			err = rpc2XML(buffer, elem.Interface())
			buffer.WriteString("</param>")
			continue
		}

		numField := elem.NumField()

		for i := 0; i < numField; i++ {
			buffer.WriteString("<param>")
			err = rpc2XML(buffer, elem.Field(i).Interface())
			buffer.WriteString("</param>")
		}
	}

	buffer.WriteString("</params>")
	return err
}

func rpc2XML(w stringWriter, value interface{}) error {
	w.WriteString("<value>")
	switch reflect.ValueOf(value).Kind() {
	case reflect.Int:
		fmt.Fprintf(w, "<int>%d</int>", value.(int))
	case reflect.Float64:
		fmt.Fprintf(w, "<double>%f</double>", value.(float64))
	case reflect.String:
		w.Write([]byte(string2XML(value.(string))))
	case reflect.Bool:
		w.Write([]byte(bool2XML(value.(bool))))
	case reflect.Struct:
		if reflect.TypeOf(value).String() != "time.Time" {
			struct2XML(w, value)
		} else {
			w.Write([]byte(time2XML(value.(time.Time))))
		}
	case reflect.Slice, reflect.Array:
		// FIXME: is it the best way to recognize '[]byte'?
		if reflect.TypeOf(value).String() != "[]uint8" {
			array2XML(w, value)
		} else {
			w.WriteString(base642XML(value.([]byte)))
		}
	case reflect.Ptr:
		if reflect.ValueOf(value).IsNil() {
			w.WriteString("<nil/>")
		}
	}
	w.WriteString("</value>")
	return nil
}

func bool2XML(value bool) string {
	var b string
	if value {
		b = "1"
	} else {
		b = "0"
	}
	return fmt.Sprintf("<boolean>%s</boolean>", b)
}

func string2XML(value string) string {
	value = strings.Replace(value, "&", "&amp;", -1)
	value = strings.Replace(value, "\"", "&quot;", -1)
	value = strings.Replace(value, "<", "&lt;", -1)
	value = strings.Replace(value, ">", "&gt;", -1)
	return fmt.Sprintf("<string>%s</string>", value)
}

func struct2XML(w stringWriter, value interface{}) {
	w.WriteString("<struct>")
	for i := 0; i < reflect.TypeOf(value).NumField(); i++ {
		field := reflect.ValueOf(value).Field(i)
		fieldType := reflect.TypeOf(value).Field(i)
		var name string
		if fieldType.Tag.Get("xml") != "" {
			name = fieldType.Tag.Get("xml")
		} else {
			name = fieldType.Name
		}
		w.WriteString("<member>")
		fmt.Fprintf(w, "<name>%s</name>", name)
		rpc2XML(w, field.Interface())
		w.WriteString("</member>")
	}
	w.WriteString("</struct>")
	return
}

func array2XML(w stringWriter, value interface{}) {
	w.WriteString("<array><data>")
	for i := 0; i < reflect.ValueOf(value).Len(); i++ {
		rpc2XML(w, reflect.ValueOf(value).Index(i).Interface())
	}
	w.WriteString("</data></array>")
}

func time2XML(t time.Time) string {
	/*
		// TODO: find out whether we need to deal
		// here with TZ
		var tz string;
		zone, offset := t.Zone()
		if zone == "UTC" {
			tz = "Z"
		} else {
			tz = fmt.Sprintf("%03d00", offset / 3600 )
		}
	*/
	return fmt.Sprintf("<dateTime.iso8601>%04d%02d%02dT%02d:%02d:%02d</dateTime.iso8601>",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func base642XML(data []byte) string {
	str := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("<base64>%s</base64>", str)
}
