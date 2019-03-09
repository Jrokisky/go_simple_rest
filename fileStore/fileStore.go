/*
 * The purpose of this package is to handle all operations
 * involved in storing data in the file system.
 */

package fileStore

import (
	"io/ioutil"
	"os"
	"strings"
)

type FileStore struct {
	prefix string
}

func (fs *FileStore) SetPrefix(prefix string) {
	fs.prefix = prefix
}

func (fs *FileStore) Load(file_name string) ([]byte, error) {
	file_data, err := ioutil.ReadFile(fs.prefix + file_name)
	return file_data, err
}

func (fs *FileStore) Write(file_name string, data []byte) error {
	err := ioutil.WriteFile(fs.prefix + file_name, data, 0666)
	return err
}

func (fs *FileStore) Delete(file_name string) error {
	return os.Remove(fs.prefix + file_name)
}

func (fs *FileStore) Exists(file_name string) (bool) {
	if _, err := os.Stat(fs.prefix + file_name); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func (fs *FileStore) GetFiles() ([]string, error) {
	files, err := ioutil.ReadDir(fs.prefix)
	var file_names []string
	if err != nil {
		return  file_names, nil
	} else {
		for _, file := range files {
			file_names = append(file_names, file.Name())
		}
		return file_names, nil
	}
}

func (fs *FileStore) RemoveTestFiles() (error) {
	files, err := ioutil.ReadDir(fs.prefix)
	if err != nil {
		return  err
	} else {
		for _, file := range files {
			file_name := file.Name()
			if strings.HasPrefix(file_name, "test_") {
				err = os.Remove(fs.prefix + file_name)
				if err != nil {
					return err
				}

			}
		}
		return nil
	}
}
