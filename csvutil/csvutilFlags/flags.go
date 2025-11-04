package csvutilFlags

import "fmt"

type File struct {
	filename string
}

func (f *File) String() string {
	return f.filename
}

func (f *File) Set(filename string) error {
	var name string
	var extension string
	fmt.Sscanf(filename, "%s.%s", &name, &extension)
	switch extension {
	case ".csv":
		f.filename = filename
		return nil
	}
	return fmt.Errorf("Unsupported file extension: %s", extension)

}
