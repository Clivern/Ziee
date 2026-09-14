// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package v1

// File is a parsed `.ziee.yml`.
type File struct {
	Version     string      `yaml:"version"`
	Labels      []Label     `yaml:"labels"`
	Teams       []Team      `yaml:"teams"`
	Knowledge   []Knowledge `yaml:"knowledge"`
	MergeQueue  MergeQueue  `yaml:"merge_queue"`
	PRReviews   PRReviews   `yaml:"pr_reviews"`
	IssueTriage IssueTriage `yaml:"issue_triage"`
}

// Label is a GitHub label defined in `.ziee.yml`.
type Label struct {
	Name        string `yaml:"name"`
	Color       string `yaml:"color"`
	Description string `yaml:"description"`
}

// Team is a named group of GitHub logins in `.ziee.yml`.
type Team struct {
	Name    string   `yaml:"name"`
	Members []string `yaml:"members"`
}

// Knowledge is a repository path indexed into the workspace knowledge base.
type Knowledge struct {
	Path string       `yaml:"path"`
	Tags KnowledgeTag `yaml:"tags"`
}

// MergeQueue is the pull-request automation block.
type MergeQueue struct {
	Enabled              bool           `yaml:"enabled"`
	Mode                 string         `yaml:"mode"`
	MaxParallelChecks    int            `yaml:"max_parallel_checks"`
	ResetOnExternalMerge string         `yaml:"reset_on_external_merge"`
	Labels               QueueLabels    `yaml:"labels"`
	Comments             string         `yaml:"comments"`
	PRTriage             PRTriage       `yaml:"pr_triage"`
	Commands             Commands       `yaml:"commands"`
	PriorityRules        []PriorityRule `yaml:"priority_rules"`
	QueueRules           []QueueRule    `yaml:"queue_rules"`
}

// QueueLabels are GitHub labels Ziee applies for queue state.
type QueueLabels struct {
	Queued   string `yaml:"queued"`
	Checking string `yaml:"checking"`
	Dequeued string `yaml:"dequeued"`
}

// PRTriage labels, assigns, and requests reviewers on pull requests.
type PRTriage struct {
	AI    AI     `yaml:"ai"`
	Rules []Rule `yaml:"rules"`
}

// IssueTriage labels, assigns, and comments on GitHub issues.
type IssueTriage struct {
	Enabled  bool     `yaml:"enabled"`
	Comments string   `yaml:"comments"`
	AI       AI       `yaml:"ai"`
	Commands Commands `yaml:"commands"`
	Rules    []Rule   `yaml:"rules"`
}

// PRReviews is parsed and unused until the review engine exists.
type PRReviews struct {
	Enabled bool `yaml:"enabled"`
}

// AI classifies intention from title, body, and (for PRs) diff.
type AI struct {
	Enabled   bool           `yaml:"enabled"`
	Knowledge []KnowledgeTag `yaml:"knowledge"`
}

// KnowledgeTag is a set of workspace knowledge label key-value pairs.
type KnowledgeTag map[string]string

// Rule is one triage rule. Every matching rule fires.
type Rule struct {
	Name        string   `yaml:"name"`
	When        Clauses  `yaml:"when"`
	Labels      Labels   `yaml:"labels,omitempty"`
	Assign      []string `yaml:"assign,omitempty"`
	Reviewers   []string `yaml:"reviewers,omitempty"`
	ReviewTeams []string `yaml:"review_teams,omitempty"`
	Comment     string   `yaml:"comment,omitempty"`
}

// Labels are GitHub labels to add or remove.
type Labels struct {
	Add    []string `yaml:"add,omitempty"`
	Remove []string `yaml:"remove,omitempty"`
}

// Commands maps `@ziee` verbs to allow lists.
type Commands map[string]Command

// Command is who may run one `@ziee` verb.
type Command struct {
	Allow Allow `yaml:"allow,omitempty"`
}

// PriorityRule is parsed; the queue engine is not run in this delivery.
type PriorityRule struct {
	Name            string  `yaml:"name"`
	When            Clauses `yaml:"when"`
	Priority        string  `yaml:"priority"`
	InterruptChecks bool    `yaml:"interrupt_checks,omitempty"`
}

// QueueRule is parsed; the queue engine is not run in this delivery.
type QueueRule struct {
	Name             string    `yaml:"name"`
	Allow            Allow     `yaml:"allow,omitempty"`
	QueueWhen        Clauses   `yaml:"queue_when,omitempty"`
	MergeWhen        Clauses   `yaml:"merge_when,omitempty"`
	BatchSize        BatchSize `yaml:"batch_size,omitempty"`
	BatchMaxWait     string    `yaml:"batch_max_wait,omitempty"`
	MergeMethod      string    `yaml:"merge_method,omitempty"`
	ChecksTimeout    string    `yaml:"checks_timeout,omitempty"`
	MaxChecksRetries int       `yaml:"max_checks_retries,omitempty"`
}

// Clauses is a list of `when` matchers. All must match.
type Clauses []Clause

// Clause is one matcher in a `when` list.
type Clause struct {
	Files             []string   `json:"files,omitempty" yaml:"files,omitempty"`
	MaxFilesChanged   *int       `json:"max_files_changed,omitempty" yaml:"max_files_changed,omitempty"`
	MinFilesChanged   *int       `json:"min_files_changed,omitempty" yaml:"min_files_changed,omitempty"`
	Title             string     `json:"title,omitempty" yaml:"title,omitempty"`
	Body              string     `json:"body,omitempty" yaml:"body,omitempty"`
	AuthorIn          []string   `json:"author_in,omitempty" yaml:"author_in,omitempty"`
	AuthorNotIn       []string   `json:"author_not_in,omitempty" yaml:"author_not_in,omitempty"`
	AuthorInTeam      []string   `json:"author_in_team,omitempty" yaml:"author_in_team,omitempty"`
	AuthorNotInTeam   []string   `json:"author_not_in_team,omitempty" yaml:"author_not_in_team,omitempty"`
	FirstContribution *bool      `json:"first_contribution,omitempty" yaml:"first_contribution,omitempty"`
	Intention         Intention  `json:"intention,omitempty" yaml:"intention,omitempty"`
	Label             string     `json:"label,omitempty" yaml:"label,omitempty"`
	Check             string     `json:"check,omitempty" yaml:"check,omitempty"`
	Approvals         *Approvals `json:"approvals,omitempty" yaml:"approvals,omitempty"`
}

// Intention is an AI classify label. A string is the name only.
type Intention struct {
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// Approvals is a `queue_when` / `merge_when` matcher.
type Approvals struct {
	Min  int      `yaml:"min,omitempty" json:"min,omitempty"`
	From []string `yaml:"from,omitempty" json:"from,omitempty"`
}

// Allow is a list of permission, team, or user entries. Any match is enough.
type Allow []AllowEntry

// AllowEntry is one ACL item.
type AllowEntry struct {
	Permission string   `json:"permission,omitempty" yaml:"permission,omitempty"`
	Teams      []string `json:"teams,omitempty" yaml:"teams,omitempty"`
	Users      []string `json:"users,omitempty" yaml:"users,omitempty"`
	Self       bool     `json:"self,omitempty" yaml:"self,omitempty"`
}

// BatchSize is either a single int or {min, max}.
type BatchSize struct {
	Value *int `json:"value,omitempty" yaml:"value,omitempty"`
	Min   int  `json:"min,omitempty" yaml:"min,omitempty"`
	Max   int  `json:"max,omitempty" yaml:"max,omitempty"`
}
