package server

import (
	"code_challenge1/db"
	"code_challenge1/log"
	"math/big"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewServer() (*Server, error) {
	r := gin.New()
	db1, err := db.Open()
	if err != nil {
		return nil, errors.Wrap(err, "open db")
	}
	return &Server{
		r:  r,
		db: db1,
	}, nil
}

type Server struct {
	r  *gin.Engine
	db *db.DB
}

func (s *Server) Serve() error {
	log.Infof("------- starting app server ---------")
	s.r.POST("/user/add", HttpHandler(s.AddUser))
	s.r.POST("/user/balance", HttpHandler(s.UserBalance))
	s.r.POST("/records", HttpHandler(s.UserRecords))
	s.r.POST("/deposit", HttpHandler(s.WithdrawOrDeposit))
	s.r.POST("/transfer", HttpHandler(s.Transfer))

	return s.r.Run(":8080")
}

type AddUserIn struct {
	Name    string `json:"name" binding:"required"`
	Balance string `json:"balance" binding:"required"`
}

func (s *Server) AddUser(c *gin.Context) (interface{}, error) {
	var in AddUserIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	log.Debugf("add user input: %+v", in)
	balance, ok := new(big.Rat).SetString(strings.TrimSpace(in.Balance))
	if !ok {
		return nil, errors.Errorf("cannot set balance: %s", in.Balance)
	}
	id, err := s.db.AddUser(strings.TrimSpace(in.Name), balance)
	if err != nil {
		return nil, err
	}
	return gin.H{"id": id}, nil
}

type UserBalanceIn struct {
	UserID int `json:"user_id" binding:"required"`
}

func (s *Server) UserBalance(c *gin.Context) (interface{}, error) {
	var in UserBalanceIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	u, err := s.db.GetUser(in.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "get user")
	}

	return gin.H{
		"name":    u.Name,
		"balance": u.Balance.FloatString(2),
	}, nil
}

type WithdrawOrDepositIn struct {
	ID     int    `json:"id" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}

func (s *Server) WithdrawOrDeposit(c *gin.Context) (interface{}, error) {
	var in WithdrawOrDepositIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	b, ok := new(big.Rat).SetString(strings.TrimSpace(in.Amount))
	if !ok {
		return nil, errors.Errorf("amount not valid: %s", in.Amount)
	}
	u, err := s.db.WithdrawOrDeposit(in.ID, b)
	if err != nil {
		return nil, errors.Wrap(err, "WithdrawOrDeposit")
	}

	return gin.H{
		"name":    u.Name,
		"balance": u.Balance.FloatString(2),
	}, nil
}

type TransferIn struct {
	FromUserID int    `json:"from_user_id" binding:"required"`
	ToUserID   int    `json:"to_user_id" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
}

func (s *Server) Transfer(c *gin.Context) (interface{}, error) {
	var in TransferIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	b, ok := new(big.Rat).SetString(strings.TrimSpace(in.Amount))
	if !ok {
		return nil, errors.Errorf("amount is not valid: %s", in.Amount)
	}
	err := s.db.Transfer(in.FromUserID, in.ToUserID, b)
	if err != nil {
		return nil, errors.Wrap(err, "transfer")
	}
	return "success", nil
}

type UserRecordsIn struct {
	UserID int `json:"user_id" binding:"required"`
}
type UserRecordsOut struct {
	FromUser int    `json:"from_user"`
	ToUser   int    `json:"to_user"`
	Amount   string `json:"amount"`
}

func (s *Server) UserRecords(c *gin.Context) (interface{}, error) {
	var in UserRecordsIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	his, err := s.db.UserRecords(in.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "query db")
	}
	var outs = make([]UserRecordsOut, 0, len(his))
	for _, r := range his {
		outs = append(outs, UserRecordsOut{
			FromUser: r.FromUser,
			ToUser:   r.ToUser,
			Amount:   r.Amount.FloatString(2),
		})
	}

	return outs, nil
}
