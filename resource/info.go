package resource

// Info represents platform information about the Cloud Foundry deployment
type Info struct {
	Build         string         `json:"build"`          // The build number of the Cloud Controller API, ends up showing the full version of the api.
	CLIVersion    InfoCLIVersion `json:"cli_version"`    // Recommended and minimum CLI versions
	Custom        map[string]any `json:"custom"`         // Custom metadata set by the operator
	Description   string         `json:"description"`    // Description of the Cloud Foundry deployment
	Name          string         `json:"name"`           // Name of the Cloud Foundry deployment
	Version       int            `json:"version"`        // Version of the Cloud Controller API (Major version number)
	OSBAPIVersion string         `json:"osbapi_version"` // Version of the Open Service Broker API in use
	RateLimits    InfoRateLimits `json:"rate_limits"`    // Rate limiting configuration
	Links         Links          `json:"links"`
}

// InfoCLIVersion contains the minimum CLI version supported and the recommend version
type InfoCLIVersion struct {
	Minimum     string `json:"minimum"`     // Minimum supported CLI version
	Recommended string `json:"recommended"` // Recommended CLI version
}

// InfoRateLimits contains rate limiting configuration
type InfoRateLimits struct {
	Enabled                bool `json:"enabled"`                   // Whether rate limiting is enabled
	GeneralLimit           int  `json:"general_limit"`             // Number of requests allowed per reset interval
	ResetIntervalInMinutes int  `json:"reset_interval_in_minutes"` // Time in minutes before rate limit counter resets
}

// InfoUsageSummary represents platform-wide usage statistics
type InfoUsageSummary struct {
	UsageSummary UsageSummary `json:"usage_summary"` // Platform usage statistics
	Links        Links        `json:"links"`
}
