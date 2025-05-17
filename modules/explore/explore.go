package explore

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/pbidwell/hippocurl/internal/config"
	"github.com/pbidwell/hippocurl/utils"

	"github.com/briandowns/spinner"
	"github.com/rodaine/table"
)

type ExploreModule struct{}

var elogger *log.Logger

func (e ExploreModule) Name() string {
	return "explore"
}

func (e ExploreModule) Description() string {
	return "Profiles a given hostname or IP address, fetching DNS records, geolocation, and port scan data."
}

func (e ExploreModule) Use() string {
	return fmt.Sprintf("%s <hostname>", e.Name())
}

func (e ExploreModule) Execute(app *config.App, args []string) {
	utils.Print(e.Name(), utils.ModuleTitle)

	if len(args) != 1 {
		utils.Print("Usage: hc explore [hostname/IP]\n", utils.NormalText)
		return
	}
	host := args[0]
	elogger = app.Logger

	explore(host)
}

func (e ExploreModule) Logo() string {
	return "🔍"
}

func explore(host string) {
	spinner := spinner.New(spinner.CharSets[35], 100*time.Millisecond)

	// Fetch DNS records
	utils.Print("DNS Records", utils.Header1)
	spinner.Start()
	dnsTable := fetchDNSRecords(host)
	spinner.Stop()
	dnsTable.Print()

	ips, err := net.LookupIP(host)
	if err != nil {
		elogger.Printf("Error resolving host: %v\n", err)
	}

	filteredIPs := make([]net.IP, 0, len(ips)) // Pre-allocate capacity
	for _, ip := range ips {
		if ip.To4() != nil {
			filteredIPs = append(filteredIPs, ip)
		}
	}

	utils.Print("Server Scans", utils.Header1)
	spinner.Start()
	serverTbl := table.New(
		"[IP Address]",
		"[Country]",
		"[Region]",
		"[City]",
		"[HTTP]",
		"[HTTPS]",
		"[SSH]",
		"[SFTP]",
		"[SSL Issuer]",
		"[SSL Expiry]",
	)

	for _, ip := range filteredIPs {
		geolocateIP(ip.String(), serverTbl)
	}

	spinner.Stop()
	serverTbl.Print()
}

func geolocateIP(ip string, serverTbl table.Table) {
	url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching geolocation data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Printf("Error decoding geolocation response: %v\n", err)
		return
	}

	openPorts := scanOpenPorts(ip)
	issuer, expiry := fetchSSLCertificate(ip)
	serverTbl.AddRow(
		ip,
		fmt.Sprintf("%v", data["country"]),
		fmt.Sprintf("%v", data["region"]),
		fmt.Sprintf("%v", data["city"]),
		openPorts[80],
		openPorts[443],
		openPorts[22],
		openPorts[115],
		issuer,
		expiry,
	)
}

func fetchDNSRecords(host string) table.Table {
	tbl := table.New("[CNAME]", "[NS Records]")

	var cnameStr, nsStr string

	cname, err := net.LookupCNAME(host)
	if err == nil {
		cnameStr = cname
	}

	nsRecords, err := net.LookupNS(host)
	if err == nil {
		for _, ns := range nsRecords {
			if nsStr == "" {
				nsStr = ns.Host
			} else {
				nsStr += ", " + ns.Host
			}
		}
	}

	tbl.AddRow(cnameStr, nsStr)
	// log.Printf("DNS Records for %s - CNAME: %s, NS: %s", host, cnameStr, nsStr)
	return tbl
}

func scanOpenPorts(host string) map[int]string {
	ports := []int{22, 80, 443, 115} // Common ports for SSH, HTTP, HTTPS, and SFTP
	result := make(map[int]string)

	for _, port := range ports {
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err != nil {
			result[port] = "Closed"
		} else {
			result[port] = "Open"
			conn.Close()
		}
	}

	return result
}

func fetchSSLCertificate(host string) (string, string) {
	// Define TLS configuration to skip verification for IP addresses
	tlsConfig := &tls.Config{
		InsecureSkipVerify: net.ParseIP(host) != nil, // Skip verification only if it's an IP
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", host), tlsConfig)
	if err != nil {
		elogger.Printf("Error fetching SSL Certificate for host %s: %v", host, err)
		return "N/A", "N/A"
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	issuer := cert.Issuer.CommonName
	expiry := cert.NotAfter.Format("2006-01-02 15:04:05")

	return issuer, expiry
}
