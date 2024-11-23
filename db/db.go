package db

import (
	"code_challenge1/log"
	"context"
	"database/sql"
	_ "embed"
	"math/big"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const CurrencyDecimal = 2

//go:embed schema.sql
var schema string

type User struct {
	ID      int
	Name    string
	Balance *big.Rat
}

type Record struct {
	ID       int
	FromUser int
	ToUser   int
	Amount   *big.Rat
}

func Open() (*DB, error) {
	connectInfo := os.Getenv("DB_CONNECT_INFO")

	return open(connectInfo)

}

func open(conninfo string) (*DB, error) {
	log.Debugf("connect string: %s", conninfo)
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		return nil, errors.Wrap(err, "open db")
	}
	_, err = db.Exec(schema)
	if err != nil {
		return nil, errors.Wrap(err, "create schema")
	}
	log.Infof("open database success!")
	return &DB{db: db}, nil
}

type DB struct {
	db *sql.DB
}

func (d *DB) AddUser(name string, balance *big.Rat) (*User, error) {
	b := balanceToInt(balance)
	tx, err := d.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec("INSERT INTO users(name, balance) VALUES ($1, $2)", name, b)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	row := tx.QueryRow("SELECT * FROM users WHERE name=$1", name)
	var u User
	var b1 int64
	err = row.Scan(&u.ID, &u.Name, &b1)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	u.Name = strings.TrimSpace(u.Name)
	u.Balance = IntToBalance(b1)

	tx.Commit()

	return &u, nil
}

func (d *DB) GetUser(id int) (*User, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE id=$1", id)
	var u User
	var b1 int64
	err := row.Scan(&u.ID, &u.Name, &b1)
	if err != nil {
		return nil, err
	}
	u.Name = strings.TrimSpace(u.Name)
	u.Balance = IntToBalance(b1)
	return &u, nil
}

func (d *DB) WithdrawOrDeposit(id int, amount *big.Rat) (*User, error) {
	if n, ok := amount.FloatPrec(); n > 2 || !ok {
		return nil, errors.Errorf("amount should only have atmost 2 decimal number, eg. 10.02")
	}

	u, err := d.GetUser(id)
	if err != nil {
		return nil, errors.Wrap(err, "get user")
	}
	b1 := new(big.Rat).Add(u.Balance, amount)
	if b1.Sign() < 0 {
		return nil, errors.Errorf("cannot withdraw larger than balance, balance is: %v", u.Balance.FloatString(2))
	}

	b := balanceToInt(b1)

	tx, err := d.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "create tx")
	}
	_, err = tx.Exec("UPDATE users SET balance=$1 WHERE id=$2", b, id)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "update balance")
	}

	_, err = tx.Exec("INSERT INTO records (from_user, to_user, amount) VALUES ($1, $1, $2)", id, balanceToInt(amount))
	if err != nil {
		return nil, errors.Wrap(err, "add record")
	}
	tx.Commit()
	u, err = d.GetUser(id)
	if err != nil {
		return nil, errors.Wrap(err, "get user")
	}

	return u, nil
}

func (d *DB) UserRecords(userID int) ([]Record, error) {
	var records []Record

	rows, err := d.db.Query("SELECT * FROM records WHERE from_user=$1 OR to_user=$1", userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var his Record
		var b int64
		err = rows.Scan(&his.ID, &his.FromUser, &his.ToUser, &b)
		if err != nil {
			return nil, err
		}
		his.Amount = IntToBalance(b)

		records = append(records, his)
	}
	return records, nil
}

func (d *DB) Transfer(fromId, toId int, amount *big.Rat) error {
	if n, ok := amount.FloatPrec(); n > 2 || !ok {
		return errors.Errorf("amount should only have atmost 2 decimal number, eg. 10.02")
	}
	if amount.Sign() < 0 {
		return errors.Errorf("transfer amount should not be negtive: %v", amount.FloatString(2))
	}
	fromUser, err := d.GetUser(fromId)
	if err != nil {
		return errors.Wrap(err, "get from user")
	}
	toUser, err := d.GetUser(toId)
	if err != nil {
		return errors.Wrap(err, "get to user")
	}
	newFromUserBalance := new(big.Rat).Sub(fromUser.Balance, amount)
	if newFromUserBalance.Sign() < 0 {
		return errors.Errorf("user balance is not sufficient")
	}
	newToUserBalance := new(big.Rat).Add(toUser.Balance, amount)

	ctx := context.TODO()
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users SET balance=$1 WHERE id=$2", balanceToInt(newFromUserBalance), fromId)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update from user")
	}
	_, err = tx.Exec("UPDATE users SET balance=$1 WHERE id=$2", balanceToInt(newToUserBalance), toId)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update to user")
	}

	b := balanceToInt(amount)
	_, err = tx.Exec("INSERT INTO records (from_user, to_user, amount) VALUES ($1, $2, $3)", fromId, toId, b)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil

}

func balanceToInt(b *big.Rat) int64 {
	return new(big.Rat).Mul(b, big.NewRat(100, 1)).Num().Int64()
}

func IntToBalance(b int64) *big.Rat {
	return big.NewRat(b, 100)
}
