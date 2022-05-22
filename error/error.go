package main

import (
	"fmt"
	"strings"

	"xorm.io/builder"
)

func main() {
	dao, err := testDao()
	if err != nil {
		if !IsNotFoundError(err) {
			// error
			fmt.Println(err)
		}
	}

	// do something
	fmt.Println(dao)
}

// Error
func (e *NotFoundError) Error() string {
	if len(e.msg) == 0 {
		return fmt.Sprintf("%s not found", e.table)
	}

	return fmt.Sprintf("%s not found, %s", e.table, strings.Join(e.msg, ","))
}

func NewNotFoundError(table Table, msg ...string) *NotFoundError {
	return &NotFoundError{table: table.TableName(), msg: msg}
}

type Table interface {
	TableName() string
}

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

type NotFoundError struct {
	table string
	msg   []string
}

func IsNotFoundError(err error) bool {
	_, ok := Cause(err).(*NotFoundError)
	return ok
}

//TestUser  model
type TestUser struct {
	id        int64
	name      string
	scores    int
	star      int
	deletedAt int32
	createdAt int32
}

func (m *TestUser) TableName() string {
	return "test_user"
}

// testDao repositry
func testDao() (interface{}, error) {
	table := &TestUser{}
	tableName := table.TableName()
	_, result, err := builder.Select("SUM(scores) as total_scores", "SUM(stars) as total_stars").
		From(tableName).
		Where(builder.
			And(builder.Eq{"user_id": 1111}).
			And(builder.Or(builder.Eq{"deleted_at": 0}, builder.IsNull{"deleted_at"})),
		).
		GroupBy("user_id").ToSQL()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, NewNotFoundError(table, "not fund table")
	}

	return result, nil
}
