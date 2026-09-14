// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"net/http"
)

// ListUserTeams lists organization teams that include login.
func (a *App) ListUserTeams(ctx context.Context, installationID int64, org, login string) ([]Team, error) {
	token, err := a.GetInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var all []Team
	for page := 1; ; page++ {
		var teams []Team
		path := fmt.Sprintf("%s/orgs/%s/members/%s/teams?per_page=%d&page=%d", a.apiURL, org, login, AppPerPage, page)
		err = Call(ctx, http.MethodGet, path, token.Token, GetHeaders(), nil, &teams)
		if err != nil {
			return nil, err
		}

		all = append(all, teams...)
		if len(teams) < AppPerPage {
			return all, nil
		}
	}
}
