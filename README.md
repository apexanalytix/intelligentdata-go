# Intelligent Data API — Go SDK

Go client for the [Intelligent Data API](https://portal.smartvmapi.com) by apexanalytix. Validate addresses, tax IDs, and bank accounts; look up business registrations; screen for sanctions and disqualified directors.

## Installation

```bash
go get github.com/apexanalytix/intelligentdata-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    intelligentdata "github.com/apexanalytix/intelligentdata-go"
)

func main() {
    client := intelligentdata.NewClient("svm...")

    resp, err := client.ValidateAddress(context.Background(), intelligentdata.AddressRequest{
        AddressLine1: "123 Main St",
        City:         "New York",
        State:        "NY",
        PostalCode:   "10001",
        Country:      "US",
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Valid: %v, Score: %v\n", resp.IsValid, resp.ConfidenceScore)
}
```

## Authentication

### API Key (recommended)

```go
client := intelligentdata.NewClient("svm...")
```

### OAuth2 Client Credentials

```go
client := intelligentdata.NewClient("",
    intelligentdata.WithOAuth2("client-id", "client-secret", ""),
)
```

## Methods

| Method | Description | Endpoint |
|--------|-------------|----------|
| `ValidateAddress()` | Validate and standardize a postal address | POST /api/validate/address |
| `ValidateTaxID()` | Validate a tax identification number | POST /api/validate/taxid |
| `ValidateBankAccount()` | Verify bank account details | POST /api/validate/bank |
| `LookupBusiness()` | Look up official business registration data | POST /api/enrich/business |
| `CheckSanctions()` | Screen against global sanctions lists | POST /api/risk/sanctions |
| `CheckDirectors()` | Check for disqualified directors | POST /api/risk/directors |

## Error Handling

```go
resp, err := client.ValidateAddress(ctx, req)
if err != nil {
    var apiErr *intelligentdata.ApiError
    if errors.As(err, &apiErr) {
        if apiErr.IsRateLimit() {
            log.Println("Rate limited")
        } else if apiErr.IsAuthError() {
            log.Println("Auth failed:", apiErr.Message)
        } else {
            log.Printf("API error [%d]: %s\n", apiErr.StatusCode, apiErr.Message)
        }
    }
}
```

All response structs include a `Raw map[string]interface{}` field with the full API response.

## Options

```go
client := intelligentdata.NewClient("svm...",
    intelligentdata.WithBaseURL("https://custom-api.example.com"),
    intelligentdata.WithTimeout(15 * time.Second),
    intelligentdata.WithHTTPClient(customClient),
)
```

## Requirements

- Go 1.22+
- Zero external dependencies (stdlib only)

## MCP Protocol Support

This API is available as a native [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server.

### Tools (37), Prompts (20), Resources (7)

| Category | Count | Description |
|----------|-------|-------------|
| MCP Tools | 37 | 1:1 mapped to API endpoints — same auth, rate limits, and plan access |
| MCP Prompts | 20 | Multi-tool workflow templates (onboarding, compliance, risk, payment) |
| MCP Resources | 7 | Reference datasets (tool inventory, risk categories, country coverage) |

Prompts may be plan-gated. See [Pricing](https://portal.smartvmapi.com/pricing) for details.

- [MCP Manifest](https://portal.smartvmapi.com/.well-known/mcp.json)
- [Resource Content](https://portal.smartvmapi.com/api/portal/mcp/resources/{name})

### Accessing Resource Content

Resources provide reference data via HTTP:

```bash
curl https://portal.smartvmapi.com/api/portal/mcp/resources/tool_inventory
curl https://portal.smartvmapi.com/api/portal/mcp/resources/supported_countries
```

Available resources: `tool_inventory`, `entity_risk_categories`, `healthcare_types`, `india_identity_types`, `certification_types`, `supported_countries`, `about`.

## Getting an API Key

1. Sign up for a free account at [portal.smartvmapi.com](https://portal.smartvmapi.com/auth/register)
2. Navigate to Dashboard → API Keys
3. Copy your API key (starts with `svm`)

## License

MIT
