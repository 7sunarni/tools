package main

import (
	"bytes"
	"io"

	"go-playground/mypackage"
)

type TestResponseWriterPtr struct {
}

func (m *TestResponseWriterPtr) Write(_ []byte) (int, error) {
	return 0, nil
}

type TestResponseWriter struct {
}

func (m TestResponseWriter) Write(_ []byte) (int, error) {
	return 0, nil
}

func (m TestResponseWriter) Read(_ []byte) (int, error) {
	return 0, nil
}

func main() {
	b := bytes.NewBuffer([]byte{})
	wPtr := &TestResponseWriterPtr{}
	io.Copy(wPtr, b)
	wPtr.Write([]byte("hello, world"))
	io.Copy(&TestResponseWriterPtr{}, b)

	w1 := &TestResponseWriter{}
	w2 := TestResponseWriter{}
	io.Copy(w1, b)
	io.Copy(w2, b)
	w1.Write([]byte("hello, world"))
	io.Copy(TestResponseWriter{}, b)
	io.Copy(&TestResponseWriter{}, b)
	io.ReadAll(w1)
	io.ReadAll(w2)

	// TODO: local package
	Mcopy(wPtr, &TestResponseWriter{})
	mypackage.AAAAMMcopy(wPtr, &TestResponseWriter{})
}

type OtherResponesWriter struct{}

func (m *OtherResponesWriter) Write(_ []byte) (int, error) {
	return 0, nil
}

func Mcopy(a io.Writer, b io.Writer) {

}
