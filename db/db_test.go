package db

import (
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalanceToInt(t *testing.T) {
	b := big.NewRat(100, 1)
	i := balanceToInt(b)
	assert.Equal(t, int64(10000), i)
}

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
	schema = s
}

func TestOpen(t *testing.T) {
	setupDbTest()
	_, err := Open()
	assert.Nil(t, err)
	schema = "CREATE "
	_, err = Open()
	assert.NotNil(t, err)

}

func TestDB_AddUser(t *testing.T) {
	setupDbTest()
	db, err := Open()
	assert.Nil(t, err)
	u, err := db.AddUser("test1", big.NewRat(100, 1))
	assert.Nil(t, err)
	assert.Equal(t, "test1", u.Name)
}

func TestDB_GetUser(t *testing.T) {
	setupDbTest()
	db, err := Open()
	assert.Nil(t, err)
	u, err := db.AddUser("test1", big.NewRat(100, 1))
	assert.Nil(t, err)
	assert.Equal(t, "test1", u.Name)

	u1, err := db.GetUser(u.ID)
	assert.Nil(t, err)
	assert.Equal(t, u1.Name, u.Name)

	_, err = db.GetUser(9999)
	assert.NotNil(t, err)
}

func TestDB_WithdrawOrDeposit(t *testing.T) {
	setupDbTest()
	db, err := Open()
	assert.Nil(t, err)
	u, err := db.AddUser("test1", big.NewRat(100, 1))
	assert.Nil(t, err)
	assert.Equal(t, "test1", u.Name)
	u1, err := db.WithdrawOrDeposit(u.ID, big.NewRat(1, 1))
	assert.Nil(t, err)
	assert.Equal(t, int64(101), u1.Balance.Num().Int64())

	_, err = db.WithdrawOrDeposit(u.ID, big.NewRat(1, 1000))
	assert.NotNil(t, err)

	_, err = db.WithdrawOrDeposit(9999, big.NewRat(1, 1))
	assert.NotNil(t, err)

	_, err = db.WithdrawOrDeposit(u.ID, big.NewRat(-1000, 1))
	assert.NotNil(t, err)

}

func TestDB_UserRecords(t *testing.T) {
	setupDbTest()
	db, err := Open()
	assert.Nil(t, err)
	u, err := db.AddUser("test1", big.NewRat(100, 1))
	assert.Nil(t, err)
	assert.Equal(t, "test1", u.Name)
	u1, err := db.WithdrawOrDeposit(u.ID, big.NewRat(1, 1))
	assert.Nil(t, err)
	assert.Equal(t, int64(101), u1.Balance.Num().Int64())
	r, err := db.UserRecords(u.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(r))
}

func TestDB_Transfer(t *testing.T) {
	setupDbTest()
	db, err := Open()
	assert.Nil(t, err)
	u1, err := db.AddUser("test1", big.NewRat(100, 1))
	assert.Nil(t, err)
	assert.Equal(t, "test1", u1.Name)
	u2, err := db.AddUser("test2", big.NewRat(100, 1))
	assert.Nil(t, err)
	assert.Equal(t, "test2", u2.Name)
	err = db.Transfer(u1.ID, u2.ID, big.NewRat(1, 1))
	assert.Nil(t, err)

	err = db.Transfer(u1.ID, u2.ID, big.NewRat(1, 1000))
	assert.NotNil(t, err)

	err = db.Transfer(u1.ID, u2.ID, big.NewRat(-1, 1))
	assert.NotNil(t, err)
	err = db.Transfer(u1.ID, u2.ID, big.NewRat(10000, 1))
	assert.NotNil(t, err)

	err = db.Transfer(9999, u2.ID, big.NewRat(1, 1))
	assert.NotNil(t, err)
	err = db.Transfer(u1.ID, 9999, big.NewRat(1, 1))
	assert.NotNil(t, err)

}
