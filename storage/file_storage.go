package storage

import (
	"encoding/base64"
	"os"
	"path"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type FileStorage struct {
	dir string
}

func NewFileStorage(dir string) (Storage, error) {
	confDir := path.Join(dir, "v1")

	return &FileStorage{
		dir: confDir,
	}, os.MkdirAll(confDir, os.ModePerm)
}

func (s *FileStorage) getFilePath(name string) string {
	return path.Join(s.dir, name+".conf")
}

func (s *FileStorage) WriteDeviceOptions(name string, options StoreDeviceOptions) error {
	filePath := s.getFilePath(name)
	f, err := os.CreateTemp(s.dir, name)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := options.Dump(f); err != nil {
		return err
	}

	if err := os.Chmod(f.Name(), 0600); err != nil {
		return err
	}

	return os.Rename(f.Name(), filePath)
}

func (s *FileStorage) WritePeerOptions(pubKey wgtypes.Key, options StorePeerOptions) error {
	safeName := base64.URLEncoding.EncodeToString(pubKey[:])
	filePath := s.getFilePath(safeName)
	f, err := os.CreateTemp(s.dir, safeName)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := options.Dump(f); err != nil {
		return err
	}

	if err := os.Chmod(f.Name(), 0600); err != nil {
		return err
	}

	return os.Rename(f.Name(), filePath)
}

func (s *FileStorage) ReadDeviceOptions(name string) (*StoreDeviceOptions, error) {
	f, err := os.OpenFile(s.getFilePath(name), os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	o := &StoreDeviceOptions{}
	return o, o.Restore(f)
}

func (s *FileStorage) ReadPeerOptions(pubKey wgtypes.Key) (*StorePeerOptions, error) {
	safeName := base64.URLEncoding.EncodeToString(pubKey[:])
	f, err := os.OpenFile(s.getFilePath(safeName), os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	o := &StorePeerOptions{}
	return o, o.Restore(f)
}
