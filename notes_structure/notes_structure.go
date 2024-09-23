package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

var (
	WorkDir_main *WorkDir = &WorkDir{}
	default_path string   = `C:\Users\arstr\OneDrive\Documents\Obsidian Vault`
	log                   = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

type Dir interface {
	FillDir()
}

// WorkDir -> Folders -> File_Note
type WorkDir struct {
	Path    string
	Notes   []File_Note
	Folders []Folder_struct
}

func (w *WorkDir) FillDir() {
	dir, err := os.ReadDir(w.Path)
	if err != nil {
		log.Error("PANIC", slog.String("ERR_MSG", err.Error()))

	}
	for _, e := range dir {
		temp_path := w.Path + "\\" + e.Name()
		fileInfo, err := os.Stat(temp_path)
		if err != nil {
			log.Error("PANIC", slog.String("ERR_MSG", err.Error()))

		}
		if fileInfo.IsDir() {
			w.Folders = append(w.Folders, Folder_struct{
				Name: e.Name(),
				Path: w.Path + "\\" + e.Name(),
			})
			//FILL THE FOLDERS STRUCT
		} else {
			w.Notes = append(w.Notes, File_Note{
				Name: e.Name(),
				Path: w.Path + "\\" + e.Name(),
			})
		}
	}
}

type Folder_struct struct {
	Name    string
	Path    string
	Notes   []File_Note
	Folders []Folder_struct
}

func (f *Folder_struct) FillDir() {
	dir, err := os.ReadDir(f.Path)
	if err != nil {
		log.Error("PANIC", slog.String("ERR_MSG", err.Error()))
	}
	for _, e := range dir {
		temp_path := f.Path + "\\" + e.Name()
		fileInfo, err := os.Stat(temp_path)
		if err != nil {
			log.Error("PANIC", slog.String("ERR_MSG", err.Error()))
		}
		if fileInfo.IsDir() {
			f.Folders = append(f.Folders, Folder_struct{
				Name: e.Name(),
				Path: f.Path + "\\" + e.Name(),
			})
			//FILL THE FOLDERS STRUCT
		} else {
			f.Notes = append(f.Notes, File_Note{
				Name: e.Name(),
				Path: f.Path + "\\" + e.Name(),
			})

		}
	}

	for i := 0; i < len(f.Folders); i++ {
		e := &f.Folders[i]
		e.FillDir()
	}

}

type File_Note struct {
	Name     string
	Path     string
	size     int
	checksum string
}

func Fill() {
	var wg sync.WaitGroup
	WorkDir_main.FillDir()
	if WorkDir_main.Folders != nil {
		for i := 0; i < len(WorkDir_main.Folders); i++ {
			e := &WorkDir_main.Folders[i]

			wg.Add(1)
			go func() {
				defer wg.Done()
				e.FillDir()
			}()
		}
	}
	wg.Wait()

}

func Init_WorkDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.ErrNotExist
	}
	if err == nil {
		WorkDir_main.Path = path
	}
	//go Fill_WorkDir()
	return nil
}
func Init_WorkDir_Default() error {

	_, err := os.Stat(default_path)
	if os.IsNotExist(err) {
		return os.ErrNotExist
	}
	if err == nil {
		WorkDir_main.Path = default_path
	}
	//go Fill_WorkDir()
	return nil
}

func main() {
	Init_WorkDir_Default()

	log.Info(
		"STRUCTER CHANGED",
		slog.String("PATH", WorkDir_main.Path),
	)
	Fill()
	fmt.Println(WorkDir_main.Folders)

	/*
		for _, e := range WorkDir_main.Folders {
			for _, e1 := range e.Notes {
				fmt.Println(e1.Name)
			}
		}
	*/
}
