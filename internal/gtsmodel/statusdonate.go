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

package gtsmodel

import "time"

// StatusDonate refers to a 'donate' in the database, from one account, targeting the status of another account
type StatusDonate struct {
	ID              string    `validate:"required,ulid" bun:"type:CHAR(26),pk,nullzero,notnull,unique"`        // id of this item in the database
	CreatedAt       time.Time `validate:"-" bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"` // when was item created
	UpdatedAt       time.Time `validate:"-" bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"` // when was item last updated
	AccountID       string    `validate:"required,ulid" bun:"type:CHAR(26),nullzero,notnull"`                  // id of the account that created ('did') the donate
	Account         *Account  `validate:"-" bun:"rel:belongs-to"`                                              // account that created the donate
	TargetAccountID string    `validate:"required,ulid" bun:"type:CHAR(26),nullzero,notnull"`                  // id the account owning the donated status
	TargetAccount   *Account  `validate:"-" bun:"rel:belongs-to"`                                              // account owning the donated status
	StatusID        string    `validate:"required,ulid" bun:"type:CHAR(26),nullzero,notnull"`                  // database id of the status that has been 'donated'
	Status          *Status   `validate:"-" bun:"rel:belongs-to"`                                              // the donated status
	Message         string    `validate:"-" bun:",nullzero"`                                                   // message for this donate, by the donater
	Amount          int64     `validate:"-" bun:",notnull,default:0"`                                          // amount of this donate
	WalletDonateID  string    `validate:"-" bun:"type:CHAR(26),nullzero,notnull"`                              // txid of this donate
	URI             string    `validate:"required,url" bun:",nullzero,notnull"`                                // ActivityPub URI of this donate
}
