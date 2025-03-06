package modules

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rodaine/table"
)

type ExploreModule struct{}

var _ HippoModule = (*ExploreModule)(nil)

func (e ExploreModule) Name() string {
	return "explore"
}

func (e ExploreModule) Description() string {
	return "Profiles a given hostname or IP address, fetching DNS records, geolocation, and port scan data."
}

func (e ExploreModule) Execute(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: hc explore [hostname/IP]")
		return
	}
	host := args[0]

	explore(host)
}

func (e ExploreModule) Logo() string {
	return "🔍"
}

func explore(host string) {
	spinner := spinner.New(spinner.CharSets[35], 100*time.Millisecond)

	// Fetch DNS records first
	fmt.Println()
	fmt.Println("=======================")
	fmt.Println("      DNS RECORDS ")
	fmt.Println("=======================")
	spinner.Start()
	dnsTable := fetchDNSRecords(host)
	spinner.Stop()
	dnsTable.Print()

	ips, err := net.LookupIP(host)
	if err != nil {
		fmt.Printf("Error resolving host: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("=======================")
	fmt.Println("      SERVER SCANS ")
	fmt.Println("=======================")
	spinner.Start()
	ipTbl := table.New(
		"[IP Address]",
		"[Country]",
		"[Region]",
		"[City]",
		"[HTTP]",
		"[HTTPS]",
		"[SSH]",
		"[SFTP]",
	)

	for _, ip := range ips {
		geolocateIP(ip.String(), ipTbl)
	}

	spinner.Stop()
	ipTbl.Print()
}

func geolocateIP(ip string, ipTbl table.Table) {
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
	ipTbl.AddRow(
		ip,
		fmt.Sprintf("%v", data["country"]),
		fmt.Sprintf("%v", data["region"]),
		fmt.Sprintf("%v", data["city"]),
		openPorts[80],
		openPorts[443],
		openPorts[22],
		openPorts[115],
	)
}

func fetchDNSRecords(host string) table.Table {
	tbl := table.New(
		"[CNAME]",
		"[NS Records]",
	)

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
