package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	ct "github.com/daviddengcn/go-colortext"
)

type OutletFactory struct {
	Padding int

	sync.Mutex
}

var colors = []ct.Color{
	ct.Cyan,
	ct.Yellow,
	ct.Green,
	ct.Magenta,
	ct.Red,
	ct.Blue,
}

func NewOutletFactory() (of *OutletFactory) {
	return new(OutletFactory)
}

func (of *OutletFactory) LineReader(wg *sync.WaitGroup, name string, index int, r io.Reader, isError bool) {
	defer wg.Done()

	color := colors[index%len(colors)]

	reader := bufio.NewReader(r)

	var buffer bytes.Buffer

	for {
		buf := make([]byte, 1024)

		if n, err := reader.Read(buf); err != nil {
			return
		} else {
			buf = buf[:n]
		}

		for {
			i := bytes.IndexByte(buf, '\n')
			if i < 0 {
				break
			}
			buffer.Write(buf[0:i])
			of.WriteLine(name, buffer.String(), color, ct.None, isError)
			buffer.Reset()
			buf = buf[i+1:]
		}

		buffer.Write(buf)
	}
}

func (of *OutletFactory) SystemOutput(str string) {
	of.WriteLine("forego", str, ct.White, ct.None, false)
}

func (of *OutletFactory) ErrorOutput(str string) {
	fmt.Printf("ERROR: %s\n", str)
	os.Exit(1)
}

// Write out a single coloured line
func (of *OutletFactory) WriteLine(left, right string, leftC, rightC ct.Color, isError bool) {
	of.Lock()
	defer of.Unlock()

	ct.ChangeColor(leftC, true, ct.None, false)
	formatter := fmt.Sprintf("%%-%ds | ", of.Padding)
	fmt.Printf(formatter, left)

	words := strings.Fields(right)
	if len(words) > 2 {
		word := words[2]
		switch words[2][:len(word)-1] {
		case "INFO":
			ct.ChangeColor(ct.Green, true, ct.None, true)
		case "WARNING":
			ct.ChangeColor(ct.Yellow, true, ct.None, true)
		case "ERROR":
			ct.ChangeColor(ct.Red, true, ct.None, true)
		case "DEBUG":
			ct.ChangeColor(ct.White, true, ct.Magenta, true)
		default:
			ct.ResetColor()
		}
	} else {
		ct.ResetColor()
	}

	fmt.Println(right)

	ct.ResetColor()
}
