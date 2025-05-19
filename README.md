# **HippoCurl - A Modular HTTP & Networking Utility** 🦛🌐

HippoCurl (`hc`) is a modular command-line utility designed for HTTP requests, network diagnostics, and security analysis. It provides a framework where users can run modules to perform API requests, DNS lookups, port scanning, geolocation analysis, etc.

## **Key Features**
 **Modular Architecture** – Extendable with custom modules  
 **HTTP & API Interaction** – Send requests, inspect responses, and analyze APIs  
 **Networking Tools** – DNS record lookup, reverse DNS, IP geolocation, Whois  
 **Security Insights** – Open port scanning, SSL/TLS certificate analysis  
 **CLI-Friendly** – Supports structured output, formatted tables, and spinners for better UX  

**Note:** HippoCurl is something I write in my free time to fit the way I like to interact with network entities. It's a simple utility that generally just works, but as with any weekend project: Mandatory "YMMV" warning.

Feel free to explore the existing modules and contribute to the project.

## Installation

### **Using Prebuilt Binaries**
You can download the latest version of HippoCurl from the [Releases](https://github.com/pbidwell/hippocurl/releases) page.

1. **Download the appropriate binary** for your operating system:
   - **Linux**: `hippocurl-linux-amd64`
   - **macOS**: `hippocurl-mac-amd64`
   - **Windows**: `hippocurl-windows-amd64.exe`
2. **Make it executable (Linux/macOS)**:
   ```sh
   chmod +x hippocurl-linux-amd64
   ```
3. **Move it to a directory in your PATH**:
   ```sh
   sudo mv hippocurl-linux-amd64 /usr/local/bin/hc
   ```
4. **Run HippoCurl**:
   ```sh
   hc
   ```

### **Building from Source**
If you prefer, you can build HippoCurl from source.

1. **Clone the repository**:
   ```sh
   git clone https://github.com/pbidwell/hippocurl.git
   cd hippocurl
   ```
2. **Ensure you have Go installed** (version 1.20+):
   ```sh
   go version
   ```
3. **Build the binary**:
   ```sh
   go build -o hc ./cmd/main.go
   ```
4. **Move it to your PATH** (optional):
   ```sh
   sudo mv hc /usr/local/bin/
   ```
5. **Verify installation**:
   ```sh
   hc
   ```

---

## Usage
HippoCurl (`hc`) is a command-line tool designed to simplify HTTP requests, API interactions, and service explorations.

### Basic Command Structure
```
hc <module> [arguments]
```
- `<module>`: The name of the module you wish to execute (e.g., `api`, `explore`, `log`).
- `[arguments]`: Optional parameters that vary by module.

### API Requests
To make an API request using a configured service:
```
hc api <service> <route> <environment>
```
If any of the parameters are omitted, HippoCurl will interactively prompt you to select the desired service, route, and environment.
Example:
```
hc api UserService GetUser Development
```
This will:
- Fetch the `GetUser` route from the `UserService` in the `Development` environment.
- Use the configured base URL, headers, and authentication.
- Display the response in a structured format.

#### API Configuration File (`~/.hc/api_config.yml`)
HippoCurl uses a YAML configuration file to define reusable HTTP services, their environments, authentication settings, and request routes. This allows you to interact with APIs using simple commands like:
```sh
hc api GitHubAPI get-user production
```
HippoCurl looks for the configuration file in the following location:
```
~/.hc/api_config.yml
```
##### Example Configuration
```yaml
services:
  - name: HttpBin
    environments:
      - name: default
        base_url: "https://httpbin.org"
        auth:
          type: "none"
        headers:
          Content-Type: "application/json"

    routes:
      - name: post-json
        description: "POST JSON test payload"
        method: POST
        path: "/post"
        body: '{"hippo": "rules"}'

      - name: get-ip
        description: "Get your IP address"
        method: GET
        path: "/ip"
        body: ""

  - name: DuckDuckGo
    environments:
      - name: default
        base_url: "https://duckduckgo.com"
        auth:
          type: "none"

    routes:
      - name: homepage
        description: "Fetch DuckDuckGo homepage"
        method: GET
        path: "/"
        body: ""
```


##### Top-Level Structure
The config file defines a list of `services`, each with their own `environments` and `routes`.

###### Service Block
Each service entry represents a logical grouping of API routes:
```yaml
services:
  - name: GitHubAPI
```
- **name**: A unique identifier for the service.

###### Environments
Each service can have one or more environments (e.g., production, staging):
```yaml
    environments:
      - name: production
        base_url: "https://api.github.com"
```
- **name**: Environment name (e.g., `production`, `default`)
- **base_url**: Base URL used for all routes under this environment.
- **auth**: Authentication details for this environment.
  - `type`: One of `none`, `basic`, or `bearer`
  - `token`, `username`, `password`: Depending on the auth type
- **headers** *(optional)*: Custom headers to include with all requests (e.g., content type, user agent)

###### Routes
Each route represents a specific API endpoint:
```yaml
    routes:
      - name: get-user
        description: "Fetch authenticated user info"
        method: GET
        path: "/user"
        body: ""
```
- **name**: Short route name used from the CLI
- **description**: What the route does
- **method**: HTTP method (`GET`, `POST`, etc.)
- **path**: URL path appended to `base_url`
- **body**: Optional JSON payload for POST/PUT requests

---
##### Services in Sample Config
- `GitHubAPI`: Uses bearer token auth to interact with GitHub
- `HttpBinTest`: Great for testing; includes a JSON POST and an IP-fetching GET
- `DuckDuckGo`: A basic GET request to the homepage without auth

---
#### Tip
You can quickly test any route like this (using sample config for illustration):
```sh
hc api HttpBinTest post-json default
```
This will perform a POST request to `https://httpbin.org/post` with the predefined JSON body.

### Exploring Hosts
```
hc explore <hostname or IP>
```
Example:
```
hc explore example.com
```
This will:
- Resolve DNS records.
- Retrieve IP-based geolocation.
- Perform common port scans.
- Detect SSL/TLS certificate details.

### Viewing Logs
```
hc log
```
This command displays:
- The location of the log file.
- The last 100 lines of logs.

---
## Feature Ideas
- [x] Modular design allowing easy extension
- [x] CLI-based utility with structured output
- [x] Configuration file support (`.hc`)
- [ ] Configuration wizard support as new module + curl command converter
- [ ] Postman config conversion support
- [ ] Support for automated API test suite with HTTP response code + response header assertions

