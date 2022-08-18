package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/nsf/jsondiff"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err.Error())
	}
}

func run() error {
	cmd := exec.Command("go", "test", "-json", "./test")

	sp, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to open stdout pipe: %w", err)
	}
	ep, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to open stderr pipe: %w", err)
	}

	r := io.MultiReader(sp, ep)

	processResult := make(chan error)
	go func() {
		processResult <- processOutput(r)
		close(processResult)
	}()

	err = cmd.Run()
	if processError := <-processResult; processError != nil {
		return processError
	}
	return err
}

func processOutput(r io.Reader) error {
	s := bufio.NewScanner(r)

	var ev testEvent
	c := newSnapshotChecker()
	for s.Scan() {
		if err := json.Unmarshal(s.Bytes(), &ev); err != nil {
			return fmt.Errorf("unmarshal test event %q: %w", s.Text(), err)
		}

		if err := c.Handle(ev); err != nil {
			return fmt.Errorf("handle test event: %w", err)
		}
	}

	return s.Err()
}

type snapshotChecker struct {
	state       func(testEvent) error
	actualJSON  *bytes.Buffer
	currentTest string
}

func newSnapshotChecker() *snapshotChecker {
	c := &snapshotChecker{
		actualJSON: &bytes.Buffer{},
	}
	c.state = c.waitTestRun
	return c
}

func (c *snapshotChecker) Handle(ev testEvent) error {
	return c.state(ev)
}

func (c *snapshotChecker) waitTestRun(ev testEvent) error {
	if ev.Action == testEventActionOutput && strings.Contains(ev.Output, "=== RUN") {
		c.state = c.waitActualJSONStart
	}
	return nil
}

func (c *snapshotChecker) waitActualJSONStart(ev testEvent) error {
	if ev.Action == testEventActionOutput && strings.Contains(ev.Output, "Actual JSON") {
		c.state = c.collectActualJSONFirstLine
	}
	return nil
}

func (c *snapshotChecker) collectActualJSONFirstLine(ev testEvent) error {
	if ev.Action == testEventActionOutput && strings.Contains(ev.Output, "{") {
		c.actualJSON.Reset()
		c.actualJSON.WriteRune('{')
		c.actualJSON.WriteRune('\n')
		c.currentTest = ev.Test
		c.state = c.collectActualJSON
	}
	return nil
}

func (c *snapshotChecker) collectActualJSON(ev testEvent) error {
	switch ev.Action {
	case testEventActionOutput:
		c.actualJSON.WriteString(strings.TrimPrefix(ev.Output, "        "))
	case testEventActionRun:
		c.state = c.waitTestRun
		if c.actualJSON.Len() != 0 {
			path := strings.ReplaceAll(c.currentTest, "TestJDocExchange", "testdata")
			expectedJSON, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read file %q: %w", path, err)
			}

			opts := jsondiff.DefaultConsoleOptions()
			difference, str := jsondiff.Compare(
				expectedJSON,
				c.actualJSON.Bytes(),
				&opts,
			)
			if difference == jsondiff.FullMatch {
				break
			}

			border := strings.Repeat("*", len(c.currentTest))

			fmt.Println(border)
			fmt.Println(c.currentTest)
			fmt.Println(border)
			fmt.Println("Diff:")
			fmt.Println()
			fmt.Println(str)
			fmt.Println()
			fmt.Println(border)
			fmt.Println(c.currentTest)
			fmt.Println("Update snapshot (y/N): ")

			should, err := shouldUpdateSnapshot()
			if err != nil {
				return fmt.Errorf("should update snapshot: %w", err)
			}

			if !should {
				break
			}

			if err := ioutil.WriteFile(path, c.actualJSON.Bytes(), 0644); err != nil { //nolint:gosec // It's safe to have 0644 permission.
				return fmt.Errorf("write file %q: %w", path, err)
			}
		}
	}
	return nil
}

func shouldUpdateSnapshot() (bool, error) {
	buf := bufio.NewReader(os.Stdin)
	sentence, err := buf.ReadBytes('\n')
	if err != nil {
		return false, err
	}

	return bytes.EqualFold(sentence, []byte("y\n")), nil
}

// testEvent represent structure of test event.
// Read `go doc test2json` for reference.
type testEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

const (
	testEventActionOutput = "output"
	testEventActionRun    = "run"
)
