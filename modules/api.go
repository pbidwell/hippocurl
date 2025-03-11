package modules

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hippocurl/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
)

// APIModule implements the HippoModule interface
type APIModule struct{}

var alogger *log.Logger

func (a APIModule) Name() string {
	return "api"
}

func (a APIModule) Description() string {
	return "Performs HTTP requests."
}

func (a APIModule) Execute(ctx context.Context, args []string) {
	// Module banner
	utils.Print(a.Name(), utils.ModuleTitle)

	config, ok := ctx.Value(utils.ConfigKey).(*utils.Config)
	if !ok || len(config.Services) == 0 {
		utils.Print("No services configured. Please check your configuration file.", utils.NormalText)
		return
	}

	alogger = ctx.Value(utils.LoggerKey).(*log.Logger)

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

	service, route, env, interactive := getServiceDetails(config, serviceName, routeName, envName)
	if service == nil || route == nil || env == nil {
		utils.Print("Invalid selection.", utils.NormalText)
		return
	}

	url := env.BaseURL + route.Path
	performHTTPRequest(url, route.Method, env.Headers, route.Body)
	if interactive {
		utils.Print(fmt.Sprintf("Use \"hc %s %s %s %s\" to re-try this API call.", a.Name(), service.Name, route.Name, env.Name), utils.Hint)
	}
}

func (e APIModule) Logo() string {
	return "📤"
}

func performHTTPRequest(url string, method string, headers map[string]string, body string) {
	spinner := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
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

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	utils.Print("HTTP Request", utils.Header1)
	utils.Print("Headers", utils.Header2)
	utils.PrintHeaders(req.Header)
	utils.Print("Body", utils.Header2)
	reqBodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		alogger.Printf("Error parsing request body: %v\n", err)
		fmt.Printf("Error parsing request body: %v\n", err)
		return
	}
	printFormattedResponse(reqBodyBytes, req.Header.Get("Content-Type"))
	spinner.Start()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	spinner.Stop()
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

	utils.Print("HTTP Response", utils.Header1)
	utils.Print("Status", utils.Header2)
	fmt.Println(resp.Status)
	utils.Print("Headers", utils.Header2)
	utils.PrintHeaders(resp.Header)
	utils.Print("Body", utils.Header2)
	printFormattedResponse(bodyBytes, resp.Header.Get("Content-Type"))
}

func getServiceDetails(config *utils.Config, serviceName, routeName, envName string) (*utils.Service, *utils.Route, *utils.Environment, bool) {
	serviceMap := make(map[string]*utils.Service)
	for i := range config.Services {
		serviceMap[config.Services[i].Name] = &config.Services[i]
	}

	if serviceName == "" || routeName == "" || envName == "" {
		service, route, env := promptUserForServiceDetails(config)
		return service, route, env, true
	}

	service, exists := serviceMap[serviceName]
	if !exists {
		return nil, nil, nil, false
	}

	routeMap := make(map[string]*utils.Route)
	for i := range service.Routes {
		routeMap[service.Routes[i].Name] = &service.Routes[i]
	}

	route, exists := routeMap[routeName]
	if !exists {
		return nil, nil, nil, false
	}

	envMap := make(map[string]*utils.Environment)
	for i := range service.Environments {
		envMap[service.Environments[i].Name] = &service.Environments[i]
	}

	env, exists := envMap[envName]
	if !exists {
		return nil, nil, nil, false
	}

	return service, route, env, false
}

func promptUserForServiceDetails(config *utils.Config) (*utils.Service, *utils.Route, *utils.Environment) {
	servicePrompt := promptui.Select{
		Label: "Select a Service",
		Items: config.GetServiceNames(),
	}
	_, serviceName, err := servicePrompt.Run()
	if err != nil {
		utils.Print("Selection cancelled.", utils.NormalText)
		return nil, nil, nil
	}

	service := config.GetServiceByName(serviceName)

	routePrompt := promptui.Select{
		Label: "Select a Route",
		Items: service.GetRouteNames(),
	}
	_, routeName, err := routePrompt.Run()
	if err != nil {
		utils.Print("Selection cancelled.", utils.NormalText)
		return nil, nil, nil
	}

	route := service.GetRouteByName(routeName)

	envPrompt := promptui.Select{
		Label: "Select an Environment",
		Items: service.GetEnvironmentNames(),
	}
	_, envName, err := envPrompt.Run()
	if err != nil {
		utils.Print("Selection cancelled.", utils.NormalText)
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
			utils.Print(prettyJSON.String(), utils.NormalText)
		} else {
			utils.Print(string(body), utils.NormalText)
		}
	case strings.Contains(contentType, "xml"):
	case strings.Contains(contentType, "html"):
		var prettyXML []byte
		err := xml.Unmarshal(body, &prettyXML) // Unmarshal XML first
		if err == nil {
			formattedXML, _ := xml.MarshalIndent(prettyXML, "", "  ")
			utils.Print(string(formattedXML), utils.NormalText)
		} else {
			utils.Print(string(body), utils.NormalText)
		}
	default:
		utils.Print(string(body), utils.NormalText)
	}
}
