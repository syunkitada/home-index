package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	Path      string
	Text      string
	UpdatedAt time.Time
}

func dirwalk(root string, dir string) (fileInfos []FileInfo) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			fileInfos = append(fileInfos, dirwalk(root, filepath.Join(dir, file.Name()))...)
			continue
		} else {
			var path string
			if dir != "." {
				path = filepath.Join(dir, file.Name())
			} else {
				path = file.Name()
			}
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("Failed ioutil.ReadFile: err=%s", err.Error())
			}
			contentType := http.DetectContentType(content)
			if contentType == "text/plain; charset=utf-8" {
				path = strings.Replace(path, root, "", 1)
				fileInfo := FileInfo{
					Path:      path,
					Text:      string(content),
					UpdatedAt: file.ModTime(),
				}
				fileInfos = append(fileInfos, fileInfo)
			}
		}
	}
	return
}

type Text struct {
	Text      string
	UpdatedAt time.Time
}

type IndexData struct {
	Entry                 int
	TokenDocumentIndexMap map[string]map[int][]int // [token]: [documentId][tokenIndex]
	IdTextMap             []Text
	IdPathMap             []string
	PathMap               map[string]int
}

type Config struct {
	Pages []Page
}

type Page struct {
	Name  string
	Root  string
	Entry string
}

func GeneratePageIndex(page *Page) {
	fileInfos := dirwalk(page.Root, page.Root)
	idTextMap := []Text{}
	idPathMap := []string{}
	pathMap := map[string]int{}
	var entry int
	for i, file := range fileInfos {
		idTextMap = append(idTextMap, Text{Text: file.Text, UpdatedAt: file.UpdatedAt})
		idPathMap = append(idPathMap, file.Path)
		pathMap[file.Path] = i
		if page.Entry == "" {
			if file.Path == "README.md" {
				entry = i
			}
		} else if file.Path == page.Entry {
			entry = i
		}
	}

	indexData := IndexData{
		Entry:     entry,
		IdTextMap: idTextMap,
		IdPathMap: idPathMap,
		PathMap:   pathMap,
	}

	bytes, err := json.Marshal(indexData)
	if err != nil {
		log.Fatalf("Failed json.Marshal: err=%s", err.Error())
	}

	if err := ioutil.WriteFile(fmt.Sprintf("./%s.json", page.Name), bytes, 0644); err != nil {
		log.Fatalf("Failed WriteFile index.json: err=%s", err.Error())
	}
}

func main() {
	var conf Config
	if bytes, err := ioutil.ReadFile("./config.json"); err != nil {
		log.Fatalf("Failed json.Marshal: %v", err)
	} else {
		if err := json.Unmarshal(bytes, &conf); err != nil {
			log.Fatalf("Failed json.Marshal: %v", err)
		}
	}
	for i := range conf.Pages {
		GeneratePageIndex(&conf.Pages[i])
	}
	return
}
