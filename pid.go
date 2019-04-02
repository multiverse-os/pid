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

func WriteToTempDirectory() (*File, error) { return Write(TempDefault()) }
func WriteToOSDefault() (*File, error)     { return Write(OSDefault()) }
func WriteToUserDefault() (*File, error)   { return Write(UserDefault()) }

///////////////////////////////////////////////////////////////////////////////
// TODO: Consider using os.Executable() as the default app name
func ValidatePath(pidPath string) string {
	if len(pidPath) < 0 || len(pidPath) > 256 {
		return OSDefault()
	}
	basename := path.Base(pidPath)
	if basename[len(basename)-1:] == "/" {
		return pidPath + serviceName() + ".pid"
	} else if filepath.Ext(basename) != ".pid" {
		return pidPath + ".pid"
	} else {
		return pidPath
	}
}

// TODO: Test this; it may not work if locked, so below fileless clean may fail
func removeFile(path string) error {
	if err := os.Remove(path); err != nil {
		return errCleanFailed
	}
	return nil
}

func serviceName() string {
	executable, _ := os.Executable()
	return path.Base(executable)
}

//[ Pid File Location Helpers ]////////////////////////////////////////////////
func OSDefault() string {
	app := serviceName()
	return ("/var/run/" + app + "/" + app + ".pid")
}
func TempDefault() string {
	return ("/var/tmp/" + serviceName() + ".pid")
}
func UserDefault() string {
	user, err := user.Current()
	if err != nil {
		return TempDefault()
	}
	return ("/run/" + user.Name + "/" + serviceName() + ".pid")
}

//[ File Methods ]/////////////////////////////////////////////////////////////
func (self *File) Clean() error {
	fmt.Println("self is:", self)
	fmt.Println("cleaning up file with pid:", self.Pid)
	fmt.Println("and path:", self.Path)
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
			Write(TempDefault())
		}
	}
	if _, err := os.Stat(pid.Path); !os.IsNotExist(err) {
		// NOTE: Exists, checking if pid is stale
		if pid.File, err = os.OpenFile(pid.Path, os.O_RDWR|os.O_CREATE, 0600); err != nil {
			return nil, errStalePid
		} else {
			if pidData, err := ioutil.ReadFile(pid.Path); err != nil {
				pidInt, _ := strconv.Atoi(string(pidData))
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
		if _, err := pid.File.WriteString(strconv.Itoa(pid.Pid)); err != nil {
			return nil, errWriteFailed
		}
	}
	// NOTE: Locking via Fd() and returning the File object
	Lock(pid.File.Fd())
	fmt.Println("pid.Pid:", pid.Pid)
	fmt.Println("pid.Path:", pid.Path)
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
