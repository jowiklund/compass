package scope

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type Permission struct {
	Access bool `json:"access,omitempty"`
	All    bool `json:"all,omitempty"`
	Create bool `json:"create,omitempty"`
	Delete bool `json:"delete,omitempty"`
	Get    bool `json:"get,omitempty"`
	List   bool `json:"list,omitempty"`
	Modify bool `json:"modify,omitempty"`
}

type UUID = openapi_types.UUID

type TokenInformation struct {
	ClientIdentifier *string  `json:"clientIdentifier,omitempty"`
	Created          *string  `json:"created,omitempty"`
	Id               *UUID    `json:"id,omitempty"`
	Scope            *SPScope `json:"scope,omitempty"`
	TokenId          *UUID    `json:"tokenId,omitempty"`
	TokenName        *string  `json:"tokenName,omitempty"`
	ValidUntil       *string  `json:"validUntil,omitempty"`
}

type Client string

const (
	ClientSSI Client = "SynkzoneSSI"
)

type SPScope struct {
	Clients     []Client               `json:"clients"`
	Description string                 `json:"description,omitempty"`
	Extensions  map[string]interface{} `json:"extensions,omitempty"`
	Limitations map[string]interface{} `json:"limitations,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Permissions map[string]Permission  `json:"permissions,omitempty"`
}

type NewTokenRequest struct {
	Password  []string `json:"password,omitempty"`
	Scope     SPScope  `json:"scope,omitempty"`
	TokenName *string  `json:"tokenName,omitempty"`
	Validity  *int64   `json:"validity,omitempty"`
}

type ZoneMember struct {
	AccessLevel   string       `json:"accessLevel,omitempty"`
	Authorization string `json:"authorization,omitempty"`
	DisplayName   string                       `json:"displayName,omitempty"`
	MemberId      string                       `json:"memberId,omitempty"`
	MemberType    string            `json:"memberType,omitempty"`
}

type ZoneData struct {
	AllowExternalSharing *bool       `json:"allowExternalSharing,omitempty"`
	CreatedAt            *string     `json:"createdAt,omitempty"`
	CreatedBy            *string     `json:"createdBy,omitempty"`
	CurrentAccess        *ZoneMember `json:"currentAccess,omitempty"`
	Name                 string      `json:"name,omitempty"`
	Description          string      `json:"description,omitempty"`
	Id                   UUID `json:"id,omitempty"`
	SsiAccessible        bool        `json:"ssiAccessible,omitempty"`
	Type                 UUID `json:"type,omitempty"`
	TypeName             string      `json:"typeName,omitempty"`
	WebAccessible        bool        `json:"webAccessible,omitempty"`
}
