package core

import "github.com/eskpil/aarhus/pkg/contracts"

type IdentityType int64

const (
	IdentityTypeUnknown IdentityType = iota
	IdentityTypeUser
	IdentityTypeNode
)

type Identity struct {
	Type IdentityType

	User *contracts.User
	Node *contracts.Node
}
