package server

import (
	"code_challenge1/db"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func setupDbTest() {
	os.Setenv("TEST_ENV", "true")
	s := `CREATE TABLE IF NOT EXISTS "users" (
	"id" INTEGER NOT NULL UNIQUE ,
	"name" CHAR(256) NOT NULL UNIQUE,
	"balance" INTEGER NOT NULL,
	PRIMARY KEY("id")
);


CREATE TABLE IF NOT EXISTS "records" (
	"id" INTEGER NOT NULL UNIQUE,
	"from_user" INTEGER NOT NULL,
	"to_user" INTEGER NOT NULL,
	"amount" INTEGER NOT NULL,
	PRIMARY KEY("id")
);
`
	db.Schema = s
}

func TestNewServer(t *testing.T) {
	_, err := NewServer()
	assert.NotNil(t, err)

	setupDbTest()

	_, err = NewServer()
	assert.Nil(t, err)
}

func TestServer_Serve(t *testing.T) {
	setupDbTest()

	ss, err := NewServer()
	assert.Nil(t, err)

	err = ss.Serve("rrweeqw")
	assert.NotNil(t, err)

	//
}

func TestServer_AddUser(t *testing.T) {
	setupDbTest()

	ss, err := NewServer()
	assert.Nil(t, err)
	ss.router()
	router := ss.r

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/add", strings.NewReader(`{"name":"test1", "balance":"100"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res := toResponse(w.Body.Bytes())
	assert.Equal(t, 0, res.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/user/add", strings.NewReader(``))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res = toResponse(w.Body.Bytes())
	assert.NotEqual(t, 0, res.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/user/add", strings.NewReader(`{"name":"test1", "balance":"qwqwqw"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res = toResponse(w.Body.Bytes())
	assert.NotEqual(t, 0, res.Code)

}

func TestServer_UserBalance(t *testing.T) {
	setupDbTest()

	ss, err := NewServer()
	assert.Nil(t, err)
	ss.router()
	router := ss.r
	u, _ := ss.db.AddUser("name1", big.NewRat(100, 1))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/balance", strings.NewReader(fmt.Sprintf(`{"user_id":%d}`, u.ID)))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res := toResponse(w.Body.Bytes())
	assert.Equal(t, 0, res.Code)
	data := res.Data.(map[string]interface{})
	assert.Equal(t, "100.00", data["balance"].(string))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/user/balance", strings.NewReader(``))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res = toResponse(w.Body.Bytes())
	assert.NotEqual(t, 0, res.Code)
}

func TestServer_WithdrawOrDeposit(t *testing.T) {
	setupDbTest()
	ss, err := NewServer()
	assert.Nil(t, err)
	ss.router()
	router := ss.r

	u1, _ := ss.db.AddUser("name1", big.NewRat(100, 1))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/deposit", strings.NewReader(fmt.Sprintf(`{"id":%d, "amount":"1"}`, u1.ID)))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res := toResponse(w.Body.Bytes())
	assert.Equal(t, 0, res.Code)

	//json error
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/deposit", strings.NewReader(``))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res = toResponse(w.Body.Bytes())
	assert.NotEqual(t, 0, res.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/deposit", strings.NewReader(fmt.Sprintf(`{"id":%d, "amount":"qwe"}`, u1.ID)))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res = toResponse(w.Body.Bytes())
	assert.NotEqual(t, 0, res.Code)

}

func TestServer_Transfer(t *testing.T) {
	setupDbTest()
	ss, err := NewServer()
	assert.Nil(t, err)
	ss.router()
	router := ss.r

	u1, _ := ss.db.AddUser("name1", big.NewRat(100, 1))
	u2, _ := ss.db.AddUser("name2", big.NewRat(100, 1))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transfer",
		strings.NewReader(fmt.Sprintf(`{"from_user_id":%d, "to_user_id":%d, "amount":"1"}`, u1.ID, u2.ID)))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res := toResponse(w.Body.Bytes())
	assert.Equal(t, 0, res.Code)

	//json error
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/transfer", strings.NewReader(``))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res = toResponse(w.Body.Bytes())
	assert.NotEqual(t, 0, res.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/transfer",
		strings.NewReader(fmt.Sprintf(`{"from_user_id":%d, "to_user_id":%d, "amount":"11qwe"}`, u1.ID, u2.ID)))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res = toResponse(w.Body.Bytes())
	assert.NotEqual(t, 0, res.Code)

}

func TestServer_UserRecords(t *testing.T) {
	setupDbTest()
	ss, err := NewServer()
	assert.Nil(t, err)
	ss.router()
	router := ss.r

	u1, _ := ss.db.AddUser("name1", big.NewRat(100, 1))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/deposit", strings.NewReader(fmt.Sprintf(`{"id":%d, "amount":"1"}`, u1.ID)))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res := toResponse(w.Body.Bytes())
	assert.Equal(t, 0, res.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/records", strings.NewReader(fmt.Sprintf(`{"user_id":%d}`, u1.ID)))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res = toResponse(w.Body.Bytes())
	assert.Equal(t, 0, res.Code)
	assert.Equal(t, 1, len(res.Data.([]interface{})))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/records", strings.NewReader(fmt.Sprintf(`{%d}`, u1.ID)))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res = toResponse(w.Body.Bytes())
	assert.NotEqual(t, 0, res.Code)

}

func toResponse(data []byte) Response {
	var r Response
	_ = json.Unmarshal(data, &r)
	return r
}
