package filemanager

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/brutalzinn/go-mqtt-integration/awshelper"
)

type FileManager struct {
	IsAWS bool
	File  File
	AWS   AWS
}

type File struct {
	Name string
	Data []byte
}
type AWS struct {
	Region string
	Bucket string
}

func New(name string) *FileManager {
	fm := &FileManager{
		IsAWS: false,
	}
	fm.File.Name = name
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

func (fm *FileManager) SetAWS(awsConfig AWS) {
	fm.AWS = awsConfig
	fm.IsAWS = true
}

func (fm *FileManager) Load() ([]byte, error) {
	if fm.IsAWS {
		data, err := awshelper.S3GetObject(context.TODO(), fm.AWS.Region, fm.AWS.Bucket, fm.File.Name)
		return data, err
	}
	file, err := os.Open(fm.File.Name)
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

func (fm *FileManager) Save() error {
	if fm.IsAWS {
		return awshelper.S3PutObject(context.TODO(), fm.AWS.Region, fm.AWS.Bucket, fm.File.Name, fm.File.Data)
	}
	file, err := os.Create(fm.File.Name)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(fm.File.Data)
	return nil
}
