package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ui(reader io.Reader) {
	scanner := bufio.NewReader(reader)
	for {
		fmt.Println("Please enter the name of the file you want to use followed by the length of words you want in your new list")
		text, err := scanner.ReadString('\n')
		text = strings.TrimSpace(text)

		if err != io.EOF && len(text) > 0 {
			cmd := strings.Split(text, " ")
			if len(cmd) != 2 {
				fmt.Println("Please start again and enter both parameters")
				continue
			}
			cmdName := cmd[0]
			cmdLength, err := strconv.Atoi(cmd[1])
			if err != nil {
				fmt.Println("Please enter a length in integer format")
			}
			begin(cmdName, cmdLength)
		}
	}
}
