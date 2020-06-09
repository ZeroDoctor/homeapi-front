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
// 192.168.1.123 - pi
// 192.168.1.122 - desktop 
var endPoint = "http://192.168.1.122:8080/api/"

// CheckFolderContent :
func CheckFolderContent(id string) error {
	response, err := http.Get(endPoint + id)
	if err != nil {
		return err
	}
	defer response.Body.Close()

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

func DownloadFolder(file model.FileFolder) (string, string) {
	user, err := user.Current()
    if err != nil {
        return "", "could not find current user"
    }
	path := user.HomeDir + "/Downloads/" + file.FullName
	request := "download/file/" + file.ID + "." + file.Type
	if file.Dir == 1 {
		request = "download/folder/" + file.ID
		path += ".zip"
	}

	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		return "", "could not create path " + path
	}

	response, err := http.Get(endPoint + request)
	if err != nil {
		return "", "could not send request " + request
	}
	defer response.Body.Close()
	
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return "", "could not copy to " + path
	}

	return path, ""
}

func UploadFolder(parent, path string) (string, error) {
	data, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer data.Close();

	response, err := http.Post(endPoint + parent, "multipart/form-data", data)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func DeleteFolder(file model.FileFolder) (string, error) {

	request := endPoint + "file/" + file.ID + "." + file.Type
	if file.Dir == 1{
		request = endPoint + "folder/" + file.ID
	}

    req, err := http.NewRequest("DELETE", request, nil)
    if err != nil {
        return "", err
    }

    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer res.Body.Close()

    // Read Response Body
    resBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return "", err
    }

    return string(resBody), nil
}
