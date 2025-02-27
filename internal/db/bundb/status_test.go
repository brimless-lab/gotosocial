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

package bundb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/superseriousbusiness/gotosocial/internal/db"
)

type StatusTestSuite struct {
	BunDBStandardTestSuite
}

func (suite *StatusTestSuite) TestGetStatusByID() {
	status, err := suite.db.GetStatusByID(context.Background(), suite.testStatuses["local_account_1_status_1"].ID)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.NotNil(status)
	suite.NotNil(status.Account)
	suite.NotNil(status.CreatedWithApplication)
	suite.Nil(status.BoostOf)
	suite.Nil(status.BoostOfAccount)
	suite.Nil(status.InReplyTo)
	suite.Nil(status.InReplyToAccount)
	suite.True(*status.Federated)
	suite.True(*status.Boostable)
	suite.True(*status.Replyable)
	suite.True(*status.Likeable)
}

func (suite *StatusTestSuite) TestGetStatusesByID() {
	ids := []string{
		suite.testStatuses["local_account_1_status_1"].ID,
		suite.testStatuses["local_account_2_status_3"].ID,
	}

	statuses, err := suite.db.GetStatuses(context.Background(), ids)
	if err != nil {
		suite.FailNow(err.Error())
	}

	if len(statuses) != 2 {
		suite.FailNow("expected 2 statuses in slice")
	}

	status1 := statuses[0]
	suite.NotNil(status1)
	suite.NotNil(status1.Account)
	suite.NotNil(status1.CreatedWithApplication)
	suite.Nil(status1.BoostOf)
	suite.Nil(status1.BoostOfAccount)
	suite.Nil(status1.InReplyTo)
	suite.Nil(status1.InReplyToAccount)
	suite.True(*status1.Federated)
	suite.True(*status1.Boostable)
	suite.True(*status1.Replyable)
	suite.True(*status1.Likeable)

	status2 := statuses[1]
	suite.NotNil(status2)
	suite.NotNil(status2.Account)
	suite.NotNil(status2.CreatedWithApplication)
	suite.Nil(status2.BoostOf)
	suite.Nil(status2.BoostOfAccount)
	suite.Nil(status2.InReplyTo)
	suite.Nil(status2.InReplyToAccount)
	suite.True(*status2.Federated)
	suite.True(*status2.Boostable)
	suite.False(*status2.Replyable)
	suite.False(*status2.Likeable)
}

func (suite *StatusTestSuite) TestGetStatusByURI() {
	status, err := suite.db.GetStatusByURI(context.Background(), suite.testStatuses["local_account_2_status_3"].URI)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.NotNil(status)
	suite.NotNil(status.Account)
	suite.NotNil(status.CreatedWithApplication)
	suite.Nil(status.BoostOf)
	suite.Nil(status.BoostOfAccount)
	suite.Nil(status.InReplyTo)
	suite.Nil(status.InReplyToAccount)
	suite.True(*status.Federated)
	suite.True(*status.Boostable)
	suite.False(*status.Replyable)
	suite.False(*status.Likeable)
}

func (suite *StatusTestSuite) TestGetStatusWithExtras() {
	status, err := suite.db.GetStatusByID(context.Background(), suite.testStatuses["admin_account_status_1"].ID)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.NotNil(status)
	suite.NotNil(status.Account)
	suite.NotNil(status.CreatedWithApplication)
	suite.NotEmpty(status.Tags)
	suite.NotEmpty(status.Attachments)
	suite.NotEmpty(status.Emojis)
	suite.True(*status.Federated)
	suite.True(*status.Boostable)
	suite.True(*status.Replyable)
	suite.True(*status.Likeable)
}

func (suite *StatusTestSuite) TestGetStatusWithMention() {
	status, err := suite.db.GetStatusByID(context.Background(), suite.testStatuses["local_account_2_status_5"].ID)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.NotNil(status)
	suite.NotNil(status.Account)
	suite.NotNil(status.CreatedWithApplication)
	suite.NotEmpty(status.MentionIDs)
	suite.NotEmpty(status.InReplyToID)
	suite.NotEmpty(status.InReplyToAccountID)
	suite.True(*status.Federated)
	suite.True(*status.Boostable)
	suite.True(*status.Replyable)
	suite.True(*status.Likeable)
}

func (suite *StatusTestSuite) TestGetStatusTwice() {
	before1 := time.Now()
	_, err := suite.db.GetStatusByURI(context.Background(), suite.testStatuses["local_account_1_status_1"].URI)
	suite.NoError(err)
	after1 := time.Now()
	duration1 := after1.Sub(before1)
	fmt.Println(duration1.Microseconds())

	before2 := time.Now()
	_, err = suite.db.GetStatusByURI(context.Background(), suite.testStatuses["local_account_1_status_1"].URI)
	suite.NoError(err)
	after2 := time.Now()
	duration2 := after2.Sub(before2)
	fmt.Println(duration2.Microseconds())

	// second retrieval should be several orders faster since it will be cached now
	suite.Less(duration2, duration1)
}

func (suite *StatusTestSuite) TestGetStatusChildren() {
	targetStatus := suite.testStatuses["local_account_1_status_1"]
	children, err := suite.db.GetStatusChildren(context.Background(), targetStatus, true, "")
	suite.NoError(err)
	suite.Len(children, 2)
	for _, c := range children {
		suite.Equal(targetStatus.URI, c.InReplyToURI)
		suite.Equal(targetStatus.AccountID, c.InReplyToAccountID)
		suite.Equal(targetStatus.ID, c.InReplyToID)
	}
}

func (suite *StatusTestSuite) TestDeleteStatus() {
	targetStatus := suite.testStatuses["admin_account_status_1"]
	err := suite.db.DeleteStatusByID(context.Background(), targetStatus.ID)
	suite.NoError(err)

	_, err = suite.db.GetStatusByID(context.Background(), targetStatus.ID)
	suite.ErrorIs(err, db.ErrNoEntries)
}

func TestStatusTestSuite(t *testing.T) {
	suite.Run(t, new(StatusTestSuite))
}
