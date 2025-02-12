package resource

type AppsQuota struct {
	// The collective memory allocation permitted for all initiated processes and active tasks.
	TotalMemoryInMB *int `json:"total_memory_in_mb"`

	// The upper limit for memory allocation per individual process or task.
	PerProcessMemoryInMB *int `json:"per_process_memory_in_mb"`

	// The cumulative log rate limit permitted for all initiated processes and active tasks.
	LogRateLimitInBytesPerSecond *int `json:"log_rate_limit_in_bytes_per_second"`

	// The maximum allowable total instances of all started processes.
	TotalInstances *int `json:"total_instances"`

	// The maximum limit for the number of tasks currently running.
	PerAppTasks *int `json:"per_app_tasks"`
}

type ServicesQuota struct {
	// Specifies if instances of paid service plans are permitted to be created.
	PaidServicesAllowed bool `json:"paid_services_allowed"`

	// The maximum number of service instances permitted.
	TotalServiceInstances *int `json:"total_service_instances"`

	// The maximum number of service keys permitted within an organization.
	TotalServiceKeys *int `json:"total_service_keys"`
}

type RoutesQuota struct {
	// The maximum number of routes permitted within an organization.
	TotalRoutes *int `json:"total_routes"`

	// The maximum number of ports that can be reserved by routes within an organization.
	TotalReservedPorts *int `json:"total_reserved_ports"`
}

type DomainsQuota struct {
	// The maximum number of domains that can be associated with an organization.
	TotalDomains *int `json:"total_domains"`
}
