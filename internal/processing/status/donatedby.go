/*
   GoToSocial
   Copyright (C) 2021-2023 GoToSocial Authors admin@gotosocial.org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package status

import (
	"context"
	"errors"
	"fmt"

	apimodel "github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/util"
)

func (p *processor) DonatedBy(ctx context.Context, requestingAccount *gtsmodel.Account, targetStatusID string) ([]*apimodel.StatusDonatedBy, gtserror.WithCode) {
	targetStatus, err := p.db.GetStatusByID(ctx, targetStatusID)
	if err != nil {
		return nil, gtserror.NewErrorNotFound(fmt.Errorf("error fetching status %s: %s", targetStatusID, err))
	}
	if targetStatus.Account == nil {
		return nil, gtserror.NewErrorNotFound(fmt.Errorf("no status owner for status %s", targetStatusID))
	}

	visible, err := p.filter.StatusVisible(ctx, targetStatus, requestingAccount)
	if err != nil {
		return nil, gtserror.NewErrorNotFound(fmt.Errorf("error seeing if status %s is visible: %s", targetStatus.ID, err))
	}
	if !visible {
		return nil, gtserror.NewErrorNotFound(errors.New("status is not visible"))
	}

	statusDonates, err := p.db.GetStatusDonates(ctx, targetStatus)
	if err != nil {
		return nil, gtserror.NewErrorNotFound(fmt.Errorf("error seeing who donated status: %s", err))
	}

	// filter the list so the user doesn't see accounts they blocked or which blocked them
	apiStatusDontedByAccounts := []*apimodel.StatusDonatedBy{}
	for _, donate := range statusDonates {
		blocked, err := p.db.IsBlocked(ctx, requestingAccount.ID, donate.AccountID, true)
		if err != nil {
			return nil, gtserror.NewErrorInternalError(fmt.Errorf("error checking blocks: %s", err))
		}
		if blocked {
			continue
		}

		apiAccount, err := p.tc.AccountToAPIAccountPublic(ctx, donate.Account)
		if err != nil {
			return nil, gtserror.NewErrorInternalError(fmt.Errorf("error converting status %s to frontend representation: %s", targetStatus.ID, err))
		}
		apiStatusDontedByAccounts = append(apiStatusDontedByAccounts, &apimodel.StatusDonatedBy{
			ID:        donate.ID,
			Account:   apiAccount,
			Message:   donate.Message,
			Amount:    donate.Amount,
			CreatedAt: util.FormatISO8601(donate.CreatedAt),
		})
	}

	return apiStatusDontedByAccounts, nil
}
