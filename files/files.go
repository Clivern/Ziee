// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package files

import _ "embed"

//go:embed .ziee.yaml
var DefaultZieeYML string

//go:embed setup_issue.md
var SetupIssueBody string

//go:embed setup_pr.md
var SetupPullRequestBody string
