package models

import (
	"encoding/json"
	"errors"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

func (d *Duration) String() string {
	return d.Duration.String()
}

func (d *Duration) Set(value string) error {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	d.Duration = duration
	return nil
}

func ParseDurationFromEnv(v string) (interface{}, error) {
	var parsedDuration Duration
	if err := parsedDuration.Set(v); err != nil {
		return nil, err
	}

	return parsedDuration, nil
}
