package resource

import "time"

type AppUsage struct {
	GUID      string    `json:"guid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	App          AppUsageGUIDNamePair `json:"app"`
	Process      AppUsageGUIDTypePair `json:"process"`
	Space        AppUsageGUIDNamePair `json:"space"`
	Organization AppUsageGUID         `json:"organization"`
	Buildpack    AppUsageGUIDNamePair `json:"buildpack"`
	Task         AppUsageGUIDNamePair `json:"task"`

	State                 AppUsageCurrentPreviousString `json:"state"`
	MemoryInMbPerInstance AppUsageCurrentPreviousInt    `json:"memory_in_mb_per_instance"`
	InstanceCount         AppUsageCurrentPreviousInt    `json:"instance_count"`

	Links map[string]Link `json:"links"`
}

type AppUsageList struct {
	Pagination Pagination  `json:"pagination"`
	Resources  []*AppUsage `json:"resources"`
}

type AppUsageCurrentPreviousString struct {
	Current  string `json:"current"`
	Previous string `json:"previous"`
}

type AppUsageCurrentPreviousInt struct {
	Current  int `json:"current"`
	Previous int `json:"previous"`
}

type AppUsageGUIDNamePair struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
}

type AppUsageGUIDTypePair struct {
	GUID string `json:"guid"`
	Type string `json:"type"`
}

type AppUsageGUID struct {
	GUID string `json:"guid"`
}
