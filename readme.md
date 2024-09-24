# SSTI Revealer
A simple tool to detect Server-Side Template Injection (SSTI) vulnerabilities.

## Features
- **Multiple Template Engine Support**: Tests for SSTI across various template engines including Jinja2, Twig, Smarty, Velocity, Freemarker, and more.
- **Customizable Payloads**: Easily modify or extend payloads via the external payloads.json file.
- **Unique Random Payloads**: Generates unique mathematical expressions using random numbers to reduce false positives and accurately detect vulnerabilities.
- **GET and POST Request Support**: Flexibility in testing different types of endpoints.
- **Proxy Support**: Ability to route requests through a proxy for inspection and debugging.
- **Custom Headers**: Add custom headers as needed for testing.

# How It Works

The tool injects unique mathematical expressions into the target application to detect SSTI vulnerabilities. By using random numbers in expressions (e.g., UNIQUE_NUM * 2), it reduces false positives and improves accuracy in identifying vulnerabilities.
Example:

- The tool generates a unique number, say 12345.
- It creates a payload like {{12345*2}}, which the template engine should evaluate to 24690.
- The tool then checks if 24690 is present in the response.
- If found, it confirms an SSTI vulnerability.

## Build

```bash
git clone git@github.com:ArthurHydr/SSTI-Revealer.git
cd SSTI-Revealer
go build .
```

## Usage

```bash
./SSTI-Revealer -u TARGER_URL [options]
```

### Options
- -u: Target URL (required).
- -x: HTTP method (GET or POST, default GET).
- -d: Data for POST requests (FUZZ as payload placeholder).
- -H: Additional headers (Header1:Value1,Header2:Value2).
- -proxy: Proxy URL (e.g., http://127.0.0.1:8080).
- -p: Path to payloads JSON file (default payloads.json).

## Example
Testing a GET request with payload injection:
```bash
./SSTI-Revealer -u "http://example.com/search?query=FUZZ"
```

Testing a POST request:
```bash
./SSTI-Revealer -u "http://example.com/login" -x POST -d "username=admin&password=FUZZ"
```

