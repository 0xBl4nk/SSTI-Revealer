package ssti

import (
    "flag"
    "fmt"
    "strings"
)

// Run parses command-line arguments and initiates the SSTI tests
func Run() error {
    // Command-line configurations
    urlPtr := flag.String("u", "", "Target URL")
    methodPtr := flag.String("x", "GET", "HTTP method (GET or POST)")
    dataPtr := flag.String("d", "", "Data for POST requests (use FUZZ as placeholder)")
    headersPtr := flag.String("H", "", "Additional headers in 'Header1:Value1,Header2:Value2' format")
    proxyPtr := flag.String("proxy", "", "Proxy to use for requests (e.g., http://127.0.0.1:8080)")
    payloadsFilePtr := flag.String("p", "src/payloads/payloads.json", "Path to JSON file containing the payloads")
    flag.Parse()

    if *urlPtr == "" {
        return fmt.Errorf("Please specify a URL with the -u parameter.")
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
    sstiChain, err := LoadPayloads(*payloadsFilePtr)
    if err != nil {
        return fmt.Errorf("Failed to load payloads: %v", err)
    }

    // Test the chain of SSTI payloads
    TestSSTIChain(*urlPtr, *methodPtr, *dataPtr, sstiChain, headers, *proxyPtr)

    return nil
}
