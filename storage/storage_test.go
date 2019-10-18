package storage_test

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/suquant/wgrest/storage"
)

type TestData struct {
	Path        string
	Data        []byte
	StorageType storage.Type
}

func TestStorages(t *testing.T) {
	data := make([]byte, 10240)
	_, err := rand.Read(data)
	if err != nil {
		t.Errorf("rand err: %s", err.Error())
	}

	testData := []TestData{
		TestData{
			Path:        "/tmp/wgrest.tmp",
			Data:        data,
			StorageType: storage.MemoryStorage,
		},
		TestData{
			Path:        "/tmp/wgrest.tmp",
			Data:        data,
			StorageType: storage.DiskStorage,
		},
	}

	for _, td := range testData {
		testStorage(t, td)
	}

}

func testStorage(t *testing.T, td TestData) {
	s := storage.NewStorage(td.StorageType)
	t.Logf("testing %T storage", s)

	rwc, err := s.Open(td.Path)
	if err != nil {
		t.Errorf("open err(%T): %s", s, err.Error())
	}

	n, err := rwc.Write(td.Data)
	if err != nil {
		t.Errorf("write err(%T): %s", s, err.Error())
	}

	err = rwc.Close()
	if err != nil {
		t.Errorf("close err(%T): %s", s, err.Error())
	}

	rwc, err = s.Open(td.Path)
	if err != nil {
		t.Errorf("open err(%T): %s", s, err.Error())
	}

	if len(td.Data) != n {
		t.Errorf("expected bytes noq equal wrote one(%T): %v != %v", s, len(td.Data), n)
	}

	data, err := ioutil.ReadAll(rwc)
	if err != nil {
		t.Errorf("read err(%T): %s", s, err.Error())
	}

	if bytes.Compare(td.Data, data) != 0 {
		t.Errorf("wrote and read data are different: %T", s)
	}

	err = rwc.Close()
	if err != nil {
		t.Errorf("close err(%T): %s", s, err.Error())
	}

	_, err = ioutil.ReadAll(rwc)
	if err == nil {
		t.Errorf("storage closed but Read still works: %T", s)
	}

	_, err = rwc.Write(td.Data)
	if err == nil {
		t.Errorf("storage closed but Write still works: %T", s)
	}

}
