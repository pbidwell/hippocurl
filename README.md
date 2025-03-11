# **HippoCurl - A Modular HTTP & Networking Utility** 🦛🌐

HippoCurl (`hc`) is a **powerful, modular command-line utility** designed for **HTTP requests, network diagnostics, and security analysis**. It provides an **extensible framework** where users can run various **modules** to perform **API requests, DNS lookups, port scanning, geolocation analysis, and more**.

## **Key Features**
✅ **Modular Architecture** – Easily extendable with custom modules  
✅ **HTTP & API Interaction** – Send requests, inspect responses, and analyze APIs  
✅ **Networking Tools** – DNS record lookup, reverse DNS, IP geolocation, Whois  
✅ **Security Insights** – Open port scanning, SSL/TLS certificate analysis  
✅ **AI-Enhanced** – API response summarization, anomaly detection, intelligent suggestions  
✅ **CLI-Friendly** – Supports structured output, formatted tables, and spinners for better UX  

## 🚧 Under Construction 🚧

HippoCurl is currently **under development**! ⚙️  
We are actively building new features, improving performance, and expanding its capabilities.  

Stay tuned for updates! 🚀  

In the meantime, feel free to explore the existing modules and contribute to the project.  

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
   git clone https://github.com/yourusername/hippocurl.git
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

### Configuration

HippoCurl uses a YAML-based configuration to define services, environments, routes, and authentication details. See the `hc_config.yml` file for customization.

```
~/.hcconfig/hc_config.yml
```

Modify this file to add new API services, routes, authentication methods, and custom headers.

#### Example Configuration

```yaml
services:
  - name: "User Service"
    environments:
      - name: "Development"
        base_url: "https://dev.api.example.com/user"
        auth:
          type: "bearer"
          token: "dev-secret-token"
        headers:
          Content-Type: "application/json"
          Authorization: "Bearer dev-secret-token"
      - name: "Production"
        base_url: "https://api.example.com/user"
        auth:
          type: "basic"
          username: "admin"
          password: "securepassword"
        headers:
          Content-Type: "application/json"
    routes:
      - name: "Get User"
        path: "/{id}"
        method: "GET"
        description: "Fetch user details by ID"
      - name: "Create User"
        description: "Create a new user"
        path: "/create"
        method: "POST"
        body: |
            {
                "name": "John Doe",
                "email": "johndoe@example.com",
                "favoriteanimal": "hippo"
            }
```

---

## Feature Ideas

### Core Features
- [x] Modular design allowing easy extension
- [x] CLI-based utility with structured output
- [x] Configuration file support (`.hcconfig`)
- [ ] Configuration wizard support as new module + curl command converter

### HTTP Features
- [x] Perform HTTP requests (GET, POST, etc.)
- [x] Support for custom headers & request bodies
- [ ] Response time measurement
- [x] Pretty-print JSON API responses

### Networking & Security
- [x] DNS record lookup (CNAME, NS, MX)
- [x] IP Geolocation lookup
- [x] Open port scanning (HTTP, SSH, SFTP, etc.)
- [x] Detect SSL/TLS certificate details
- [ ] Reverse DNS lookup
- [ ] Detect potential security risks in API responses
- [ ] Whois lookup for domains

### AI-Powered Enhancements
- [ ] AI-based API response summarization
- [ ] AI-powered content classification (Success, Error, Warning)
- [ ] Detect security vulnerabilities using AI
- [ ] NLP-based query support (e.g., `"Check status of example.com"`)

### Performance & Monitoring
- [ ] HTTP request benchmarking (track response times over multiple requests)
- [ ] Load testing for API endpoints
- [ ] Monitor uptime for a given endpoint
- [ ] Background monitoring with notifications

### Data & Exporting
- [ ] Save results to JSON or CSV format
- [ ] Generate reports from API responses
- [ ] Webhook integration to send results to Slack, Discord, etc.

### Utility & UX Enhancements
- [x] Interactive spinner for long-running tasks
- [ ] Auto-complete for module names
- [x] Consistent, color-coded output for better readability across modules
- [x] Easily-accessible log-viewing mode
- [ ] Interactive mode for easier configuration
- [ ] Config validation (no spaces in names, valid URLs, etc)

