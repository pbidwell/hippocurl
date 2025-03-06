package modules

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
)

// APIModule implements the HippoModule interface
type APIModule struct{}

var _ HippoModule = (*APIModule)(nil)

func (a APIModule) Name() string {
	return "api"
}

func (a APIModule) Description() string {
	return "Performs HTTP requests and displays the response."
}

func (a APIModule) Logo() string {
	return "L"
}

func (a APIModule) Execute(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: hc api <url> [method=GET] [body]")
		return
	}

	url := args[0]
	method := "GET"
	body := ""

	if len(args) > 1 {
		method = args[1]
	}
	if len(args) > 2 {
		body = args[2]
	}

	performHTTPRequest(url, method, body)
}

func performHTTPRequest(url, method, body string) {
	spinner := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	spinner.Start()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		spinner.Stop()
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		spinner.Stop()
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		spinner.Stop()
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	spinner.Stop()
	fmt.Println("\n=======================")
	fmt.Println("    HTTP RESPONSE    ")
	fmt.Println("=======================")
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Headers: %v\n", resp.Header)
	fmt.Println("Body:")
	fmt.Println(string(bodyBytes))
}
