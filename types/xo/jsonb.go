package xo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Jsonb json.RawMessage

type NullJsonb struct {
	Jsonb json.RawMessage
	Valid bool
}

// Scan implements the Scanner interface.
func (nj *NullJsonb) Scan(rawValue interface{}) error {
	if rawValue == nil {
		nj.Jsonb, nj.Valid = nil, false
		return nil
	}
	nj.Valid = true
	value, err := castToBytes(rawValue)
	if err != nil {
		return err
	}

	return json.Unmarshal(value, nj.Jsonb)
}

// Value implements the driver Valuer interface.
func (nj NullJsonb) Value() (driver.Value, error) {
	if !nj.Valid {
		return nil, nil
	}
	return []byte(nj.Jsonb), nil
}

// Scan implements the Scanner interface.
func (nj *Jsonb) Scan(rawValue interface{}) error {
	if rawValue == nil {
		return errors.New("expected not null jsonb")
	}

	value, err := castToBytes(rawValue)
	if err != nil {
		return err
	}

	return json.Unmarshal(value, nj)
}

// Value implements the driver Valuer interface.
func (nj Jsonb) Value() (driver.Value, error) {
	if len(nj) == 0 {
		return nil, errors.New("null is not supported by jsonb")
	}
	return []byte(nj), nil
}

func castToBytes(src interface{}) ([]byte, error) {
	var data []byte
	switch rawData := src.(type) {
	case []byte:
		data = rawData
	case string:
		data = []byte(rawData)
	default:
		return nil, errors.New("unexpected type for jsonb")
	}

	return data, nil
}
