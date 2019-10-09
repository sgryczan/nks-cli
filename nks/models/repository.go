package models

type CreateRepoInput []CheckRepositoryInput

type RepositoryS []Repository

type Repository struct {
	ID             int          `json:"pk"`
	Name           string       `json:"name,omitempty"`
	Source         string       `json:"source"`
	SourceDisplay  string       `json:"Github"`
	OrganizationId int          `json:"org"`
	Path           string       `json:"path"`
	URL            string       `json:"url"`
	IsSystem       bool         `json:"is_system"`
	IsPrivate      bool         `json:"is_private"`
	KeysetId       *int         `json:"keyset"`
	ChartIndex     []ChartIndex `json:"chart_index"`
	State          string       `json:"state"`
	Owner          int          `json:"owner"`
	IsAccessible   bool         `json:"is_accessible"`
	Synced         *string      `json:"synced"`
	Created        string       `json:"created"`
	Updated        string       `json:"updated"`
}

type ChartIndex struct {
	Name    string                 `json:"name"`
	Sha     string                 `json:"sha"`
	Chart   Chart                  `json:"Chart"`
	Values  string                 `json:"values"`
	Path    string                 `json:"path"`
	Spec    map[string]interface{} `json:"spec"`
	Version string                 `json:"version,omitempty"`
}

type Chart struct {
	Name        string              `json:"name"`
	Created     string              `json:"created"`
	Description string              `json:"description"`
	AppVersion  string              `json:"appVersion"`
	Sources     []string            `json:"sources"`
	Maintainers []map[string]string `json:"maintainers"`
	Version     string              `json:"version"`
	URLs        []string            `json:"urls"`
	Keywords    []string            `json:"keywords"`
	Home        string              `json:"home"`
	Digest      string              `json:"digest"`
	Icon        string              `json:"icon"`
}

type CheckRepositoryInput struct {
	Name      string             `json:"name"`
	Source    string             `json:"source"`
	Path      string             `json:"path"`
	URL       string             `json:"url"`
	KeysetId  *int               `json:"keyset"`
	IsPrivate bool               `json:"is_private"`
	Config    map[string]*string `json:"config"`
}

type CheckRepositoryResponse struct {
	Accessible   bool                      `json:"accessible"`
	Directories  []string                  `json:"directories"`
	IsMultiChart bool                      `json:"is_multi_chart"`
	Error        *string                   `json:"error"`
	Contents     []CheckRepositoryContents `json:"contents"`
}

type CheckRepositoryContents struct {
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	HtmlUrl     string            `json:"html_url"`
	DownloadURL string            `json:"download_url"`
	Sha         string            `json:"sha"`
	Links       map[string]string `json:"_links"`
	GitURL      string            `json:"git_url"`
	Path        string            `json:"path"`
	Type        string            `json:"type"`
	Size        int               `json:"size"`
}
