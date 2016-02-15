/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */

package deb

import (
	"fmt"
	"io"
	"strings"

	"archive/tar"

	"compress/bzip2"
	"compress/gzip"
	"xi2.org/x/xz"
)

type compressionReader func(io.Reader) (io.Reader, error)

func gzipNewReader(r io.Reader) (io.Reader, error) {
	return gzip.NewReader(r)
}

func xzNewReader(r io.Reader) (io.Reader, error) {
	return xz.NewReader(r, 0)
}

func bzipNewReader(r io.Reader) (io.Reader, error) {
	return bzip2.NewReader(r), nil
}

var knownCompressionAlgorithms = map[string]compressionReader{
	".tar.gz":  gzipNewReader,
	".tar.bz2": bzipNewReader,
	".tar.xz":  xzNewReader,
}

func (e *DebEntry) IsTarfile() bool {
	return e.getCompressionReader() != nil
}

func (e *DebEntry) getCompressionReader() *compressionReader {
	for key, decompressor := range knownCompressionAlgorithms {
		if strings.HasSuffix(e.Name, key) {
			return &decompressor
		}
	}
	return nil
}

func (e *DebEntry) Tarfile() (*tar.Reader, error) {
	decompressor := e.getCompressionReader()
	if decompressor == nil {
		return nil, fmt.Errorf("%s appears to not be a tarfile", e.Name)
	}
	reader, err := (*decompressor)(e.Data)
	if err != nil {
		return nil, err
	}
	return tar.NewReader(reader), nil
}

// vim: foldmethod=marker
