//go:build !solution

package otp

import (
	"errors"
	"io"
)

type CipherReader struct {
	r    io.Reader
	prng io.Reader
}

type CipherWriter struct {
	w    io.Writer
	prng io.Reader
}

func (cipherReader CipherReader) Read(stream []byte) (nbytesCipher int, e error) {
	nbytesStream, err := cipherReader.r.Read(stream)
	var cipher = make([]byte, nbytesStream)

	nbytesCipher, _ = cipherReader.prng.Read(cipher)

	for i := 0; i < nbytesCipher; i++ {
		stream[i] = stream[i] ^ cipher[i]
	}

	if nbytesCipher < nbytesStream {
		err = errors.New("not enough []byte")
	}

	return nbytesCipher, err
}

func (cipherWriter CipherWriter) Write(stream []byte) (nbytesCipher int, e error) {
	bytes := make([]byte, len(stream))
	nbytes, _ := cipherWriter.prng.Read(bytes)
	result := make([]byte, nbytes)

	if nbytes == 0 {
		return nbytesCipher, errors.New("[]byte is empty")
	}

	for i := 0; i < nbytes; i++ {
		result[i] = stream[i] ^ bytes[i]
	}

	nbytesCipher, err := cipherWriter.w.Write(result)

	return nbytesCipher, err
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return CipherReader{r, prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return CipherWriter{w, prng}
}
