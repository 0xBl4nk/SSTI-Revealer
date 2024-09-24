package ssti

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "strings"
    "time"
)

// LoadPayloads loads SSTI payloads from a JSON file
func LoadPayloads(filePath string) (map[string]string, error) {
    payloadsData, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("Failed to read payloads file: %v", err)
    }

    var sstiChain map[string]string
    err = json.Unmarshal(payloadsData, &sstiChain)
    if err != nil {
        return nil, fmt.Errorf("Failed to parse payloads file: %v", err)
    }

    return sstiChain, nil
}

// TestSSTIChain tests a series of SSTI payloads against the target
func TestSSTIChain(url, method, data string, sstiChain map[string]string, headers map[string]string, proxy string) {
    vulnerable := false

    for engine, payloadTemplate := range sstiChain {
        uniqueNum := (time.Now().UnixNano() % 10000) + 10000 // Generates a 5-digit number
        expectedResult := fmt.Sprintf("%d", uniqueNum*2)

        // Replace placeholders in the payload
        payload := strings.Replace(payloadTemplate, "UNIQUE_NUM", fmt.Sprintf("%d", uniqueNum), -1)

        resp, err := TestSSTI(url, method, data, payload, headers, proxy)
        if err != nil {
            fmt.Printf("[-] Error testing %s: %v\n", engine, err)
            continue
        }

        if CheckVulnerability(resp, expectedResult, engine, payload) {
            vulnerable = true
        } else {
            fmt.Printf("[-] Not vulnerable to %s SSTI.\n", engine)
        }
    }

    if !vulnerable {
        fmt.Println("[-] No SSTI vulnerabilities detected.")
    }
}

// CheckVulnerability checks if the response contains the expected result
func CheckVulnerability(resp, expectedResult, engine, payload string) bool {
    cleanResp := strings.ReplaceAll(resp, " ", "")
    if strings.Contains(cleanResp, expectedResult) {
        fmt.Printf("[+] Vulnerable to %s SSTI with payload: %s\n", engine, payload)
        return true
    }
    return false
}
