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
var Output *bufio.Writer = bufio.NewWriter(os.Stdout)

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

	Output.Flush()
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

func handleResponse(response string, box *Box) {
  switch response {
  case "test":
    Output.Flush()
    fillScreen("Z")
    mainMenu(response, box)
  case "exit":
    Output.Flush()
    MoveCursor(1, 1)
    // clearScreen()
    Clear()
    c := exec.Command("clear")
    c.Stdout = os.Stdout
    c.Run()
    os.Exit(1)
  default:
    mainMenu(response, box)
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

func mainMenu(text string, box *Box) {
  MoveCursor(1, 1)
  Print(box.String())

  MoveCursor(2, 2)
  MoveCursor(2, 2)
  Print("Hello ", text)
  MoveCursor(2, 4)
  fmt.Fprintln(Output, strings.Repeat(" ", 25))
  MoveCursor(2, 4)
  Print("set> ")
}

func slowPrinter(text string, spacing int) {
  MoveCursor(6, 3)
  for _, r := range text {
    Print(string(r))
    // time.Sleep(time.Millisecond)
  }
}

func main() {
  scanner := bufio.NewScanner(os.Stdin)
  Clear()
  Output.Flush()
  box := NewBox(Width() / 2, Height() - 1)
  writer := New()
  writer.Start()

  for i := 0; i <= 100; i++ {
    fmt.Fprintf(writer, "Loading.. (%d/%d) GB\n", i, 100)
    time.Sleep(time.Millisecond * 25)
  }

  writer.Stop()

  for {
    MoveCursor(1, 1)
    handleResponse(scanner.Text(), box)
    Output.Flush()
    clearScreen()
    scanner.Scan()
  }
}
