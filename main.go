package main

import (
	"os"
	"fmt"
	"strings"
	"errors"
	"io"
)

var buffer = make([]byte, 512*1024)

func main() {
	var err error

	workDir, err := GetWorkDir()
	if err != nil {
		fmt.Println(err)
		return
	}
	files, folders, err := GetFilesFromWorkDir(workDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = DirChecker(files, folders, workDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = MoveFiles(files, workDir)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetWorkDir() (string, error) {
	var err error
	workDir := ""
	if len(os.Args) == 1 {
		workDir = "./"
	} else if len(os.Args) == 2 {
		info, err := os.Stat(os.Args[1])
		if err != nil {
			err = errors.New("второй аргумент не является существующей директорией")
			return workDir, err
		}
		if info.IsDir() {
			workDir = os.Args[1] + "/"
		} else {
			err = errors.New("второй аргумент не является существующей директорией")
			return workDir, err
		}
	}
	if workDir == "" {
		err = errors.New("некорректный ввод")
		return workDir, err
	}
	return workDir, nil
}

func GetFilesFromWorkDir(workDir string) ([]string, []string, error) {
	var err error
	var files []string
	var folders []string
	entries, err := os.ReadDir(workDir)
	if err != nil {
		err = errors.New("ошибка при чтении каталога")
		return nil, nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		} else if (entry.Name() != "go.mod") && (entry.Name() != "main.go") {
			files = append(files, entry.Name())
		}
	}
	return files, folders, nil
}

func DirChecker(files []string, folders []string, workDir string) (error) {
	var err error
	var extensions []string
	isEmptyExtension := false
	for _, value := range files {
		dotIndex := strings.Index(value, ".")
		if dotIndex != -1 {
			extensions = append(extensions, value[dotIndex + 1:])
		} else {
			isEmptyExtension = true
		}
	}
	err = CreateDir(extensions, folders, isEmptyExtension, workDir)
	if err != nil {
		return err
	}
	return nil
}

func CreateDir(extensions []string, folders []string, isEmptyExtension bool, workDir string) (error) {
	var err error
	flag := false
	if isEmptyExtension {
		for _, value := range folders {
			if value == "for_empty_extension_folder" {
				flag = true
				break
			}
		}
		if !flag {
			err = os.Mkdir(workDir + "for_empty_extension_folder", 0755)
			if err != nil {
				err = errors.New("ошибка при создании директории")
				return err
			}
			folders = append(folders, "for_empty_extension_folder")
		}
		flag = false
	}
	for _, value := range extensions {
		for _, dirValue := range folders {
			if value + "_folder" == dirValue {
				flag = true
				break
			}
		}
		if !flag {
			err = os.Mkdir(workDir + value + "_folder", 0755)
			if err != nil {
				err = errors.New("ошибка при создании директории")
				return err
			}
			folders = append(folders, value)
		}
		flag = false
	}
	return nil
}

func MoveFiles(files []string, workDir string) error {
	var err error
	for _, value := range files {
		err = CopyFiles(value, workDir)
		if err != nil {
			return err
		}
	}
	err = DeleteFiles(files, workDir)
	if err != nil {
		return err
	}
	return nil
}

func CopyFiles(file string, workDir string) error {
	var err error
	var extension string
	dotIndex := strings.Index(file, ".")
	if dotIndex != -1 {
		extension = file[dotIndex+1:]
	} else {
		extension = "for_empty_extension"
	}
	filePathDst := workDir + extension + "_folder" + "/" + file
	filePathSrc := workDir + file
	srcFile, err := os.Open(filePathSrc)
	if err != nil {
		err = errors.New("невозможно открыть файл")
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(filePathDst)
	if err != nil {
		err = errors.New("невозможно создать файл")
		return err
	}
	defer dstFile.Close()
	for {
		n, err := srcFile.Read(buffer)
		if err != nil && err != io.EOF {
			err = errors.New("ошибка при чтении файла")
			return err
		}
		if n == 0 {
			break
		}
		_, err = dstFile.Write(buffer[:n])
		if err != nil {
			err = errors.New("ошибка при записи файла")
			return err
		}
	}
	return nil
}

func DeleteFiles(files []string, workDir string) error {
	var err error
	for _, value := range files {
		filePath := workDir + "/" + value
		err = os.Remove(filePath)
		if err != nil {
			err = errors.New("не удалось удалить файл")
			return err
		}
	}
	return nil
}
