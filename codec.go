// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggroll

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Key that identifies the codec type.
type CodecKey [4]byte

// String representation of a codec key.
func (k CodecKey) String() string {
	return common.Bytes2Hex(k[:])
}

// Decode raw payload into Go value and encode Go value into raw payload.
type Codec interface {

	// Get the codec key used to identify which coded the payload uses.
	Key() CodecKey

	// Get the Go type of the codec.
	Type() reflect.Type

	// Try to decode the given payload into Go value.
	// The type of the value should be the same one returned by Type().
	Decode(payload []byte) (any, error)

	// Encode a given Go value to payload.
	// The type of the value should be the same one returned by Type().
	Encode(w io.Writer, value any) error
}

// JSON codec for values of type T.
type JSONCodec struct {
	type_ reflect.Type
}

// Create a new JSON codec for the struct type.
func NewJSONCodec[T any]() *JSONCodec {
	// Check if T is struct
	type_ := reflect.TypeOf((*T)(nil))
	if type_.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("type must be a struct; is %v\n", type_))
	}
	return &JSONCodec{
		type_: type_,
	}
}

// The JSON codec uses the first 4 bytes of the keccak of the type name as the codec key.
func (c *JSONCodec) Key() CodecKey {
	hash := crypto.Keccak256Hash([]byte(c.type_.Elem().Name()))
	return CodecKey(hash[:4])
}

// Return the Go type (pointer to struct).
func (c *JSONCodec) Type() reflect.Type {
	return c.type_
}

// Try to decode the given payload into Go value.
// Return a pointer to the struct.
func (c *JSONCodec) Decode(payload []byte) (any, error) {
	value := reflect.New(c.type_.Elem()).Interface()
	if err := json.Unmarshal(payload, value); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	return value, nil
}

// Encode a given Go value to payload.
// Receives a pointer of the struct.
func (c *JSONCodec) Encode(w io.Writer, value any) error {
	// sanity check
	if reflect.TypeOf(value) != c.type_ {
		return fmt.Errorf("wrong encode type: %v", reflect.TypeOf(value))
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("marshal error: %v", err)
	}
	return nil
}

// Manages codecs.
type codecManager struct {
	decoderMap map[CodecKey]Codec
	encoderMap map[reflect.Type]Codec
}

// Create a new codec manager.
func newCodecManager(codecs []Codec) *codecManager {
	m := &codecManager{
		decoderMap: make(map[CodecKey]Codec),
		encoderMap: make(map[reflect.Type]Codec),
	}
	for _, codec := range codecs {
		_, ok := m.decoderMap[codec.Key()]
		if ok {
			panic(fmt.Sprintf("two codecs with same key: %v", codec.Key()))
		}
		_, ok = m.encoderMap[codec.Type()]
		if ok {
			panic(fmt.Sprintf("two codecs with same type: %v", codec.Type()))
		}
		m.decoderMap[codec.Key()] = codec
		m.encoderMap[codec.Type()] = codec
	}
	return m
}

// Try to decode a value.
// If fails, return the error in the place of the value.
func (m *codecManager) decode(payload []byte) any {
	if len(payload) < 4 {
		return fmt.Errorf("payload too small")
	}
	key := CodecKey(payload[:4])
	payload = payload[4:]
	codec, ok := m.decoderMap[key]
	if !ok {
		return fmt.Errorf("codec not found for %v", key)
	}
	input, err := codec.Decode(payload)
	if err != nil {
		return err
	}
	return input
}

// Try to encode a value.
func (m *codecManager) encode(value any) ([]byte, error) {
	type_ := reflect.TypeOf(value)
	codec, ok := m.encoderMap[type_]
	if !ok {
		return nil, fmt.Errorf("codec not found for %v", type_)
	}
	var buffer bytes.Buffer
	key := codec.Key()
	_, err := buffer.Write(key[:])
	if err != nil {
		return nil, fmt.Errorf("error when writting to buffer: %v", err)
	}
	err = codec.Encode(&buffer, value)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
