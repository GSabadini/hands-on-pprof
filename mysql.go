package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// NewMySQLConnection creates a new mysql connection
func NewMySQLConnection() *sql.DB {
	// Register an OTel driver
	driverName, err := otelsql.Register("mysql", otelsql.WithAttributes(semconv.DBSystemMySQL))
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := sql.Open(driverName, fmt.Sprintf("dev:dev@tcp(localhost:3306)/dev"))
	if err != nil {
		log.Fatal(err)
	}

	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemMySQL))
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 1)

	fmt.Println("DB Connect")
	return db
}

func Insert(name string, db *sql.DB, tracer trace.TracerProvider) (int64, error) {
	ctx, span := tracer.Tracer("instrumentationName").Start(context.Background(), "Insert")
	defer span.End()

	fail := func(err error) (int64, error) {
		return 0, fmt.Errorf("insert: %v", err)
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, "INSERT INTO Authorizations (name) VALUES (?)", name)
	if err != nil {
		return fail(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return id, nil
}

func Select(id int, db *sql.DB, tracer trace.TracerProvider) (string, error) {
	ctx, span := tracer.Tracer("instrumentationName").Start(context.Background(), "Select")
	defer span.End()

	fail := func(err error) (string, error) {
		return "", fmt.Errorf("select: %v", err)
	}

	var name string
	err := db.QueryRowContext(ctx, "SELECT name FROM Authorizations WHERE id=?", id).Scan(&name)
	switch {
	case err == sql.ErrNoRows:
		return fail(err)
	case err != nil:
		return fail(err)
	}

	return name, nil
}
