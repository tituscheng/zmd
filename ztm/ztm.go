package ztm

import (
	"strconv"
	"zmd/cache"
	"zmd/database"
)

type Course struct {
	Name string
	ID   int
}

type Lecture struct {
	Id            int
	SectionName   string
	Name          string
	CourseId      int
	VideoOrder    int
	LecturePath   string
	LocalPath     string
	Duration      string
	FileUrl       string
	VideoUrl      string
	HasVideo      int
	HasDownloaded int
}

func ListCourses(dbPath string) []Course {
	courses := []Course{}

	sqlDB := database.SQLite{Path: dbPath}
	db := sqlDB.Connect()
	db.Table("courses").Find(&courses)
	return courses
}

func ListLectures(dbPath string) []Lecture {
	cacheFile := cache.GetCache()

	sqlDB := database.SQLite{Path: dbPath}
	db := sqlDB.Connect()
	// fmt.Printf("course id = %d\n", cacheFile.CourseID)

	lectures := []Lecture{}
	courseId := strconv.Itoa(cacheFile.CourseID)
	db.Table("lectures").Where("course_id = ? AND has_video = ?", courseId, 1).Find(&lectures)

	return lectures
}

func GetNextLecture(dbPath string, c cache.Cache) (Lecture, bool) {

	sqlDB := database.SQLite{Path: dbPath}
	db := sqlDB.Connect()

	var lecture Lecture
	result := db.Table("lectures").Order("video_order").Where("course_id = ? AND has_video = ? AND video_order > ?", c.CourseID, 1, c.LectureOrder).First(&lecture)
	if result.RowsAffected == 1 {
		return lecture, true
	}
	return lecture, false
}
