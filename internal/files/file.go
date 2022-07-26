package files

import (
	"os"
)

func WriteDataToFile(data []byte, to string) error {
  if err :=  EnsureParentDirectory(to); err != nil {
    return err
  }

  return os.WriteFile(to, data, 0644)
}
