package config

const APIConfigSampleYaml = `
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
`
