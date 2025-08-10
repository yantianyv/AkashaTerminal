package types

// FileState 表示文件的当前状态
type FileState struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Checksum string `json:"checksum"` // SHA256 校验和
}

// ProjectState 表示项目的整体状态
type ProjectState struct {
	CWD          string              `json:"cwd"`
	Directory    string              `json:"directory_tree"`
	FileStates   map[string]FileState `json:"file_states"`
	ScannedDirs  map[string]string    `json:"scanned_dirs"`
	MaxDepth     int                 `json:"max_depth"`
	Truncated    bool                `json:"truncated"`
}

// FileOperation 定义文件操作
type FileOperation struct {
	Action   string `json:"action"`  // read/write/create/scan
	Path     string `json:"path"`
	Content  string `json:"content,omitempty"`
	Mode     string `json:"mode,omitempty"`   // replace/insert/append
	OldText  string `json:"old_text,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}

// APIConfig 表示 API 配置
type APIConfig struct {
	Name       string  `json:"-"`
	Provider   string  `json:"provider"`
	APIKey     string  `json:"api_key"`
	APIBase    string  `json:"api_base"`
	Model      string  `json:"model"`
	Deployment string  `json:"deployment,omitempty"`
	MaxTokens  int     `json:"max_tokens,omitempty"`
	Version    string  `json:"version,omitempty"`
	AppID      string  `json:"app_id,omitempty"`
	AgentID    string  `json:"agent_id,omitempty"`
	AuthType   string  `json:"auth_type,omitempty"`
	AuthKey    string  `json:"auth_key,omitempty"`
}

// AIProvider 是所有供应商实现的接口
type AIProvider interface {
	SendRequest(prompt string, state interface{}) (string, error)
	GetName() string
	GetModel() string
	SupportsFeature(feature string) bool
}
