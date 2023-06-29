package contracts

type CommonNetworking struct {
	Ipv6Address string `json:"ipv6_address" bson:"ipv6_address"`
	Ipv4Address string `json:"ipv4_address" bson:"ipv4_address"`
	Hostname    string `json:"hostname"`
}

type PortForward struct {
	Target uint16 `json:"target"`
	Host   uint16 `json:"host"`
}
