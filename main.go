package main

import (
	"os"
	"fmt"
	"strings"
	"errors"
	"io"
)

const (
    ErrInvalidInput    = "некорректный ввод"
    ErrNotADirectory   = "второй аргумент не является существующей директорией"
    ErrReadDir         = "ошибка при чтении каталога"
    ErrCreateDir       = "ошибка при создании директории"
    ErrOpenFile        = "невозможно открыть файл"
    ErrCreateFile      = "невозможно создать файл"
    ErrDeleteFile      = "не удалось удалить файл"
    ErrReadFile        = "ошибка при чтении файла"
    ErrWriteFile       = "ошибка при записи файла"
)

var buffer = make([]byte, 512*1024)

type fileInfo struct {
	filename string
	extension string
}

type dirInfo struct {
	files 	[]fileInfo
	folders map[string]struct{}
	workDir string
}

func main() {
	var objects dirInfo
	var err error

	objects, err = GetWorkDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = GetObjectsFromWorkDir(&objects)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = GetExtensions(&objects)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = CreateDir(&objects)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = CopyAndDeleteFiles(objects)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetWorkDir() (objects dirInfo, err error) {
	objects.workDir = ""
	objects.folders = make(map[string]struct{})

	if len(os.Args) == 1 {
		objects.workDir = "./"
		return
	}
	
	if len(os.Args) != 2 {
		return objects, errors.New(ErrInvalidInput)
	}
	
	info, err := os.Stat(os.Args[1])
	if err != nil {
		return objects, errors.New(ErrNotADirectory)
	}

	if info.IsDir() {
		objects.workDir = os.Args[1] + "/"
	} else {
		err = errors.New(ErrNotADirectory)
	}
	return
}

func GetObjectsFromWorkDir(objects *dirInfo) (err error) {
	entries, err := os.ReadDir(objects.workDir)
	if err != nil {
		return errors.New(ErrReadDir)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			objects.folders[entry.Name()] = struct{}{}
		} else if (entry.Name() != "go.mod") && (entry.Name() != "main.go") {
			file := fileInfo{filename: entry.Name(), extension: ""}
			objects.files = append(objects.files, file)
		}
	}

	return
}

func GetExtensions(objects *dirInfo) (err error) {
	for i := range objects.files {
		dotIndex := strings.Index(objects.files[i].filename, ".")
		if dotIndex != -1 {
			objects.files[i].extension = objects.files[i].filename[dotIndex + 1:] + "_folder"
		} else {
			objects.files[i].extension = "for_empty_extension_folder"
		}
	}
	return
}

func CreateDir(objects *dirInfo) (err error) {
	for _, value := range objects.files {
		if _, exists := objects.folders[value.extension]; exists {
			continue
		} else {
			err = os.Mkdir(objects.workDir + value.extension, 0755)
			if err != nil {
				return errors.New(ErrCreateDir)
			}
			objects.folders[value.extension] = struct{}{}
		}
	}
	return
}

func CopyAndDeleteFiles(objects dirInfo) (err error) {
	for _, value := range objects.files {
		filePathDst := objects.workDir + value.extension + "/" + value.filename
		filePathSrc := objects.workDir + value.filename

		err = CopyFiles(filePathDst, filePathSrc)
		if err != nil {
			return
		}

		err = os.Remove(filePathSrc)
		if err != nil {
			return errors.New(ErrDeleteFile)
		}
	}
	return
}

func CopyFiles(filePathDst string, filePathSrc string) (err error) {
	srcFile, err := os.Open(filePathSrc)
	if err != nil {
		return errors.New(ErrOpenFile)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(filePathDst)
	if err != nil {
		return errors.New(ErrCreateFile)
	}
	defer dstFile.Close()

	for {
		n, err := srcFile.Read(buffer)
		if err != nil && err != io.EOF {
			return errors.New(ErrReadFile)
		}
		if n == 0 {
			return nil
		}
		_, err = dstFile.Write(buffer[:n])
		if err != nil {
			return errors.New(ErrWriteFile)
		}
	}
}
