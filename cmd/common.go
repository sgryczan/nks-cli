package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

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
