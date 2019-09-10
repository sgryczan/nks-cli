package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/viper"
)

type PrintObjInput struct {
	Properties []string
	Objects    interface{}
}

func httpRequest(method string, url string) ([]byte, error) {
	var httpClient = &http.Client{}
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", "Bearer "+viper.GetString("api_token"))

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

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

func PrintObj(i *PrintObjInput) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	fields := i.Properties
	s := ""
	for _, f := range fields {
		s += f + "\t"
	}
	s += "\n"
	fmt.Fprintf(w, s)

	return nil
}
