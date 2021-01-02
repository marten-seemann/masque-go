package masque

import (
	"bytes"
	"fmt"
	"io"

	"github.com/lucas-clemente/quic-go/quicvarint"
)

const maxStreamChunkSize = 1500

type streamChunk struct {
	Type uint64
	Data []byte
}

func parseStreamChunk(r *bytes.Reader) (*streamChunk, error) {
	t, err := quicvarint.Read(r)
	if err != nil {
		return nil, err
	}
	l, err := quicvarint.Read(r)
	if err != nil {
		return nil, err
	}
	if l > maxStreamChunkSize {
		return nil, fmt.Errorf("too large stream chunk: %d bytes (allowed: %d bytes)", l, maxStreamChunkSize)
	}
	b := make([]byte, l)
	if _, err := io.ReadFull(r, b); err != nil {
		if err == io.ErrUnexpectedEOF {
			err = io.EOF
		}
		return nil, err
	}
	return &streamChunk{
		Type: t,
		Data: b,
	}, nil
}

func (c *streamChunk) Write(b *bytes.Buffer) error {
	quicvarint.Write(b, c.Type)
	quicvarint.Write(b, uint64(len(c.Data)))
	_, err := b.Write(c.Data)
	return err
}
