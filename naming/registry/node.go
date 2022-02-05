package registry

// Node 服务节点信息
type Node struct {
	ServiceName string
	Address     string
	Network     string
	Protocol    string
	Weight      int
	Metadata    map[string]interface{}
}
