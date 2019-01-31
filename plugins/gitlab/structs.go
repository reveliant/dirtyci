package main

type PushEvent struct {
	ObjectKind          string      `json:"object_kind"`
	Before              string      `json:"before"`
	After               string      `json:"after"`
	Ref                 string      `json:"ref"`
	CheckoutSha         string      `json:"checkout_sha"`
	UserId              int         `json:"user_id"`
	UserName            string      `json:"user_name"`
	UserUsername        string      `json:"user_username"`
	UserEmail           string      `json:"user_email"`
	UserAvatar          string      `json:"user_avatar"`
	ProjectId           int         `json:"project_id"`
	Project             Project     `json:"project"`
	Repository          Repository  `json:"repository"`
	Commits             []Commit    `json:"commits"`
	TotalCommitsCount   int         `json:"total_commits_count"`
}

type Project struct {
	Id                  int         `json:"id"`
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	WebUrl              string      `json:"web_url"`
	AvatarUrl           string      `json:"avatar_url"`
	GitSshUrl           string      `json:"git_ssh_url"`
	GitHttpUrl          string      `json:"git_http_url"`
	Namespace           string      `json:"namespace"`
	VisibilityLevel     int         `json:"visibility_level"`
	PathWithNamespace   string      `json:"path_with_namespace"`
	DefaultBranch       string      `json:"default_branch"`
	Homepage            string      `json:"homepage"`
	Url                 string      `json:"url"`
	SshUrl              string      `json:"ssh_url"`
	HttpUrl             string      `json:"http_url"`
}

type Repository struct {
	Name                string      `json:"name"`
	Url                 string      `json:"url"`
	Description         string      `json:"description"`
	Homepage            string      `json:"homepage"`
	GitHttpUrl          string      `json:"git_http_url"`
	GitSshUrl           string      `json:"git_ssh_url"`
	VisibilityLevel     int         `json:"visibility_level"`
}

type Commit struct {
	Id                  string      `json:"id"`
	Message             string      `json:"message"`
	Timestamp           string      `json:"timestamp"`
	Url                 string      `json:"url"`
	Author              Author      `json:"author"`
	Added               []string    `json:"added"`
	Modified            []string    `json:"modified"`
	Removed             []string    `json:"removed"`
}

type Author struct {
	Name                string      `json:"name"`
	Email               string      `json:"email"`
}
