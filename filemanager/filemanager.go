package filemanager

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/brutalzinn/go-mqtt-integration/awshelper"
)

type FileManager struct {
	File File
	AWS  AWS
}

type File struct {
	Name string
	Dir  string
	Data []byte
}
type AWS struct {
	Region  string
	Bucket  string
	Enabled bool
}

func New(name string) *FileManager {
	fm := &FileManager{
		File: File{
			Name: name,
		},
		AWS: AWS{
			Enabled: false,
		},
	}
	return fm
}

func (fm *FileManager) SetBytes(data []byte) {
	fm.File.Data = data
}

func (fm *FileManager) SetReader(r io.Reader) error {
	bodyBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	fm.File.Data = bodyBytes
	return nil
}

func (fm *FileManager) SetAWS(aws *AWS) {
	fm.AWS = *aws
}

// /TODO: finish the implementation of the Open, Write and Delete methods
func (fm *FileManager) Open() ([]byte, error) {
	if fm.AWS.Enabled {
		data, err := awshelper.S3GetObject(fm.AWS.Region, fm.AWS.Bucket, fm.File.Name)
		return data, err
	}
	file, err := os.Open(filepath.Join(fm.File.Dir, fm.File.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

func (fm *FileManager) Write() error {
	if fm.AWS.Enabled {
		return awshelper.S3PutObject(fm.AWS.Region, fm.AWS.Bucket, fm.File.Name, fm.File.Data)
	}
	file, err := os.Create(filepath.Join(fm.File.Dir, fm.File.Name))
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(fm.File.Data)
	return nil
}

func (fm *FileManager) Delete() error {
	if fm.AWS.Enabled {
		return awshelper.S3DeleteObject(fm.AWS.Region, fm.AWS.Bucket, fm.File.Name)
	}
	return os.Remove(filepath.Join(fm.File.Dir, fm.File.Name))
}
