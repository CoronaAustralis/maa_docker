package scheduler

import (
	"io"
	"maa-server/config"
	"os"
	"path/filepath"
)

func ReadProfile() (string,error) {
	fs, err := os.Open(filepath.Join(config.D.ProfilesDir, "default.toml"))
	if(err != nil){
		return "",err
	}
	defer fs.Close()
	bytes, err := io.ReadAll(fs)
	if(err != nil){
		return "",err
	}
	return string(bytes),nil
}
