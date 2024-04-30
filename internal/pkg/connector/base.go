package connector

import "errors"

var (
	ErrInvalidType = errors.New("invalid connector type provided")
)

type IConnector interface {
	Create() error
	Read() ([]byte, error)
	Update(string) error
	Delete() error
}

func GetConnector(conn_type string, conn_string string) (IConnector, error) {
	switch conn_type {
	case "file":
		fileConnector := FileConnector{
			filepath: conn_string,
		}
		return fileConnector, nil
	default:
		return nil, ErrInvalidType
	}
}
