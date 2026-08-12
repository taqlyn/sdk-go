package taqlyn

// LinkMode controls how a short link routes visitors.
type LinkMode string

const (
	ModeWebOnly     LinkMode = "web_only"
	ModeAppAware    LinkMode = "app_aware"
	ModeDeferredApp LinkMode = "deferred_app"
)

// Environment identifies a Taqlyn credential environment.
type Environment string

const (
	EnvironmentSandbox    Environment = "sandbox"
	EnvironmentProduction Environment = "production"
)

// CreateShortLinkRequest is the payload accepted by CreateShortLink.
type CreateShortLinkRequest struct {
	DestinationWeb  string         `json:"destinationWeb"`
	Mode            LinkMode       `json:"mode,omitempty"`
	DestinationPath string         `json:"destinationPath,omitempty"`
	Params          map[string]any `json:"params,omitempty"`
	Environment     Environment    `json:"env,omitempty"`
	OGTitle         string         `json:"ogTitle,omitempty"`
	OGDescription   string         `json:"ogDescription,omitempty"`
	OGImage         string         `json:"ogImage,omitempty"`
}

// PatchShortLinkRequest contains fields that may be updated on a short link.
// Pointer fields distinguish an omitted value from a value being cleared.
type PatchShortLinkRequest struct {
	DestinationWeb  *string         `json:"destinationWeb,omitempty"`
	DestinationPath *string         `json:"destinationPath,omitempty"`
	Params          *map[string]any `json:"params,omitempty"`
	Mode            *LinkMode       `json:"mode,omitempty"`
	OGTitle         *string         `json:"ogTitle,omitempty"`
	OGDescription   *string         `json:"ogDescription,omitempty"`
	OGImage         *string         `json:"ogImage,omitempty"`
}

// ShortLink is returned by the Taqlyn short-link API.
type ShortLink struct {
	ID             string      `json:"id"`
	Code           string      `json:"code"`
	ShortURL       string      `json:"shortUrl"`
	Host           string      `json:"host"`
	Mode           LinkMode    `json:"mode"`
	DestinationWeb string      `json:"destinationWeb"`
	Environment    Environment `json:"env"`
	OrgID          string      `json:"orgId"`
	AppID          string      `json:"appId"`
	OGTitle        string      `json:"ogTitle,omitempty"`
	OGDescription  string      `json:"ogDescription,omitempty"`
	OGImage        string      `json:"ogImage,omitempty"`
}
