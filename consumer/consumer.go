package consumer

import (
	"encoding/json"
	"io/ioutil"
	"io"
	"net/http"
	"os/user"
	"os"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/zerodoctor/homeapi-front/model"
	"github.com/bradfitz/slice"
)

var fileFolder = cmap.New()
var endPoint = "http://192.168.1.97:8080/api/"

// CheckFolderContent :
func CheckFolderContent(id string) error {
	_, ok := fileFolder.Get(id)
	if !ok {
		response, err := http.Get(endPoint + id)
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil
		}

		var object []model.FileFolder
		json.Unmarshal(data, &object)

		// We're just going to sneak this in here
		slice.Sort(object, func(i, j int) bool {
			return object[i].Dir > object[j].Dir
		})
		fileFolder.Set(id, object)
	}
	return nil
}

// GetFolderContent :
func GetFolderContent(id string) []model.FileFolder {
	content, ok := fileFolder.Get(id)
	if !ok {
		return nil
	}
	return content.([]model.FileFolder)
}

// GetLabelContent :
func GetLabelContent(id string) []model.Label {
	var result []model.Label
	content, ok := fileFolder.Get(id)
	if ok {
		for _, c := range content.([]model.FileFolder) {
			isFolder := true
			if c.Dir == 0 {
				isFolder = false
			}

			result = append(result, model.NewLabel(
				c.FullName, c.ID,
				isFolder, false,
				0, 0, nil,
			))
		}
	}

	return result
}

func DownloadFolder(file model.FileFolder) (string, error) {
	user, err := user.Current()
    if err != nil {
        return "", err
    }
	path := user.HomeDir + "/Downloads/" + file.FullName
	request := endPoint + "download/"
	if file.Dir == 0 {
		request += "file/" + file.ID + "." + file.Type
	} else {
		request += "folder/" + file.ID
		path += ".zip"
	}

	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		return "", err
	}

	response, err := http.Get(request)
	defer response.Body.Close()
	if err != nil {
		return "", err
	}
	
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return "", err
	}

	return path, nil
}
