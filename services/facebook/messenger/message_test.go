package messenger

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func newPayload(jsonStr ...string) io.Reader {
	const common = `{
  "object":"page",
  "entry":[
    {
      "id":"123456",
      "time":1458692752478,
      "messaging":[%s]
    }
  ]
}`
	s := fmt.Sprintf(common, strings.Join(jsonStr, ","))
	return bytes.NewBufferString(s)
}

func Test_AccountLinkingMessage(t *testing.T) {
	a := assert.New(t)
	messages, err := Parse(newPayload(`{
	"sender":{
		"id":"USER_ID"
	},
	"recipient":{
		"id":"PAGE_ID"
	},
	"timestamp":1234567890,
	"account_linking":{
		"status":"linked",
		"authorization_code":"PASS_THROUGH_AUTHORIZATION_CODE"
	}
}`))
	a.Nil(err)
	a.EqInt(1, len(messages))
	linking, _ := messages[0].Content.(*AccountLinking)
	a.NotNil(linking, "Account Linking")
	a.EqStr(string(AccountStatusLinked), string(linking.Status))
	a.EqStr("PASS_THROUGH_AUTHORIZATION_CODE", linking.AuthorizationCode)
}
