package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
)

// DB is a Postgres database which implements Repository interface.
type DB struct {
	conn *sql.DB
}

func (p *DB) SaveToDisk(filePath string) error {
	return errors.New("saving to disk is not supported")
}

func (p *DB) RestoreFromDisk(filePath string) error {
	return errors.New("restoring from disk is not supported")
}

func (p *DB) SupportsSavingToDisk() bool {
	return false
}

func (p *DB) SupportsTx() bool {
	return true
}

func (p *DB) Updates(ctx context.Context, metrics []models.Metric) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	updateGaugeStmt, err := tx.PrepareContext(ctx, "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET type = $2, value = $3")
	if err != nil {
		return err
	}
	defer updateGaugeStmt.Close()

	updateCounterStmt, err := tx.PrepareContext(ctx, "INSERT INTO metrics (name, type, delta) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET type = $2, delta = metrics.delta + $3")
	if err != nil {
		return err
	}
	defer updateCounterStmt.Close()

	for _, metric := range metrics {
		switch metric.MType {
		case models.Gauge:
			if metric.Value == nil {
				return apierror.InvalidValue
			}
			_, err := updateGaugeStmt.Exec(metric.Name, metric.MType, *metric.Value)
			if err != nil {
				return err
			}
		case models.Counter:
			if metric.Delta == nil {
				return apierror.InvalidValue
			}
			_, err := updateCounterStmt.Exec(metric.Name, metric.MType, *metric.Delta)
			if err != nil {
				return err
			}
		default:
			return apierror.UnknownMetricType
		}
	}

	return tx.Commit()
}

func (p *DB) Ping(ctx context.Context) error {
	return p.conn.PingContext(ctx)
}

func (p *DB) Update(ctx context.Context, metric models.Metric) error {
	switch metric.MType {
	case models.Gauge:
		_, err := p.conn.ExecContext(ctx, "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET type = $2, value = $3", metric.Name, metric.MType, *metric.Value)
		if err != nil {
			return err
		}
	case models.Counter:
		_, err := p.conn.ExecContext(ctx, "INSERT INTO metrics (name, type, delta) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET type = $2, delta = metrics.delta + $3", metric.Name, metric.MType, *metric.Delta)
		if err != nil {
			return err
		}
	default:
		return apierror.UnknownMetricType
	}

	return nil
}

func (p *DB) UpdateWithStruct(ctx context.Context, metric *models.Metric) error {
	if metric == nil || len(metric.Name) == 0 {
		return apierror.InvalidValue
	}

	var err error
	switch metric.MType {
	case models.Gauge:
		if metric.Value == nil || metric.Delta != nil {
			return apierror.InvalidValue
		}

		_, err = p.conn.ExecContext(ctx, "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET type = $2, value = $3", metric.Name, metric.MType, *metric.Value)
	case models.Counter:
		if metric.Delta == nil || metric.Value != nil {
			return apierror.InvalidValue
		}

		_, err = p.conn.ExecContext(ctx, "INSERT INTO metrics (name, type, delta) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET type = $2, delta = metrics.delta + $3", metric.Name, metric.MType, *metric.Delta)
	default:
		return apierror.UnknownMetricType
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *DB) Get(ctx context.Context, key string, valueType models.MetricType) (*models.Metric, error) {
	var metric models.Metric
	var result *sql.Row
	var err error
	switch valueType {
	case models.Gauge:
		result = p.conn.QueryRowContext(ctx, "SELECT name, type, value FROM metrics WHERE name = $1 AND type = $2", key, valueType)
		metric.Value = new(float64)
		err = result.Scan(&metric.Name, &metric.MType, &metric.Value)
	case models.Counter:
		result = p.conn.QueryRowContext(ctx, "SELECT name, type, delta FROM metrics WHERE name = $1 AND type = $2", key, valueType)
		metric.Delta = new(int64)
		err = result.Scan(&metric.Name, &metric.MType, &metric.Delta)
	default:
		return nil, apierror.NotFound
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apierror.NotFound
		}
		return nil, err
	}

	return &metric, nil
}

func (p *DB) GetAll(ctx context.Context) ([]models.Metric, error) {
	rows, err := p.conn.QueryContext(ctx, "SELECT name, type, value, delta FROM metrics")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]models.Metric, 0)
	for rows.Next() {
		var metric models.Metric
		var value sql.NullFloat64
		var delta sql.NullInt64

		err := rows.Scan(&metric.Name, &metric.MType, &value, &delta)
		if err != nil {
			return nil, err
		}

		switch metric.MType {
		case models.Gauge:
			if !value.Valid {
				return nil, fmt.Errorf("invalid float64 value, got error: %s", err)
			}

			metric.Value = &value.Float64
		case models.Counter:
			if !delta.Valid {
				return nil, fmt.Errorf("invalid int64 value, got error: %s", err)
			}

			metric.Delta = &delta.Int64
		default:
			return nil, apierror.UnknownMetricType
		}

		result = append(result, metric)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func New(url string) (*DB, error) {
	if len(url) == 0 {
		log.Printf("No database url provided, skipping database initialization")
		return nil, nil
	}

	conn, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := conn.ExecContext(ctx, CreateTable); err != nil {
		return nil, fmt.Errorf("couldn't create default sqlite tables with error %s", err.Error())
	}

	return &DB{conn}, nil
}

func NewFromDB(db *sql.DB) (*DB, error) {
	return &DB{db}, nil
}
