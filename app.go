package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Function to perform an HTTP request with the given payload and parameters
func testSSTI(targetURL, method, data, payload string, headers map[string]string, proxy string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Configure proxy if provided
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return "", fmt.Errorf("invalid proxy URL: %v", err)
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	var req *http.Request
	var err error

	// Replace "FUZZ" with the payload in URL or data
	if method == "GET" {
		if strings.Contains(targetURL, "FUZZ") {
			targetURL = strings.Replace(targetURL, "FUZZ", url.QueryEscape(payload), -1)
		} else {
			separator := "?"
			if strings.Contains(targetURL, "?") {
				separator = "&"
			}
			targetURL = fmt.Sprintf("%s%spayload=%s", targetURL, separator, url.QueryEscape(payload))
		}
		req, err = http.NewRequest("GET", targetURL, nil)
	} else if method == "POST" {
		if strings.Contains(data, "FUZZ") {
			data = strings.Replace(data, "FUZZ", payload, -1)
		} else {
			if data != "" {
				data += "&"
			}
			data += "payload=" + url.QueryEscape(payload)
		}
		body := strings.NewReader(data)
		req, err = http.NewRequest("POST", targetURL, body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		return "", fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Add provided headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Execute the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(respBody), nil
}

// Function to check if the response contains the expected result
func checkVulnerability(resp, expectedResult, engine, payload string) bool {
	cleanResp := strings.ReplaceAll(resp, " ", "")
	if strings.Contains(cleanResp, expectedResult) {
		fmt.Printf("[+] Vulnerable to %s SSTI with payload: %s\n", engine, payload)
		return true
	}
	return false
}

// Function to test a series of SSTI payloads
func testSSTIChain(url, method, data string, sstiChain map[string]string, headers map[string]string, proxy string) {
	vulnerable := false

	for engine, payloadTemplate := range sstiChain {
		uniqueNum := (time.Now().UnixNano() % 10000) + 10000 // Ensures a 5-digit number
		expectedResult := fmt.Sprintf("%d", uniqueNum*2)

		// Replace placeholders in the payload
		payload := strings.Replace(payloadTemplate, "UNIQUE_NUM", fmt.Sprintf("%d", uniqueNum), -1)

		resp, err := testSSTI(url, method, data, payload, headers, proxy)
		if err != nil {
			fmt.Printf("[-] Error testing %s: %v\n", engine, err)
			continue
		}

		if checkVulnerability(resp, expectedResult, engine, payload) {
			vulnerable = true
		} else {
			fmt.Printf("[-] Not vulnerable to %s SSTI.\n", engine)
		}
	}

	if !vulnerable {
		fmt.Println("[-] No SSTI vulnerabilities detected.")
	}
}

func main() {
	// Command-line configurations
	urlPtr := flag.String("u", "", "Target URL")
	methodPtr := flag.String("x", "GET", "HTTP method (GET or POST)")
	dataPtr := flag.String("d", "", "Data for POST requests, use FUZZ as placeholder for payload")
	headersPtr := flag.String("H", "", "Additional headers in 'Header1:Value1,Header2:Value2' format")
	proxyPtr := flag.String("proxy", "", "Proxy to use for requests (e.g., http://127.0.0.1:8080)")
	payloadsFilePtr := flag.String("p", "payloads.json", "Path to JSON file containing SSTI payloads")
	flag.Parse()

	if *urlPtr == "" {
		fmt.Println("Please specify a URL with the -u parameter.")
		return
	}

	// Parse headers
	headers := make(map[string]string)
	if *headersPtr != "" {
		parts := strings.Split(*headersPtr, ",")
		for _, part := range parts {
			header := strings.SplitN(part, ":", 2)
			if len(header) == 2 {
				headers[strings.TrimSpace(header[0])] = strings.TrimSpace(header[1])
			}
		}
	}

	// Load SSTI payloads from external JSON file
	payloadsData, err := ioutil.ReadFile(*payloadsFilePtr)
	if err != nil {
		fmt.Printf("Failed to read payloads file: %v\n", err)
		return
	}

	var sstiChain map[string]string
	err = json.Unmarshal(payloadsData, &sstiChain)
	if err != nil {
		fmt.Printf("Failed to parse payloads file: %v\n", err)
		return
	}

	// Test the chain of SSTI payloads
	testSSTIChain(*urlPtr, *methodPtr, *dataPtr, sstiChain, headers, *proxyPtr)
}
