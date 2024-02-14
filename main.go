package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"zmd/cache"
	"zmd/cmd"
	"zmd/gui"
	"zmd/media"
	"zmd/util"
	"zmd/ztm"

	"github.com/manifoldco/promptui"
)

// folders
const (
	ZTM_ROOT_FOLDER = "/Volumes/Extreme SSD/ZTM"
)

// commands
const (
	LIST_COURSES_COMMAND   = "courses"
	LIST_LECTURES_COMMAND  = "lectures"
	SELECT_COURSE_COMMAND  = "course"
	SELECT_LECTURE_COMMAND = "lecture"
	INFO_COMMAND           = "info"
	PLAY_COMMAND           = "play"
	NEXT_COMMAND           = "next"
	GUI_COMMAND            = "gui"
)

var (
	ztmVideoFolder = "videos"
	ztmDBFolder    = "db"
	ztmDBFile      = "ztm.db"
	ztmDBFilePath  = ""
	commandManager = cmd.CommandManager{
		List: []cmd.Command{
			{
				Name: LIST_COURSES_COMMAND,
			},
			{
				Name: LIST_LECTURES_COMMAND,
			},
			{
				Name: SELECT_COURSE_COMMAND,
			},
			{
				Name: SELECT_LECTURE_COMMAND,
			},
			{
				Name: SELECT_LECTURE_COMMAND,
			},
			{
				Name: INFO_COMMAND,
			},
			{
				Name: PLAY_COMMAND,
			},
			{
				Name: NEXT_COMMAND,
			},
			{
				Name: GUI_COMMAND,
			},
		}}
)

func main() {
	initialVariableSetup()
	argLength := len(os.Args)
	switch argLength {
	case 1:
		fmt.Println("zmd is a command line application to run ztm courses")
		commandManager.Show()
	case 2:
		runCommand(os.Args[1:])
	case 3:
		runCommand(os.Args[1:])
	default:
		fmt.Println("Nothing to do....")
		return
	}
}

func runCommand(args []string) {
	if commandName, commandFound := findCommandInList(args); commandFound {
		switch commandName {
		case "courses":
			listCourses()
		case "lectures":
			listLectures()
		case "course":
			execSelectCourseCommand()
		case "lecture":
			execSelectLectureCommand()
		case "info":
			runShowCache()
		case "play":
			execPlayCommand()
		case "next":
			execNextCommand()
		case "gui":
			gui.Run()
		}
	} else {
		fmt.Printf("'%s' command not found\n", args[0])
	}
}

func execNextCommand() {
	cache.Init()

	cacheFile := cache.GetCache()

	if nextLecture, hasNext := ztm.GetNextLecture(ztmDBFilePath, cacheFile); hasNext {
		cache.WriteAndUpdateLecture(nextLecture.Id, nextLecture.Name, nextLecture.LecturePath, nextLecture.VideoOrder)

		lectureFilePath := filepath.Join(ztmVideoFolder, createLocalZTMVideoFileNameFormat(nextLecture.CourseId, nextLecture.Id))

		media.Play(lectureFilePath)
	} else {
		fmt.Println("Cannot find next lecture")
	}
}

func createLocalZTMVideoFileNameFormat(courseId, lectureId int) string {
	return fmt.Sprintf("%d.%d.mp4", courseId, lectureId)
}

func execPlayCommand() {
	cache.Init()

	cacheFile := cache.GetCache()

	lectureFilePath := filepath.Join(ztmVideoFolder, createLocalZTMVideoFileNameFormat(cacheFile.CourseID, cacheFile.LectureID))

	fmt.Printf("Playing video: \n%s -  %s\n", cacheFile.CourseName, cacheFile.LectureName)
	media.Play(lectureFilePath)
}

func runShowCache() {
	cache.Init()

	cacheFile := cache.GetCache()

	data, err := json.MarshalIndent(cacheFile, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

	if util.HasFlag("modify") {
		cache.Edit()
	}
}

func execSelectCourseCommand() {
	cache.Init()

	courses := ztm.ListCourses(ztmDBFilePath)
	courseNames := []string{}

	for _, course := range courses {
		courseNames = append(courseNames, course.Name)
	}

	prompt := promptui.Select{
		Label: "Select Course",
		Items: courseNames,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	courseName := result
	courseId := 0

	for _, course := range courses {
		if course.Name == courseName {
			courseId = course.ID
			break
		}
	}

	cache.WriteAndUpdateCourse(courseName, courseId)
}

func execSelectLectureCommand() {
	cache.Init()
	lectures := ztm.ListLectures(ztmDBFilePath)
	lectureNames := []string{}

	for _, lecture := range lectures {
		lectureNames = append(lectureNames, lecture.Name)
	}

	prompt := promptui.Select{
		Label: "Select lecture",
		Items: lectureNames,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	var selectedLecture ztm.Lecture
	for _, lecture := range lectures {
		if lecture.Name == result {
			selectedLecture = lecture
			break
		}
	}

	cache.WriteAndUpdateLecture(selectedLecture.Id, selectedLecture.Name, selectedLecture.LecturePath, selectedLecture.VideoOrder)

}

func listCourses() {
	for _, course := range ztm.ListCourses(ztmDBFilePath) {
		fmt.Println(course.Name)
	}
}

func listLectures() {
	for _, lecture := range ztm.ListLectures(ztmDBFilePath) {
		fmt.Println(lecture.Name)
	}
}

func findCommandInList(list []string) (string, bool) {
	commandDictionary := make(map[string]bool, 0)
	for _, cmd := range commandManager.List {
		commandDictionary[cmd.Name] = true
	}

	for _, cmd := range list {
		if _, exists := commandDictionary[cmd]; exists {
			return cmd, true
		}
	}
	return "", false
}

func initialVariableSetup() {
	ztmVideoFolder = filepath.Join(ZTM_ROOT_FOLDER, ztmVideoFolder)
	ztmDBFolder = filepath.Join(ZTM_ROOT_FOLDER, ztmDBFolder)
	ztmDBFilePath = filepath.Join(ztmDBFolder, ztmDBFile)
}
