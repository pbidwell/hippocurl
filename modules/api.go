package modules

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hippocurl/utils"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

// APIModule implements the HippoModule interface
type APIModule struct{}

func (a APIModule) Name() string {
	return "api"
}

func (a APIModule) Description() string {
	return "Performs HTTP requests."
}

func (a APIModule) Execute(ctx context.Context, args []string) {
	config, ok := ctx.Value(utils.ConfigKey).(*utils.Config)
	if !ok || len(config.Services) == 0 {
		fmt.Println("No services configured. Please check your configuration file.")
		return
	}

	var serviceName, routeName, envName string
	if len(args) > 0 {
		serviceName = args[0]
	}
	if len(args) > 1 {
		routeName = args[1]
	}
	if len(args) > 2 {
		envName = args[2]
	}

	service, route, env := getServiceDetails(config, serviceName, routeName, envName)
	if service == nil || route == nil || env == nil {
		fmt.Println("Invalid selection.")
		return
	}

	url := env.BaseURL + route.Path
	performHTTPRequest(url, route.Method, "")
}

func (e APIModule) Logo() string {
	return "📤"
}

func performHTTPRequest(url, method, body string) {
	var reqBody *bytes.Reader
	if body != "" {
		reqBody = bytes.NewReader([]byte(body))
	} else {
		reqBody = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Println("\n=======================")
	fmt.Println("    HTTP RESPONSE    ")
	fmt.Println("=======================")
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Headers: %v\n", resp.Header)
	fmt.Println("Body (Raw Response):")
	printFormattedResponse(bodyBytes, resp.Header.Get("Content-Type"))
}

func getServiceDetails(config *utils.Config, serviceName, routeName, envName string) (*utils.Service, *utils.Route, *utils.Environment) {
	if serviceName == "" || routeName == "" || envName == "" {
		return promptUserForServiceDetails(config)
	}

	for _, service := range config.Services {
		if service.Name == serviceName {
			for _, route := range service.Routes {
				if route.Name == routeName {
					for _, env := range service.Environments {
						if env.Name == envName {
							return &service, &route, &env
						}
					}
				}
			}
		}
	}
	return nil, nil, nil
}

func promptUserForServiceDetails(config *utils.Config) (*utils.Service, *utils.Route, *utils.Environment) {
	servicePrompt := promptui.Select{
		Label: "Select a Service",
		Items: config.GetServiceNames(),
	}
	_, serviceName, err := servicePrompt.Run()
	if err != nil {
		fmt.Println("Selection cancelled.")
		return nil, nil, nil
	}

	service := config.GetServiceByName(serviceName)

	routePrompt := promptui.Select{
		Label: "Select a Route",
		Items: service.GetRouteNames(),
	}
	_, routeName, err := routePrompt.Run()
	if err != nil {
		fmt.Println("Selection cancelled.")
		return nil, nil, nil
	}

	route := service.GetRouteByName(routeName)

	envPrompt := promptui.Select{
		Label: "Select an Environment",
		Items: service.GetEnvironmentNames(),
	}
	_, envName, err := envPrompt.Run()
	if err != nil {
		fmt.Println("Selection cancelled.")
		return nil, nil, nil
	}

	environment := service.GetEnvironmentByName(envName)

	return service, route, environment
}
func printFormattedResponse(body []byte, contentType string) {
	switch {
	case strings.Contains(contentType, "json"):
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, body, "", "  ")
		if err == nil {
			fmt.Println(prettyJSON.String())
		} else {
			fmt.Println(string(body))
		}
	case strings.Contains(contentType, "xml"):
		var prettyXML []byte
		err := xml.Unmarshal(body, &prettyXML) // Unmarshal XML first
		if err == nil {
			formattedXML, _ := xml.MarshalIndent(prettyXML, "", "  ")
			fmt.Println(string(formattedXML))
		} else {
			fmt.Println(string(body))
		}
	default:
		fmt.Println(string(body))
	}
}
