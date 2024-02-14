package cache

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"zmd/util"
)

const (
	CACHE_FOLDER    = ".zmd"
	USER_CACHE_FILE = "user.json"
)

type Cache struct {
	MediaLocation string `json:"mediaLocation"`
	CourseName    string `json:"courseName"`
	CourseID      int    `json:"courseId"`
	LectureID     int    `json:"lectureId"`
	LectureName   string `json:"lectureName"`
	LecturePath   string `json:"lecturePath"`
	LectureOrder  int    `json:"lectureOrder"`
}

var (
	cacheFilePath   = ""
	cacheFolderPath = ""

	cacheFile = Cache{}
)

func Init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// setup cache file and folder path
	if len(cacheFilePath) == 0 && len(cacheFolderPath) == 0 {
		cacheFolderPath = filepath.Join(userHomeDir, CACHE_FOLDER)
		cacheFilePath = filepath.Join(cacheFolderPath, USER_CACHE_FILE)
	}

	// Make cache directory if not exists
	if _, err := os.Stat(cacheFolderPath); errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(cacheFolderPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	// Create an empty cache file it one does not exists
	if _, err := os.Stat(cacheFilePath); errors.Is(err, fs.ErrNotExist) {
		if err := writeCacheToFile(Cache{}, cacheFilePath); err != nil {
			panic(err)
		}
	}

	// Read data from cache file
	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		panic(err)
	}
	unmarshalErr := json.Unmarshal(data, &cacheFile)
	if unmarshalErr != nil {
		panic(unmarshalErr)
	}
}

func WriteAndUpdateCourse(courseName string, courseId int) {
	cacheFile.CourseName = courseName
	cacheFile.CourseID = courseId

	if err := writeCacheToFile(cacheFile, cacheFilePath); err != nil {
		panic(err)
	}
}

func WriteAndUpdateLecture(lectureId int, lectureName, lecturePath string, lectureOrder int) {
	cacheFile.LectureID = lectureId
	cacheFile.LectureName = lectureName
	cacheFile.LecturePath = lecturePath
	cacheFile.LectureOrder = lectureOrder

	if err := writeCacheToFile(cacheFile, cacheFilePath); err != nil {
		panic(err)
	}
}

func GetCache() Cache {
	return cacheFile
}

func Edit() {
	if err := util.ExecCommand("code", []string{cacheFilePath}); err != nil {
		panic(err)
	}
}

func WriteAndUpdateMediaLocation(location string) {
	cacheFile.MediaLocation = location

	if err := writeCacheToFile(cacheFile, cacheFilePath); err != nil {
		panic(err)
	}
}

func writeCacheToFile(cache interface{}, cacheFile string) error {
	data, marshallErr := json.MarshalIndent(cache, "", "    ")
	if marshallErr != nil {
		return marshallErr
	}

	writeErr := os.WriteFile(cacheFile, data, 0755)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
