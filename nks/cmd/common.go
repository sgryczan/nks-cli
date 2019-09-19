package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	nks "github.com/NetApp/nks-sdk-go/nks"
	"github.com/spf13/viper"
)

var flagClusterId int
var flagSolutionId int
var flagOrganizationId int
var flagSolutionName string
var flagSolutionRepoName string
var flagSolutionReleaseName string
var flagForce bool
var flagGenerateCompletionBash bool
var flagGenerateCompletionZsh bool
var FlagDebug bool

var SDKClient *nks.APIClient

func newClient() *nks.APIClient {
	client := nks.NewClient(viper.GetString("api_token"), viper.GetString("api_url"))
	return client
}

func httpRequest(method string, url string) ([]byte, error) {
	var httpClient = http.Client{}

	req, _ := http.NewRequest(method, url, nil)

	//fmt.Printf("http_request:\n%+v", req)

	req.Header.Add("Authorization", "Bearer "+viper.GetString("api_token"))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func httpRequestPost(method string, url string, b []byte) ([]byte, error) {
	var httpClient = http.Client{}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(b))
	//fmt.Printf("http_request:\n%+v", req)

	req.Header.Add("Authorization", "Bearer "+viper.GetString("api_token"))
	req.Header.Set("Content-Type", "application/json")
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
		os.Exit(1)
	}
}

func readApiToken() string {
	fmt.Printf("Enter your NKS API Token: ")
	reader := bufio.NewReader(os.Stdin)
	token, _ := reader.ReadString('\n')

	r := strings.NewReplacer("\r", "", "\n", "")
	// convert CRLF to LF
	token = r.Replace(token)

	return token
}
