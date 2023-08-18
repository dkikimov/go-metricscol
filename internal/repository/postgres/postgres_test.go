package postgres

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/server/apierror"
	"testing"
	"time"
)

// TODO: Кажется не совсем правильно, что я пишу запрос ручками. А вдруг он изменится? Поискать другой способ.

func TestDB_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery(`SELECT name, type, value FROM metrics`).
		WithArgs("Alloc", models.Gauge).
		WillReturnRows(sqlmock.NewRows([]string{"name", "type", "value"}).AddRow("Alloc", models.Gauge, 101.42))

	mock.ExpectQuery(`SELECT name, type, delta FROM metrics`).
		WithArgs("PollCount", models.Counter).
		WillReturnRows(sqlmock.NewRows([]string{"name", "type", "delta"}).AddRow("PollCount", models.Counter, 1))

	mock.ExpectQuery(`SELECT name, type, delta FROM metrics`).
		WithArgs("Alloc", models.Counter).
		WillReturnError(apierror.NotFound)

	postgres, err := NewFromDB(db)
	require.NoError(t, err)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	repository.TestGet(ctx, t, postgres)
	require.NoError(t, mock.ExpectationsWereMet())

}

func TestDB_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery(`SELECT name, type, value, delta FROM metrics`).
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "type", "value", "delta"}).
				AddRow("Alloc", models.Gauge, 101.42, sql.NullInt64{}).
				AddRow("PollCount", models.Counter, sql.NullFloat64{}, 1),
		)

	postgres, err := NewFromDB(db)
	require.NoError(t, err)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	repository.TestGetAll(ctx, t, postgres)
	require.NoError(t, mock.ExpectationsWereMet())

}

func TestDB_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()
	mock.MatchExpectationsInOrder(false)

	// TODO: Подумать как сделать проверку на то, что в запросе есть все нужные поля
	mock.ExpectExec("INSERT INTO metrics").
		WithArgs("Alloc", models.Gauge, 120.123).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO metrics").
		WithArgs("PollCount", models.Counter, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	postgres, err := NewFromDB(db)
	require.NoError(t, err)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	repository.TestUpdate(ctx, t, postgres)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_UpdateWithStruct(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()
	mock.MatchExpectationsInOrder(false)

	mock.ExpectExec("INSERT INTO metrics").
		WithArgs("Alloc", models.Gauge, 120.123).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT name, type, value FROM metrics").
		WithArgs("Alloc", models.Gauge).
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "type", "value"}).
				AddRow("Alloc", models.Gauge, 120.123),
		)

	mock.ExpectExec("INSERT INTO metrics").
		WithArgs("PollCount", models.Counter, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT name, type, delta FROM metrics").
		WithArgs("PollCount", models.Counter).
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "type", "delta"}).
				AddRow("PollCount", models.Counter, 2),
		)

	postgres, err := NewFromDB(db)
	require.NoError(t, err)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	repository.TestUpdateWithStruct(ctx, t, postgres)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_Updates(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	// #1
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO metrics")
	mock.ExpectPrepare("INSERT INTO metrics")

	mock.ExpectExec("INSERT INTO metrics").
		WithArgs("Alloc", models.Gauge, 120.123).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO metrics").
		WithArgs("PollCount", models.Counter, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery("SELECT name, type, value, delta FROM metrics").
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "type", "value", "delta"}).
				AddRow("Alloc", models.Gauge, 120.123, sql.NullInt64{}).
				AddRow("PollCount", models.Counter, sql.NullFloat64{}, 1),
		)
	// #2

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO metrics")
	mock.ExpectPrepare("INSERT INTO metrics")

	mock.ExpectExec("INSERT INTO metrics").
		WithArgs("Alloc", models.Gauge, 120.123).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectRollback()

	mock.ExpectQuery("SELECT name, type, value, delta FROM metrics").
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "type", "value", "delta"}),
		)
	//	#3

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO metrics")
	mock.ExpectPrepare("INSERT INTO metrics")

	mock.ExpectExec("INSERT INTO metrics").
		WithArgs("Alloc", models.Gauge, 120.123).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectRollback()

	mock.ExpectQuery("SELECT name, type, value, delta FROM metrics").
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "type", "value", "delta"}),
		)

	postgres, err := NewFromDB(db)
	require.NoError(t, err)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	repository.TestUpdates(ctx, t, postgres)

	require.NoError(t, mock.ExpectationsWereMet())
}
