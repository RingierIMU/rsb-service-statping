// Statping
// Copyright (C) 2018.  Hunter Long and the project contributors
// Written by Hunter Long <info@socialeck.com> and the project contributors
//
// https://github.com/hunterlong/statping
//
// The licenses for most software and other practical works are designed
// to take away your freedom to share and change the works.  By contrast,
// the GNU General Public License is intended to guarantee your freedom to
// share and change all versions of a program--to make sure it remains free
// software for all its users.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package notifiers

import (
	"fmt"
	"github.com/statping/statping/database"
	"github.com/statping/statping/types/notifications"
	"github.com/statping/statping/types/null"
	"github.com/statping/statping/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var (
	EMAIL_HOST     string
	EMAIL_USER     string
	EMAIL_PASS     string
	EMAIL_OUTGOING string
	EMAIL_SEND_TO  string
	EMAIL_PORT     int64
)

var testEmail *emailOutgoing

func init() {
	EMAIL_HOST = os.Getenv("EMAIL_HOST")
	EMAIL_USER = os.Getenv("EMAIL_USER")
	EMAIL_PASS = os.Getenv("EMAIL_PASS")
	EMAIL_OUTGOING = os.Getenv("EMAIL_OUTGOING")
	EMAIL_SEND_TO = os.Getenv("EMAIL_SEND_TO")
	EMAIL_PORT = utils.ToInt(os.Getenv("EMAIL_PORT"))
}

func TestEmailNotifier(t *testing.T) {
	db, err := database.OpenTester()
	require.Nil(t, err)
	db.AutoMigrate(&notifications.Notification{})
	notifications.SetDB(db)

	if EMAIL_HOST == "" || EMAIL_USER == "" || EMAIL_PASS == "" {
		t.Log("email notifier testing skipped, missing EMAIL_ environment variables")
		t.SkipNow()
	}

	t.Run("New email", func(t *testing.T) {
		email.Host = EMAIL_HOST
		email.Username = EMAIL_USER
		email.Password = EMAIL_PASS
		email.Var1 = EMAIL_OUTGOING
		email.Var2 = EMAIL_SEND_TO
		email.Port = int(EMAIL_PORT)
		email.Delay = time.Duration(100 * time.Millisecond)
		email.Enabled = null.NewNullBool(true)

		Add(email)
		assert.Equal(t, "Hunter Long", email.Author)
		assert.Equal(t, EMAIL_HOST, email.Host)

		testEmail = &emailOutgoing{
			To:       email.GetValue("var2"),
			Subject:  fmt.Sprintf("Service %v is Failing", exampleService.Name),
			Template: mainEmailTemplate,
			Data:     exampleService,
			From:     email.GetValue("var1"),
		}
	})

	t.Run("email Within Limits", func(t *testing.T) {
		ok := email.CanSend()
		assert.True(t, ok)
	})

	t.Run("email OnFailure", func(t *testing.T) {
		err := email.OnFailure(exampleService, exampleFailure)
		assert.Nil(t, err)
	})

	t.Run("email OnSuccess", func(t *testing.T) {
		err := email.OnSuccess(exampleService)
		assert.Nil(t, err)
	})

	t.Run("email Check Back Online", func(t *testing.T) {
		assert.True(t, exampleService.Online)
	})

	t.Run("email OnSuccess Again", func(t *testing.T) {
		err := email.OnSuccess(exampleService)
		assert.Nil(t, err)
	})

	t.Run("email Test", func(t *testing.T) {
		err := email.OnTest()
		assert.Nil(t, err)
	})

}
