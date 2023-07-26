package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"log"
	"strconv"
	"time"
)

type DB struct {
	conn *pgx.Conn
}

func (p DB) Ping(ctx context.Context) error {
	return p.conn.Ping(ctx)
}

func (p DB) Update(name string, valueType models.MetricType, value string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return apierror.NumberParse
	}

	_, err = p.conn.Exec(ctx, "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3)", name, valueType, floatVal)
	if err != nil {
		return fmt.Errorf("couldn't insert metric with error %s", err.Error())
	}

	return nil
}

func (p DB) Get(key string, valueType models.MetricType) (*models.Metric, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 9999999*time.Second)
	defer cancelFunc()

	result := p.conn.QueryRow(ctx, "SELECT name, type, value FROM metrics WHERE name = $1 AND type = $2", key, valueType)

	var metric models.Metric
	var value pgtype.Numeric

	err := result.Scan(&metric.Name, &metric.MType, &value)
	if err != nil {
		return nil, fmt.Errorf("couldn't get metric with error %s", err.Error())
	}

	switch metric.MType {
	case models.Gauge:
		val, err := numericToFloat64(&value)
		if err != nil {
			return nil, fmt.Errorf("couldn't get metric with error %s", err.Error())
		}
		metric.Value = &val
	case models.Counter:
		val, err := numericToInt64(&value)
		if err != nil {
			return nil, fmt.Errorf("couldn't get metric with error %s", err.Error())
		}
		metric.Delta = &val
	default:
		return nil, apierror.UnknownMetricType
	}

	if err != nil {
		return nil, fmt.Errorf("couldn't get metric with error %s", err.Error())
	}

	return &metric, nil
}

func (p DB) UpdateWithStruct(metric *models.Metric) error {
	return p.Update(metric.Name, metric.MType, metric.StringValue())
}

func (p DB) GetAll() ([]models.Metric, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	rows, err := p.conn.Query(ctx, "SELECT name, type, value FROM metrics")
	if err != nil {
		return nil, err
	}

	result := make([]models.Metric, 0)
	for rows.Next() {
		var metric models.Metric
		var value pgtype.Numeric

		err := rows.Scan(&metric.Name, &metric.MType, &value)
		if err != nil {
			return nil, err
		}

		switch metric.MType {
		case models.Gauge:
			val, err := numericToFloat64(&value)
			if err != nil {
				return nil, err
			}
			metric.Value = &val
		case models.Counter:
			val, err := numericToInt64(&value)
			if err != nil {
				return nil, err
			}
			metric.Delta = &val
		default:
			return nil, apierror.UnknownMetricType
		}

		if err != nil {
			log.Printf("couldn't scan metrics row with error %s", err.Error())
			return nil, nil
		}

		result = append(result, metric)
	}
	return result, nil
}

func New(url string) (*DB, error) {
	if len(url) == 0 {
		log.Printf("No database url provided, skipping database initialization")
		return nil, nil
	}

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := conn.Exec(ctx, CreateTable); err != nil {
		log.Panicf("Couldn't create default sqlite tables with error %s", err.Error())
	}

	return &DB{conn}, nil
}
