package core

import (
	"crypto/sha256"
	"hash"
	"io"
)

type Storage interface {
	Add()
}

type StorageWriter interface {
	hash.Hash
}

func NewStorage(w io.Writer) StorageWriter {
	return newStorageWriter(w)
}

type storageWriter struct {
	writer io.Writer
	w io.Writer
	h hash.Hash
}

func newStorageWriter(w io.Writer) *storageWriter {
	h := sha256.New()
	return &storageWriter{
		w: w,
		h: h,
		writer: io.MultiWriter(w, h),
	}
}

func (s *storageWriter) Write(p []byte) (int, error) {
	return s.writer.Write(p)
}

func (w *storageWriter) Sum(b []byte) []byte {
	return w.h.Sum(b)
}

func (w *storageWriter) Reset() {
	w.h.Reset()
}

func (w *storageWriter) Size() int {
	return w.h.Size()
}

func (w *storageWriter) BlockSize() int {
	return w.h.BlockSize()
}