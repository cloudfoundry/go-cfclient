package cfclient

import (
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListV3UserByQuery(t *testing.T) {
	Convey("List V3 Users by Query", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/users", []string{listV3UsersPayload}, "", http.StatusOK, "", nil},
			{"GET", "/v3/userspage2", []string{listV3UsersPayloadPage2}, "", http.StatusOK, "page=2&per_page=2", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		query := url.Values{}
		users, err := client.ListV3UsersByQuery(query)
		So(err, ShouldBeNil)
		So(users, ShouldHaveLength, 3)

		So(users[0].Username, ShouldEqual, "smoke_tests")
		So(users[1].Username, ShouldEqual, "test1")
		So(users[2].Username, ShouldEqual, "test2")
	})

}
