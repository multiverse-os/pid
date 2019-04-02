package pid

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
)

const defaultAppName = "service"

type File struct {
	*os.File
	Path string
	Pid  int
}

func New(pidPath string) *File {
	pidPath = ValidatePath(pidPath)
	return &File{
		Path: pidPath,
		Pid:  os.Getpid(),
	}
}

func ValidatePath(pidPath string) string {
	if len(pidPath) < 0 || len(pidPath) > 256 {
		return OSDefault(defaultAppName)
	}
	basename := path.Base(pidPath)
	if basename[len(basename)-1:] == "/" {
		return pidPath + defaultAppName + ".pid"
	} else if filepath.Ext(basename) != ".pid" {
		return pidPath + ".pid"
	} else {
		return pidPath
	}
}

func appName(pidPath string) string {
	basename := path.Base(pidPath)
	extension := filepath.Ext(basename)
	return basename[0 : len(basename)-len(extension)]
}

// TODO: Test this; it may not work if locked, so below fileless clean may fail
func removeFile(path string) error {
	if err := os.Remove(path); err != nil {
		return errCleanFailed
	}
	return nil
}

//[ Pid File Location Helpers ]////////////////////////////////////////////////
func OSDefault(app string) string   { return ("/var/run/" + app + "/" + app + ".pid") }
func TempDefault(app string) string { return ("/var/tmp/" + app + ".pid") }
func UserDefault(app string) string {
	user, err := user.Current()
	if err != nil {
		return TempDefault(app)
	}
	return ("/run/" + user.Name + "/" + app + ".pid")
}

//[ File Methods ]/////////////////////////////////////////////////////////////
func (self *File) Clean() error {
	Unlock(self.File.Fd())
	self.File.Close()
	return removeFile(self.Path)
}

func Write(pidPath string) (*File, error) {
	pid := New(pidPath)

	// NOTE: Confirm path exists, if does not exist write it
	directory := filepath.Dir(pid.Path)
	if _, err := os.Stat(pid.Path); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0700); err != nil {
			fmt.Println("can not access specified path, writing to temp:", TempDefault(appName(pidPath)))
			Write(TempDefault(appName(pidPath)))
		}
	}

	if _, err := os.Stat(pid.Path); !os.IsNotExist(err) {
		// NOTE: Exists, checking if pid is stale
		if pid.File, err = os.OpenFile(pid.Path, os.O_RDWR|os.O_CREATE, 0600); err != nil {
			return nil, errStalePid
		} else {
			if pidData, err := ioutil.ReadFile(pid.Path); err != nil {
				pidInt, _ := strconv.Atoi(string(pidData))
				fmt.Println("pidData as string:", pidData)
				if isProcessRunning(pidInt) {
					return nil, errFileLocked
				} else {
					Clean(pid.Path)
				}
			}
		}
	}

	// NOTE: Standard creation, file locking, and return Pid file object
	var err error
	if pid.File, err = os.OpenFile(pid.Path, os.O_RDWR|os.O_CREATE, 0600); err != nil {
		return nil, errOpenFailed
	} else {
		if bytesWritten, err := pid.File.WriteString(strconv.Itoa(pid.Pid) + "\n"); err != nil {
			return nil, errWriteFailed
		} else {
			fmt.Println("bytes written:", bytesWritten)
			fmt.Println("len(pid):", pid.Pid)
			if pid.Pid != bytesWritten {
				Clean(pid.Path)
				return nil, errShortWrite
			}
		}
	}

	// NOTE: Locking via Fd() and returning the File object
	Lock(pid.File.Fd())
	return pid, nil
}

func Clean(pidPath string) error {
	if _, err := os.Stat(pidPath); os.IsNotExist(err) {
		return nil
	} else {
		// TODO: Experimental way to try to close it out with just path
		if err := removeFile(pidPath); err != nil {
			file, err := os.OpenFile(pidPath, os.O_RDWR, 0600)
			if err != nil {
				return nil
			}
			Unlock(file.Fd())
			return removeFile(pidPath)
		}
	}
	return nil
}
