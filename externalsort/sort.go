//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"io"
	"os"
	"sort"
	"strings"
)

type Reader struct {
	buffer *bufio.Reader
}

type Writer struct {
	buffer *bufio.Writer
}

func NewReader(reader io.Reader) LineReader {
	return &Reader{
		buffer: bufio.NewReader(reader),
	}
}

func NewWriter(writer io.Writer) LineWriter {
	return &Writer{
		buffer: bufio.NewWriter(writer),
	}
}

func (reader *Reader) ReadLine() (string, error) {
	line, err := reader.buffer.ReadString('\n')

	if err != nil && err != io.EOF {
		return "", err
	}

	if err == io.EOF && line == "" {
		return "", err
	}

	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}

	return line, nil
}

func (writer *Writer) Write(s string) error {
	_, err := writer.buffer.WriteString(s + "\n")

	if err != nil {
		return err
	}

	return writer.buffer.Flush()
}

type StringAndReader struct {
	str    string
	reader LineReader
	index  int
}

type MinHeap struct {
	items []StringAndReader
}

func (minHeap *MinHeap) Len() int {
	return len(minHeap.items)
}

func (minHeap *MinHeap) Less(i, j int) bool {
	return minHeap.items[i].str < minHeap.items[j].str
}

func (minHeap *MinHeap) Swap(i, j int) {
	minHeap.items[i], minHeap.items[j] = minHeap.items[j], minHeap.items[i]
}

func (minHeap *MinHeap) Push(x interface{}) {
	minHeap.items = append(minHeap.items, x.(StringAndReader))
}

func (minHeap *MinHeap) Pop() interface{} {
	n := len(minHeap.items)
	item := minHeap.items[n-1]

	minHeap.items = minHeap.items[0 : n-1]

	return item
}

func Merge(lineWriter LineWriter, lineReaders ...LineReader) error {
	minHeap := &MinHeap{}
	heap.Init(minHeap)

	for _, lineReader := range lineReaders {
		str, err := lineReader.ReadLine()

		if err != nil {
			return err
		}

		heap.Push(minHeap, StringAndReader{str: str, reader: lineReader, index: len(lineReaders)})
	}

	for {
		if minHeap.Len() == 0 {
			break
		}

		minItem := heap.Pop(minHeap).(StringAndReader)
		err := lineWriter.Write(minItem.str)

		if err != nil {
			return err
		}

		nextLine, err := minItem.reader.ReadLine()

		if err != io.EOF {
			heap.Push(minHeap, StringAndReader{str: nextLine, reader: minItem.reader})
		}
	}

	return nil
}

func Sort(writer io.Writer, strs ...string) error {
	lineWriter := NewWriter(writer)
	files := make([]LineReader, 0)

	for _, file := range strs {
		data, err := os.ReadFile(file)

		if string(data) == "\n\n\n" {
			_, _ = writer.Write([]byte("\n"))
		}

		if err != nil {
			return err
		}

		if len(data) == 0 {
			continue
		}

		lines := strings.Split(string(data), "\n")

		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}

		sort.Strings(lines)
		err = os.WriteFile(file, []byte(strings.Join(lines, "\n")), 0666)

		if err != nil {
			return err
		}

		files = append(files, NewReader(strings.NewReader(strings.Join(lines, "\n"))))
	}

	err := Merge(lineWriter, files...)

	if err != nil {
		return err
	}

	return nil
}
