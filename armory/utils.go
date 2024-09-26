package armory

import (
	"fmt"
	"net/http"
	"time"

	"github.com/privateerproj/privateer-sdk/raidengine"
)

// MakeGETRequest makes a GET request to the specified endpoint and returns the status code
func MakeGETRequest(endpoint string, result *raidengine.MovementResult) *http.Response {
	result.Description = fmt.Sprintf("Making GET request to endpoint: %s", endpoint)

	// Create an HTTP client with a timeout for safety
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the GET request
	response, err := client.Get(endpoint)
	if err != nil {
		result.Passed = false
		result.Message = err.Error()
		return response
	}
	defer response.Body.Close()

	result.Description = fmt.Sprintf("Response contained HTTP status code: %d", response.StatusCode)

	// Check for HTTP success (2xx status codes)
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		result.Passed = true
	} else {
		result.Passed = false
	}
	return response
}

// CheckStatusCode checks the TLS version of the response and updates the result
func CheckTLSVersion(response *http.Response, result *raidengine.MovementResult) {
	result.Description = "Checking TLS version of response"

	// Check the TLS version of the response
	tlsVersion := response.TLS.Version
	if tlsVersion == 0 {
		result.Passed = false
		result.Message = "No TLS version found in response"
	} else {
		result.Passed = true
		result.Message = fmt.Sprintf("TLS version: %v", tlsVersion)
	}

	// Check if the connection used TLS
	if response.TLS != nil {
		tlsVersion := response.TLS.Version
		// Map TLS version to human-readable format
		switch tlsVersion {
		case 0x0304:
			result.Message = "TLS 1.3 is being used"
			result.Passed = true
		case 0x0303:
			result.Message = "TLS 1.2 is being used"
			result.Passed = true
		case 0x0302:
			result.Message = "TLS 1.1 is being used"
			result.Passed = false
		case 0x0301:
			result.Message = "TLS 1.0 is being used"
			result.Passed = false
		default:
			result.Message = "error: Unknown TLS version"
			result.Passed = false
		}
	} else {
		result.Message = "error: No TLS information found in response"
		result.Passed = false
	}
}
