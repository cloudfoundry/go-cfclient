package resource

type Sidecar struct {
	// Human-readable name for the sidecar
	Name string `json:"name"`

	// The command used to start the sidecar
	Command string `json:"command"`

	// A list of process types the sidecar applies to
	ProcessTypes []string `json:"process_types"`

	// Reserved memory for sidecar in MB
	MemoryInMB int `json:"memory_in_mb"`

	// Specifies whether the sidecar was created by the user or via the buildpack
	Origin string `json:"origin"`

	// The app the sidecar is associated with
	Relationships AppRelationship `json:"relationships"`

	Resource `json:",inline"`
}

type SidecarList struct {
	Pagination Pagination `json:"pagination"`
	Resources  []*Sidecar `json:"resources"`
}

type SidecarCreate struct {
	// Human-readable name for the sidecar
	Name string `json:"name"`

	// The command used to start the sidecar
	Command string `json:"command"`

	// A list of process types the sidecar applies to
	ProcessTypes []string `json:"process_types"`

	// Reserved memory for sidecar in MB
	MemoryInMB *int `json:"memory_in_mb,omitempty"`
}

type SidecarUpdate struct {
	// Human-readable name for the sidecar
	Name *string `json:"name,omitempty"`

	// The command used to start the sidecar
	Command *string `json:"command,omitempty"`

	// A list of process types the sidecar applies to
	ProcessTypes []string `json:"process_types,omitempty"`

	// Reserved memory for sidecar in MB
	MemoryInMB *int `json:"memory_in_mb,omitempty"`
}

func NewSidecarCreate(name, command string, processTypes []string) *SidecarCreate {
	return &SidecarCreate{
		Name:         name,
		Command:      command,
		ProcessTypes: processTypes,
	}
}

func (s *SidecarCreate) WithMemoryInMB(mb int) *SidecarCreate {
	s.MemoryInMB = &mb
	return s
}

func NewSidecarUpdate() *SidecarUpdate {
	return &SidecarUpdate{}
}

func (s *SidecarUpdate) WithMemoryInMB(mb int) *SidecarUpdate {
	s.MemoryInMB = &mb
	return s
}

func (s *SidecarUpdate) WithName(name string) *SidecarUpdate {
	s.Name = &name
	return s
}

func (s *SidecarUpdate) WithCommand(command string) *SidecarUpdate {
	s.Command = &command
	return s
}

func (s *SidecarUpdate) WithProcessTypes(processTypes []string) *SidecarUpdate {
	s.ProcessTypes = processTypes
	return s
}
