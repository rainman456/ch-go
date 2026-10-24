package compress

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/go-faster/city"
	"github.com/stretchr/testify/require"
)

// encodeNoneBlock builds a valid None-encoded block whose declared data size
// may differ from the payload length, which the Writer never produces.
func encodeNoneBlock(payload []byte, dataSize int) []byte {
	raw := make([]byte, headerSize+len(payload))
	raw[hMethod] = byte(encodedNone)
	binary.LittleEndian.PutUint32(raw[hRawSize:], uint32(len(payload)+compressHeaderSize))
	binary.LittleEndian.PutUint32(raw[hDataSize:], uint32(dataSize))
	copy(raw[headerSize:], payload)
	h := city.CH128(raw[hMethod:])
	binary.LittleEndian.PutUint64(raw[0:8], h.Low)
	binary.LittleEndian.PutUint64(raw[8:16], h.High)
	return raw
}

func TestGrowBuf(t *testing.T) {
	t.Run("Reuse", func(t *testing.T) {
		b := bytes.Repeat([]byte{0xAA}, 32)
		got := growBuf(b, 16)
		require.Len(t, got, 16)
		require.Equal(t, bytes.Repeat([]byte{0xAA}, 16), got)
	})
	t.Run("Grow", func(t *testing.T) {
		got := growBuf(bytes.Repeat([]byte{0xAA}, 8), 32)
		require.Len(t, got, 32)
		require.Equal(t, make([]byte, 32), got)
	})
}

// TestReaderNoneShortPayload checks that a None block declaring more data than
// it carries does not expose bytes from the previous block through the tail.
// The first block is deliberately larger so the second one reuses its buffer.
func TestReaderNoneShortPayload(t *testing.T) {
	first := encodeNoneBlock(bytes.Repeat([]byte{0xAA}, 32), 32)
	second := encodeNoneBlock(bytes.Repeat([]byte{0xBB}, 4), 16)
	stream := append(first, second...)

	r := NewReader(bytes.NewReader(stream))
	out := make([]byte, 32+16)
	_, err := io.ReadFull(r, out)
	require.NoError(t, err)

	require.Equal(t, bytes.Repeat([]byte{0xAA}, 32), out[:32])
	require.Equal(t, bytes.Repeat([]byte{0xBB}, 4), out[32:36])
	require.Equal(t, make([]byte, 12), out[36:])
}
