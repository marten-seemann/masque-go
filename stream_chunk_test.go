package masque

import (
	"bytes"
	"io"

	"github.com/lucas-clemente/quic-go/quicvarint"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreamChunk", func() {
	It("parses", func() {
		b := &bytes.Buffer{}
		quicvarint.Write(b, 1337)
		quicvarint.Write(b, 6)
		b.Write([]byte("foobar"))

		c, err := parseStreamChunk(bytes.NewReader(b.Bytes()))
		Expect(err).ToNot(HaveOccurred())
		Expect(c.Type).To(BeEquivalentTo(1337))
		Expect(c.Data).To(Equal([]byte("foobar")))
	})

	It("rejects chunks that are too large", func() {
		b := &bytes.Buffer{}
		quicvarint.Write(b, 1337)
		quicvarint.Write(b, maxStreamChunkSize+1)
		b.Write(make([]byte, maxStreamChunkSize+1))

		_, err := parseStreamChunk(bytes.NewReader(b.Bytes()))
		Expect(err).To(MatchError("too large stream chunk: 1501 bytes (allowed: 1500 bytes)"))
	})

	It("errors on EOFs", func() {
		b := &bytes.Buffer{}
		quicvarint.Write(b, 1337)
		quicvarint.Write(b, 6)
		b.Write([]byte("foobar"))

		data := b.Bytes()
		for i := range data {
			_, err := parseStreamChunk(bytes.NewReader(data[0:i]))
			Expect(err).To(MatchError(io.EOF))
		}
	})

	It("writes", func() {
		c := &streamChunk{
			Type: 42,
			Data: []byte("foobar"),
		}
		b := &bytes.Buffer{}
		Expect(c.Write(b)).To(Succeed())

		parsed, err := parseStreamChunk(bytes.NewReader(b.Bytes()))
		Expect(err).ToNot(HaveOccurred())
		Expect(parsed).To(Equal(c))
	})
})
