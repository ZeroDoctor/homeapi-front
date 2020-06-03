package model

// FileFolder :
type FileFolder struct {
	Dir      int8   `json:"fdirectory"`
	FullName string `json:"ffull_name"`
	Size     int    `json:"file_size"`
	IDName   string `json:"fname"`
	Path     string `json:"fpath"`
	Type     string `json:"ftype"`
	LastMod  string `json:"last_modified"`
	ID       string `json:"path_id"`
	Index	 int    `json:"-"`
}

// NewFileFolder :
func NewFileFolder(dir int8, size int, fullName, idName, path, ftype, lastMod, id string, index int) FileFolder {
	return FileFolder{
		Dir:      dir,
		FullName: fullName,
		Size:     size,
		IDName:   idName,
		Path:     path,
		Type:     ftype,
		LastMod:  lastMod,
		ID:       id,
		Index:    index,
	}
}

// EmptyFileFolder :
func EmptyFileFolder() FileFolder {
	return FileFolder{
		Dir:      0,
		FullName: "",
		Size:     0,
		IDName:   "",
		Path:     "",
		Type:     "",
		LastMod:  "",
		ID:       "$EMPTY",
		Index:    0,
	}
}
