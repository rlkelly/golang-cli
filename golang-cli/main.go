package main

import (
  "bufio"
	"bytes"
  "fmt"
  "os"
  "os/exec"
  "strings"
	"text/tabwriter"
  "time"
)

const shift = uint(^uint(0)>>63) << 4
const PCT = 0x8000 << shift

var Screen *bytes.Buffer = new(bytes.Buffer)
// var Output *bufio.Writer = bufio.NewWriter(os.Stdout)
var Output = os.Stdout
var first = true

type Table struct {
	tabwriter.Writer

	Buf *bytes.Buffer
}

func NewTable(minwidth, tabwidth, padding int, padchar byte, flags uint) *Table {
	tbl := new(Table)
	tbl.Buf = new(bytes.Buffer)
	tbl.Init(tbl.Buf, minwidth, tabwidth, padding, padchar, flags)

	return tbl
}

func (t *Table) String() string {
	t.Flush()
	return t.Buf.String()
}

func Flush() {
	for idx, str := range strings.Split(Screen.String(), "\n") {
		if idx > Height() {
			return
		}

		Output.WriteString(str + "\n")
	}

	// Output.Flush()
	Screen.Reset()
}

func Clear() {
	Output.WriteString("\033[2J")
}

func MoveCursor(x int, y int) {
	fmt.Fprintf(Output, "\033[%d;%dH", y, x)
}

func Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(Output, a...)
}

func Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(Output, a...)
}

func PrintWithOffset(offset int, a ...interface{}) (n int, err error) {
  return fmt.Fprintln(Output, a...)
}

func handleResponse(response string, box *Box, writer *Writer) {
  switch response {
  case "test":
    // Output.Flush()
    fillScreen("Z")
    mainMenu(response, box, writer)
  case "exit":
    // Output.Flush()
    MoveCursor(1, 1)
    // clearScreen()
    Clear()
    c := exec.Command("clear")
    c.Stdout = os.Stdout
    c.Run()
    os.Exit(1)
  case "1":
    first = false
    MoveCursor(2, 2)
    slowPrinter("Select from the menu:", 125, first)
    MoveCursor(2, 4)
    slowPrinter("1)  First", 50, first)
    MoveCursor(2, 5)
    slowPrinter("2)  Second", 50, first)

    MoveCursor(2, 10)
    fmt.Fprintln(Output, strings.Repeat(" ", 25))
    MoveCursor(2, 10)
    Print("set> ")
    first = true

  default:
    mainMenu(response, box, writer)
  }
}

func fillScreen(value string) {
  MoveCursor(1, 1)
  if len(value) > 1 {
    return
  }
  for i:=1; i < Height(); i++ {
    Println(strings.Repeat(value, Width()))
  }
}

func clearScreen() {
  MoveCursor(1, 1)
  for i:=1; i < Height(); i++ {
    Println(strings.Repeat(" ", Width()))
  }
}

func mainMenu(text string, box *Box, writer *Writer) {
  MoveCursor(1, 1)
  Print(box.String())

  MoveCursor(2, 2)
  slowPrinter("Select from the menu:", 125, first)
  MoveCursor(2, 4)
  slowPrinter("1)  Social-Engineer", 50, first)
  MoveCursor(2, 5)
  slowPrinter("2)  Other Stuff", 50, first)
  MoveCursor(2, 6)
  slowPrinter("3)  Third Party Modules", 50, first)

  MoveCursor(2, 10)
  fmt.Fprintln(Output, strings.Repeat(" ", 25))
  MoveCursor(2, 10)
  Print("set> ")
}

func slowPrinter(text string, spacing int, first bool) {
  if first {
    for _, r := range text {
      Print(string(r))
      time.Sleep(time.Millisecond * time.Duration(spacing))
    }
    time.Sleep(time.Millisecond * time.Duration(spacing))
    return
  }
  Print(text)
}

func main() {
  scanner := bufio.NewScanner(os.Stdin)
  Clear()
  // Output.Flush()
  box := NewBox(Width() / 2, Height() - 1)
  writer := New()
  writer.Start()

  // for i := 0; i <= 100; i++ {
  //   fmt.Fprintf(writer, "Loading.. (%d/%d) GB\n", i, 100)
  //   time.Sleep(time.Millisecond * 25)
  // }

  for {
    MoveCursor(1, 1)
    handleResponse(scanner.Text(), box, writer)

    MoveCursor(7, 10)

    scanner.Scan()
    clearScreen()
    first = false
  }
}
