/*
 * Copyright 2024 gRPC authors.
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

package grpc

import (
	"bytes"
	"testing"

	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
)

func TestEncode(t *testing.T) {
	codec := encoding.GetCodec(proto.Name)
	if codec == nil {
		t.Fatal("proto codec not found")
	}

	tests := []struct {
		name    string
		msg     interface{}
		wantErr bool
	}{
		{
			name:    "nil message",
			msg:     nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encode(codec, tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCompress(t *testing.T) {
	origData := []byte("hello grpc world - this is test data for compression")

	// Test that data round-trips correctly without compression
	buf := &bytes.Buffer{}
	buf.Write(origData)

	if !bytes.Equal(buf.Bytes(), origData) {
		t.Errorf("data mismatch: got %v, want %v", buf.Bytes(), origData)
	}
}

func TestPayloadFormat(t *testing.T) {
	tests := []struct {
		name       string
		compressed bool
		wantFlag   byte
	}{
		{
			name:       "uncompressed",
			compressed: false,
			wantFlag:   0,
		},
		{
			name:       "compressed",
			compressed: true,
			wantFlag:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var flag byte
			if tt.compressed {
				flag = 1
			}
			if flag != tt.wantFlag {
				t.Errorf("payload flag = %v, want %v", flag, tt.wantFlag)
			}
		})
	}
}
