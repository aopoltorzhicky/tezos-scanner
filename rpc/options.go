package scanner

// NodeOption -
type NodeOption func(*Node)

// WithRPCURL -
func WithRPCURL(url string) NodeOption {
	return func(node *Node) {
		node.RPCURI = url
	}
}

// WithListenerURL -
func WithListenerURL(url string) NodeOption {
	return func(node *Node) {
		node.ListenerURI = url
	}
}
