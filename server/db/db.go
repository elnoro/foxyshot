package db

import "log"

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

// ListenOnImages waits on any paths to image descriptions from a channel and tries to save it into the database
func ListenOnImages(manipulator Manipulator, pathsAndDescriptions chan []string) {
	for pd := range pathsAndDescriptions {
		if len(pd) == 2 {
			_, err := manipulator.Add(pd[0], pd[1])
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Printf("unexpected format! got %d values while expecting 2\n", len(pd))
		}
	}
}
