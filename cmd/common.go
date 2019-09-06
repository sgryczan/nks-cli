package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func readApiToken() string {
	fmt.Printf("Enter your NKS API Token:")
	reader := bufio.NewReader(os.Stdin)
	token, _ := reader.ReadString('\n')

	r := strings.NewReplacer("\r", "", "\n", "")
	// convert CRLF to LF
	token = r.Replace(token)

	return token
}
