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

