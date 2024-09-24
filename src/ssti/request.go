package ssti

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
    "time"
)

// TestSSTI performs an HTTP request with the given payload
func TestSSTI(targetURL, method, data, payload string, headers map[string]string, proxy string) (string, error) {
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Configure proxy if provided
    if proxy != "" {
        proxyURL, err := url.Parse(proxy)
        if err != nil {
            return "", fmt.Errorf("Invalid proxy: %v", err)
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
            // Append payload as a query parameter if "FUZZ" is not in URL
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
            // Append payload to data if "FUZZ" is not in data
            if data != "" {
                data += "&"
            }
            data += "payload=" + url.QueryEscape(payload)
        }
        body := strings.NewReader(data)
        req, err = http.NewRequest("POST", targetURL, body)
        req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    } else {
        return "", fmt.Errorf("Unsupported HTTP method: %s", method)
    }

    if err != nil {
        return "", fmt.Errorf("Failed to create HTTP request: %v", err)
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
        return "", fmt.Errorf("Failed to read response body: %v", err)
    }

    return string(respBody), nil
}
