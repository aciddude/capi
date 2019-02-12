package datastore

import "os"

// fileExists checks to see if a file exists at the provided filepath.
func fileExists(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		// file does not exist.
		return false, nil
	}

	if err != nil {
		// Any other error indicates system timeout/permission error
		return false, err
	}

	return true, nil
}
