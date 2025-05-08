package global_types

type IResponseAPI struct {
	Message      string `json:"message,omitempty"`
	ErrorSection string `json:"error_section,omitempty"`
	Data         any    `json:"data,omitempty"`
}
