/*
 * Copyright 2014 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package grpc implements an RPC system called gRPC.
package grpc

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The format of the payload: compressed or not?
type payloadFormat uint8

const (
	compressionNone payloadFormat = 0 // no compression
	compressionMade payloadFormat = 1 // compressed
)

// maxRecvMsgSize is the default maximum message size the client/server can receive.
const defaultMaxRecvMsgSize = 1024 * 1024 * 4 // 4 MiB

// maxSendMsgSize is the default maximum message size the client can send.
const defaultMaxSendMsgSize = math.MaxInt32

// parser reads complete gRPC messages from the underlying reader.
type parser struct {
	// r is the underlying reader.
	r io.Reader

	// The header of a gRPC message. See:
	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
	header [5]byte
}

// recvMsg reads a complete gRPC message from the stream.
//
// It returns the message and its payload format. The caller owns the returned
// msg memory.
func (p *parser) recvMsg(maxRecvMsgSize int) (pf payloadFormat, msg []byte, err error) {
	if _, err := p.r.Read(p.header[:]); err != nil {
		return 0, nil, err
	}

	pf = payloadFormat(p.header[0])
	length := binary.BigEndian.Uint32(p.header[1:])

	if length == 0 {
		return pf, nil, nil
	}
	if int64(length) > int64(maxInt) {
		return 0, nil, status.Errorf(codes.ResourceExhausted, "grpc: received message larger than max length allowed on current machine (%d vs. %d)", length, maxInt)
	}
	if int(length) > maxRecvMsgSize {
		return 0, nil, status.Errorf(codes.ResourceExhausted, "grpc: received message larger than max (%d vs. %d)", length, maxRecvMsgSize)
	}

	msg = make([]byte, int(length))
	if _, err := io.ReadFull(p.r, msg); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return 0, nil, err
	}
	return pf, msg, nil
}

// encode serializes msg, returning the payload format and serialized bytes.
func encode(c Codec, msg interface{}) ([]byte, error) {
	if msg == nil {
		return nil, nil
	}
	b, err := c.Marshal(msg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "grpc: error while marshaling: %v", err.Error())
	}
	return b, nil
}

// compress compresses the given data using gzip.
var gzPool = sync.Pool{
	New: func() interface{} {
		w, err := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
		if err != nil {
			panic(err)
		}
		return w
	},
}

func compress(in []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzPool.Get().(*gzip.Writer)
	defer gzPool.Put(w)
	w.Reset(&buf)
	if _, err := w.Write(in); err != nil {
		return nil, fmt.Errorf("grpc: failed to compress: %v", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("grpc: failed to close gzip writer: %v", err)
	}
	return buf.Bytes(), nil
}

// decompress decompresses the given gzip-compressed data.
func decompress(in []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, fmt.Errorf("grpc: failed to create gzip reader: %v", err)
	}
	defer r.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("grpc: failed to decompress: %v", err)
	}
	return out, nil
}

// maxInt is the maximum int value on the current platform.
const maxInt = int(^uint(0) >> 1)
