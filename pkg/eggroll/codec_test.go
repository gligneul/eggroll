// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type mock struct {
	Value string `json:"value"`
}

func TestNewJSONCodec(t *testing.T) {
	codec := NewJSONCodec[mock]()
	if codec.Type().Elem().Name() != "mock" {
		t.Fatalf("wrong type")
	}
	if codec.Key().String() != "fa88ac58" {
		t.Fatalf("wrong key: %v", codec.Key().String())
	}
}

func TestJSONCodecDecode(t *testing.T) {
	codec := NewJSONCodec[mock]()

	value, err := codec.Decode([]byte(""))
	if value != nil || err == nil {
		t.Fatalf("expected nil, err; got %v, %v", value, err)
	}
	if err.Error() != "unmarshal error: unexpected end of JSON input" {
		t.Fatalf("wrong err: %v", err)
	}

	value, err = codec.Decode([]byte(`{"value": "eggroll"}`))
	if value == nil || err != nil {
		t.Fatalf("expected value, nil; got %v, %v", value, err)
	}
	mockValue := value.(*mock)
	if mockValue.Value != "eggroll" {
		t.Fatalf("wrong decoded value")
	}
}

func TestJSONCodecEncode(t *testing.T) {
	codec := NewJSONCodec[mock]()

	var buffer bytes.Buffer
	err := codec.Encode(&buffer, mock{Value: "eggroll"})
	if err == nil {
		t.Fatalf("expected err; got nil")
	}
	if err.Error() != "wrong encode type: eggroll.mock" {
		t.Fatalf("wrong err: %v", err)
	}

	err = codec.Encode(&buffer, &mock{Value: "eggroll"})
	if err != nil {
		t.Fatalf("expected nil; got %v", err)
	}
	if buffer.String() != "{\"value\":\"eggroll\"}\n" {
		t.Fatalf("wrong value: %v", buffer.String())
	}
}

func TestCodecManagerPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("the code did not panic")
		}
		msg := r.(string)
		if msg != "two codecs with same key: fa88ac58" {
			t.Fatalf("wrong panic: %v", msg)
		}
	}()
	codecs := []Codec{
		NewJSONCodec[mock](),
		NewJSONCodec[mock](),
	}
	_ = newCodecManager(codecs)
}

func TestCodecManagerDecode(t *testing.T) {
	codecs := []Codec{NewJSONCodec[mock]()}
	manager := newCodecManager(codecs)

	err := manager.decode(nil).(error)
	if err.Error() != "payload too small" {
		t.Fatalf("wrong error: %v", err)
	}

	payload := common.Hex2Bytes("deadbeef")
	err = manager.decode(payload).(error)
	if err.Error() != "codec not found for deadbeef" {
		t.Fatalf("wrong error: %v", err)
	}

	payload = append(common.Hex2Bytes("fa88ac58"), []byte(`{"value": "eggroll"}`)...)
	value := manager.decode(payload).(*mock)
	if value.Value != "eggroll" {
		t.Fatalf("wrong value: %v", value.Value)
	}
}

func TestCodecManagerEncode(t *testing.T) {
	codecs := []Codec{NewJSONCodec[mock]()}
	manager := newCodecManager(codecs)

	payload, err := manager.encode(struct{}{})
	if payload != nil || err == nil {
		t.Fatalf("expected nil, err; got %v, %v", payload, err)
	}
	if err.Error() != "codec not found for struct {}" {
		t.Fatalf("wrong error: %v", err)
	}

	payload, err = manager.encode(&mock{Value: "eggroll"})
	if payload == nil || err != nil {
		t.Fatalf("expected payload, nil; got %v, %v", payload, err)
	}
	expected := append(common.Hex2Bytes("fa88ac58"), []byte(`{"value": "eggroll"}`)...)
	if reflect.DeepEqual(payload, expected) {
		t.Fatalf("wrong payload: %v", string(payload))
	}
}
