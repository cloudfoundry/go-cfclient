package resource

// Stack implements stack object. Stacks are the base operating system and file system that your
// application will execute in. A stack is how you configure applications to run against different
// operating systems (like Windows or Linux) and different versions of those operating systems.
type Stack struct {
	Name             string    `json:"name"`
	Description      *string   `json:"description"`
	RunRootfsImage   string    `json:"run_rootfs_image"`
	BuildRootfsImage string    `json:"build_rootfs_image"`
	Default          bool      `json:"default"`
	Metadata         *Metadata `json:"metadata"`
	Resource         `json:",inline"`
}

type StackCreate struct {
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Metadata    *Metadata `json:"metadata,omitempty"`
}

type StackUpdate struct {
	Metadata *Metadata `json:"metadata,omitempty"`
}

type StackList struct {
	Pagination Pagination `json:"pagination"`
	Resources  []*Stack   `json:"resources"`
}
