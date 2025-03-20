package utils

import (
	"io"
	"maa-server/config"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"github.com/electricbubble/gadb"
	cp "github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
)

func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = CreateNestedFile(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}

	err = dstfd.Sync()
	if err != nil {
		return err
	}
	return nil
}

func CopyDir(src string, dst string) error {
	if err := cp.Copy(src, dst); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CreateNestedDirectory(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Errorf("can't create folder, %s", err)
	}
	return err
}

// CreateNestedFile create nested file
func CreateNestedFile(path string) (*os.File, error) {
	basePath := filepath.Dir(path)
	if err := CreateNestedDirectory(basePath); err != nil {
		return nil, err
	}
	return os.Create(path)
}

func DeleteDirSub(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Println("Error reading directory:", err)
		return err
	}

	for _, entry := range entries {
		entryPath := path + "/" + entry.Name()
		err := os.RemoveAll(entryPath)
		if err != nil {
			log.Println("Error removing", entryPath, ":", err)
			return err
		}
	}
	return nil
}

func DeleteFileOrDir(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		log.Println(err)
	}
	return err
}

func RenameFile(oldName string, newName string) error {
	err := os.Rename(oldName, newName)
	if err != nil {
		log.Println(err)
	}
	return err
}

func PopSlice[T any](s *[]T, index int) {
	if index < 0 || index >= len(*s) {
		return
	}
	*s = append((*s)[:index], (*s)[index+1:]...)
}

func AddOneMonth(t time.Time) time.Time {
	year, month, day := t.Date()
	location := t.Location()

	// 获取下个月的时间
	nextMonth := month + 1
	if nextMonth > 12 {
		nextMonth = 1
		year++
	}

	// 获取下个月的最后一天
	firstOfNextMonth := time.Date(year, nextMonth, 1, 0, 0, 0, 0, location)
	lastOfNextMonth := firstOfNextMonth.AddDate(0, 1, -1).Day()

	// 如果当前日期是本月的最后一天
	if day == t.AddDate(0, 1, -1).Day() {
		return time.Date(year, nextMonth, lastOfNextMonth, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), location)
	}

	// 如果下个月没有当前日期，则使用下个月的最后一天
	if day > lastOfNextMonth {
		return time.Date(year, nextMonth, lastOfNextMonth, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), location)
	}

	// 否则使用下个月的相应日期
	return time.Date(year, nextMonth, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), location)
}



func CheckGame() string {
	
	device := config.Profiles.Connection.Device
	res := strings.Split(device, ":")
	var port int
	if(len(res) == 2){
		var err error
		port, err = strconv.Atoi(res[1])
		if err != nil {
			log.Errorln(err)
			return "adb address configuration error"
		}
	}else{
		port = 5555
	}

    StartAdbDeamon()

	adbClient, err := gadb.NewClient()
	if err != nil {
		log.Println(err)
		return "gadb error"
	}
    err = adbClient.Connect(res[0],port)
    if err != nil {
		log.Println(err)
		return "adb connect error"
	}
    devices,err := adbClient.DeviceList()
    if err != nil {
		log.Println(err)
		return "gadb error"
	}
	var d *gadb.Device
    for _, de := range devices {
		if de.Serial() == device{
			d = &de
			break
		}
	}
	if d == nil {
		log.Println("device not found")
		return "device not found"
	}
	output, err := d.RunShellCommand("pm list packages | grep 'com.hypergryph.arknights.bilibili'")
	if output == "" || err != nil {
		return "game not found"
	}
	return "game ready"
}

func StopGame(d *gadb.Device) {
	_, err := d.RunShellCommand("am force-stop com.hypergryph.arknights.bilibili")
	if err != nil {
		log.Println(err)
	}
}
