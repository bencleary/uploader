package storage

import "io"

type chainedReadCloser struct {
	reader  io.Reader
	closers []io.Closer
}

func newChainedReadCloser(reader io.Reader, closers ...io.Closer) io.ReadCloser {
	return &chainedReadCloser{
		reader:  reader,
		closers: closers,
	}
}

func (c *chainedReadCloser) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

func (c *chainedReadCloser) Close() error {
	var firstErr error
	for _, closer := range c.closers {
		if closer == nil {
			continue
		}
		if err := closer.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
