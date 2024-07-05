package test_table

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	cloudRedis "github.com/louvri/gold/cloud_redis"
	"github.com/louvri/gosl"
	model "github.com/louvri/gosl-gen/test/generated/model/test_table"
	base "github.com/louvri/gosl-gen/test/generated/repository"
	"github.com/louvri/gosl/builder"
)

var r base.Repository
var ctx context.Context
var dbType = "mysql"

type key int

const NOSQL_KEY key = 1988

func init() {
	var logger log.Logger
	r = New(logger)
	switch dbType {
	case "redis":
		db, err := cloudRedis.New("localhost", "", "6379", time.Duration(48*60*60))
		if err != nil {
			panic(err)
		}
		ctx = context.WithValue(ctx, NOSQL_KEY, db)
	case "pgsql":
		user := "root"
		password := "abcd"
		dbname := "test"
		host := "localhost"
		port := 5432

		// Construct the connection string
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
			user, password, dbname, host, port)

		// Connect to the database using sqlx
		db, err := sqlx.Connect("postgres", connStr)
		if err != nil {
			panic(err)
		}
		ctx = context.WithValue(ctx, gosl.SQL_KEY, db)
	default:
		db, err := sqlx.Connect("mysql", fmt.Sprintf(
			"%s:%s@(%s:%s)/%s",
			"root",
			"abcd",
			"localhost",
			"3306",
			"test"))
		if err != nil {
			panic(err)
		}
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(0)
		db.SetConnMaxLifetime(0)
		db.SetConnMaxIdleTime(0)
		ctx = context.WithValue(ctx, gosl.SQL_KEY, db)
	}
}

func TestGet(t *testing.T) {
	db := ctx.Value(gosl.SQL_KEY).(*sqlx.DB)

	// empty data
	if _, err := db.Exec("DELETE FROM test_table WHERE field_value IN(?,?,?)", "'testGet1'", "'testGet2'", "'testGet3'"); err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	// input yout data here
	data := []model.Model{
		{
			FieldValue: "testGet1",
		},
		{
			FieldValue: "testGet2",
		},
		{
			FieldValue: "testGet3",
		},
	}
	query := `INSERT INTO test_table(`
	for _, item := range data {
		tmp := query
		i := item.ToMap(nil)
		keys := []string{}
		values := []string{}

		for key, val := range i {
			keys = append(keys, key)
			values = append(values, fmt.Sprintf("'%s'", val))
		}
		tmp = fmt.Sprintf("%s%s) VALUES(%s)", tmp, strings.Join(keys, ","), strings.Join(values, ","))
		_, err := db.Exec(tmp)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
		}
		keys = nil
		values = nil
	}

	res, err := r.Get(ctx, builder.QueryParams{}, nil)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if res == nil {
		t.Log("data shouldn't be empty")
		t.FailNow()
	}
}

func TestAll(t *testing.T) {
	db := ctx.Value(gosl.SQL_KEY).(*sqlx.DB)

	// empty data
	if _, err := db.Exec("DELETE FROM test_table WHERE field_value IN(?,?,?)", "'testAll1'", "'testAll2'", "'testAll3'"); err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	// input yout data here
	data := []model.Model{
		{
			FieldValue: "testAll1",
		},
		{
			FieldValue: "testAll2",
		},
		{
			FieldValue: "testAll3",
		},
	}
	query := `INSERT INTO test_table(`
	for _, item := range data {
		tmp := query
		i := item.ToMap(nil)
		keys := []string{}
		values := []string{}

		for key, val := range i {
			keys = append(keys, key)
			values = append(values, fmt.Sprintf("'%s'", val))
		}
		tmp = fmt.Sprintf("%s%s) VALUES(%s)", tmp, strings.Join(keys, ","), strings.Join(values, ","))
		_, err := db.Exec(tmp)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
		}
		keys = nil
		values = nil
	}

	res, err := r.Get(ctx, builder.QueryParams{}, nil)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if res == nil {
		t.Log("data shouldn't be empty")
		t.FailNow()
	}
}
func TestQuery(t *testing.T) {
	db := ctx.Value(gosl.SQL_KEY).(*sqlx.DB)

	// empty data
	if _, err := db.Exec("DELETE FROM test_table WHERE field_value IN(?,?,?)", "'testQuery1'", "'testQuery2'", "'testQuery3'"); err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	// input yout data here
	data := []model.Model{
		{
			FieldValue: "testQuery1",
		},
		{
			FieldValue: "testQuery2",
		},
		{
			FieldValue: "testQuery3",
		},
	}
	query := `INSERT INTO test_table(`
	for _, item := range data {
		tmp := query
		i := item.ToMap(nil)
		keys := []string{}
		values := []string{}

		for key, val := range i {
			keys = append(keys, key)
			values = append(values, fmt.Sprintf("'%s'", val))
		}
		tmp = fmt.Sprintf("%s%s) VALUES(%s)", tmp, strings.Join(keys, ","), strings.Join(values, ","))
		_, err := db.Exec(tmp)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
		}
		keys = nil
		values = nil
	}

	res, err := r.Get(ctx, builder.QueryParams{}, nil)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if res == nil {
		t.Log("data shouldn't be empty")
		t.FailNow()
	}
}

func TestSet(t *testing.T) {

}

func TestCount(t *testing.T) {

}

func TestInsert(t *testing.T) {}

func TestUpdate(t *testing.T) {}

func TestUpsert(t *testing.T) {}

func TestDelete(t *testing.T) {}
