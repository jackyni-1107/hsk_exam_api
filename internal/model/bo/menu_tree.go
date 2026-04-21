package bo

type MenuTreeNode struct {
	Id            int64           `json:"id"`
	Name          string          `json:"name"`
	Permission    string          `json:"permission"`
	Type          int             `json:"type"`
	Sort          int             `json:"sort"`
	ParentId      int64           `json:"parent_id"`
	Path          string          `json:"path"`
	Icon          string          `json:"icon"`
	Component     string          `json:"component"`
	ComponentName string          `json:"component_name"`
	Status        int             `json:"status"`
	Visible       bool            `json:"visible"`
	KeepAlive     bool            `json:"keep_alive"`
	AlwaysShow    bool            `json:"always_show"`
	Children      []*MenuTreeNode `json:"children,omitempty"`
}
