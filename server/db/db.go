package db

type ScreenshotsDb interface {
	Manipulator
	Finder
	Lister
}

type Manipulator interface {
	Add(path, desc string) (string, error)
	Remove(id string) error
}

type Finder interface {
	FindByDesc(desc string) ([]ImageRecord, error)
}

type Lister interface {
	All() ([]ImageRecord, error)
}

type ImageRecord struct {
	Id   string
	Data *ImageData
}

type ImageData struct {
	Path string
	Desc string
}
