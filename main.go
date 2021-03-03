package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	exeFilePath := "C:/Users/むずでょ/go/src/github.com/muzudho/go-echo-next-char/go-echo-next-char.exe"
	parameters := strings.Split("", " ")

	cmd := exec.Command(exeFilePath, parameters...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		panic(fmt.Errorf("cmd.Start() --> [%s]", err))
	}

	go receiveStdout(stdout)

	go receiveStdin(stdin)

	cmd.Wait()
}

// `epStdin` - External process stdin
func receiveStdin(epStdin io.WriteCloser) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		epStdin.Write([]byte(command))
		epStdin.Write([]byte("\n"))
	}
}

// `epStdout` - External process stdout
func receiveStdout(epStdout io.ReadCloser) {
	var buffer [1]byte // これが満たされるまで待つ。1バイト。

	p := buffer[:]

	for {
		n, err := epStdout.Read(p)

		if nil != err {
			if fmt.Sprintf("%s", err) == "EOF" {
				// End of file
				return
			}

			panic(err)
		}

		if 0 < n {
			bytes := p[:n]

			// 思考エンジンから１文字送られてくるたび、表示。
			print(string(bytes))
		}
	}
}
