package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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
			File := &File_Note{Name: e.Name(), Path: w.Path + "\\" + e.Name()}
			writeFileSize(File)
			writeCheckSum(File)
			w.Notes = append(w.Notes, *File)
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
			File := &File_Note{Name: e.Name(), Path: f.Path + "\\" + e.Name()}
			writeFileSize(File)
			writeCheckSum(File)
			f.Notes = append(f.Notes, *File)
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
	//fmt.Println(WorkDir_main.Folders[0].Notes[5].checksum)
}
func writeCheckSum(f *File_Note) {
	file, err := os.Open(f.Path)
	if err != nil {
		log.Error(err.Error())
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Error(err.Error())
	}
	checksum_hex := hex.EncodeToString(h.Sum(nil))
	f.checksum = string(checksum_hex)

}
func writeFileSize(f *File_Note) {
	file_stat, err := os.Stat(f.Path)
	if err != nil {
		log.Error(err.Error())
	}
	f.size = int(file_stat.Size())
}
