package connector

import (
	"io/fs"
	"os"
)

type FileConnector struct {
	filepath string
}

func (f FileConnector) Create() error {
	if _, err := os.Create(f.filepath); err != nil {
		return err
	}
	return nil
}

func (f FileConnector) Read() ([]byte, error) {
	data, err := os.ReadFile(f.filepath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (f FileConnector) Update(dataString string) error {
	if err := os.WriteFile(f.filepath, []byte(dataString), fs.FileMode(os.O_WRONLY)); err != nil {
		return err
	}

	return nil

}

func (f FileConnector) Delete() error {
	if err := os.Remove(f.filepath); err != nil {
		return err
	}

	return nil
}
