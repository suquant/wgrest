package storage_test

import (
	"github.com/suquant/wgrest/storage"
	"io/ioutil"
	"os"
	"testing"
)

func TestStorage(t *testing.T) {
	t.Run("read non existing device options", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatalf("failed to create temp dir: %s", err)
		}

		s, err := storage.NewFileStorage(dir)
		if err != nil {
			t.Fatalf("failed to create file storage: %s", err)
		}

		_, err = s.ReadDeviceOptions("xyz")
		if os.IsNotExist(err) != true {
			t.Errorf("got %s, want not exist error", err)
		}
	})
}
