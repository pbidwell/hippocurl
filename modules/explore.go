package modules

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"hippocurl/utils"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rodaine/table"
)

type ExploreModule struct{}

var logger *log.Logger

func (e ExploreModule) Name() string {
	return "explore"
}

func (e ExploreModule) Description() string {
	return "Profiles a given hostname or IP address, fetching DNS records, geolocation, and port scan data."
}

func (e ExploreModule) Execute(ctx context.Context, args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: hc explore [hostname/IP]")
		return
	}
	host := args[0]
	logger = ctx.Value(utils.LoggerKey).(*log.Logger)

	explore(host)
}

func (e ExploreModule) Logo() string {
	return "🔍"
}

func explore(host string) {
	spinner := spinner.New(spinner.CharSets[35], 100*time.Millisecond)

	// Fetch DNS records
	fmt.Println("\n=======================")
	fmt.Println("      DNS RECORDS ")
	fmt.Println("=======================")
	spinner.Start()
	dnsTable := fetchDNSRecords(host)
	spinner.Stop()
	dnsTable.Print()

	ips, err := net.LookupIP(host)
	filteredIPs := []net.IP{}
	for _, ip := range ips {
		if ip.To4() != nil { // Only include IPv4 addresses
			filteredIPs = append(filteredIPs, ip)
		}
	}
	if err != nil {
		fmt.Printf("Error resolving host: %v\n", err)
		return
	}

	fmt.Println("\n=======================")
	fmt.Println("      SERVER SCANS ")
	fmt.Println("=======================")
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

	// log.Printf("Scanned %s - Ports: HTTP:%s, HTTPS:%s, SSH:%s, SFTP:%s, SSL Issuer: %s, Expiry: %s",
	// 	ip, openPorts[80], openPorts[443], openPorts[22], openPorts[115], issuer, expiry)
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
		logger.Printf("Error fetching SSL Certificate for host %s: %v", host, err)
		return "N/A", "N/A"
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	issuer := cert.Issuer.CommonName
	expiry := cert.NotAfter.Format("2006-01-02 15:04:05")

	return issuer, expiry
}
