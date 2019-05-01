package crawl

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"shuCourse/model/campus"
	"shuCourse/model/class"
	"shuCourse/model/course"
	"shuCourse/model/courseByTeacher"
	"shuCourse/model/courseSelectionUrl"
	"shuCourse/model/week"
	"strconv"
	"strings"
)

func loginProxy(fromurl string, username string, password string) string {
	loginUrl := os.Getenv("PROXY_AUTH_ADDRESS")
	data, _ := json.Marshal(struct {
		FromUrl  string `json:"from_url"`
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		fromurl,
		username,
		password,
	})
	response, _ := http.Post(loginUrl, "application/json", bytes.NewReader(data))
	result, _ := ioutil.ReadAll(response.Body)
	return string(result)
}

func getCoursePage(semesterId uint64, username string, password string) goquery.Document {
	url, err := courseSelectionUrl.GetBySemesterId(semesterId)
	if err != nil {
		panic(err)
	}
	token := loginProxy(url, username, password)
	client := &http.Client{}
	type requestForm struct {
		DataCount   string `json:"DataCount"`
		MinCapacity string `json:"MinCapacity"`
		MaxCapacity string `json:"MaxCapacity"`
		PageIndex   string `json:"PageIndex"`
		PageSize    string `json:"PageSize"`
	}
	form := struct {
		Url     string      `json:"url"`
		Content requestForm `json:"content"`
	}{
		url + "/StudentQuery/CtrlViewQueryCourse",
		requestForm{
			"0",
			"0",
			"100000",
			"1",
			"100000",
			//"100",
		},
	}
	data, _ := json.Marshal(form)
	request, err := http.NewRequest("POST", os.Getenv("PROXY_ADDRESS")+"post-form", bytes.NewReader(data))
	if err != nil || request == nil {
		panic(err)
	}
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	result, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		panic(err)
	}
	return *result
}

func parseWeeks(weekString string) []int64 {
	if strings.Index(weekString, "单") != -1 {
		return []int64{1, 3, 5, 7, 9}
	} else if strings.Index(weekString, "双") != -1 {
		return []int64{2, 4, 6, 8, 10}
	}
	weekString = strings.Replace(weekString, "，", ",", -1)
	discreteWeeksRegex := regexp.MustCompile(`(\d+\s*(,\s*\d+\s*)+)周`)
	discreteWeeksMatchResult := discreteWeeksRegex.FindAllStringSubmatch(weekString, -1)
	if discreteWeeksMatchResult != nil {
		weeks := strings.Split(discreteWeeksMatchResult[0][1], ",")
		var result []int64
		for _, week := range weeks {
			weekNum, _ := strconv.Atoi(week)
			result = append(result, int64(weekNum))
		}
		return result
	}
	continuousWeeksRegex := regexp.MustCompile(`(\d+)\s*-\s*(\d+)\s*周`)
	continuousMatchResult := continuousWeeksRegex.FindAllStringSubmatch(weekString, -1)
	if continuousMatchResult != nil {
		start, _ := strconv.Atoi(continuousMatchResult[0][1])
		end, _ := strconv.Atoi(continuousMatchResult[0][2])
		var result []int64
		for week := start; week <= end; week++ {
			result = append(result, int64(week))
		}
		return result
	}
	return []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
}

func parseTimes(timeString string) []class.Time {
	reg := regexp.MustCompile(`([一二三四五六日])(\d+)-(\d+)([^一二三四五]*)`)
	subMatches := reg.FindAllStringSubmatch(timeString, -1)
	var result []class.Time
	for _, subMatch := range subMatches {
		weekDay := week.Weekday[subMatch[1]]
		weeks := parseWeeks(subMatch[4])
		beginSector, _ := strconv.Atoi(subMatch[2])
		endSector, _ := strconv.Atoi(subMatch[3])
		result = append(result, class.Time{
			Weeks:       weeks,
			Weekday:     weekDay,
			BeginSector: uint8(beginSector),
			EndSector:   uint8(endSector),
		})
	}
	return result
}

func analyzeCourseColumn(courseByTeacherId uint64, column []string) {
	place := column[3]
	campusName := column[6]
	campusId := campus.Campuses[campusName]
	times := parseTimes(column[2])
	for _, time := range times {
		class.Save(class.Class{
			Time:              time,
			CourseByTeacherId: courseByTeacherId,
			CampusId:          campusId,
			Place:             place,
		})
	}
}

func analyzeCourseColumns(semesterId uint64, courseId string, courseName string, courseCredit string, columns [][]string) {
	courseCreditNumber, _ := strconv.Atoi(courseCredit)
	courseObject := course.Course{
		Id:     courseId,
		Name:   courseName,
		Credit: uint8(courseCreditNumber),
	}
	course.Save(courseObject)
	for _, column := range columns {
		courseByTeacherId, err := courseByTeacher.Save(courseByTeacher.CourseByTeacher{
			CourseId:          courseObject.Id,
			SemesterId:        semesterId,
			InCourseTeacherId: column[0],
			TeacherId:         column[1],
		})
		if err != nil {
			panic(err)
		}
		analyzeCourseColumn(courseByTeacherId, column)
	}
}

func FetchAllCourses(semesterId uint64, username string, password string) {
	doc := getCoursePage(semesterId, username, password)
	allColumns := map[string][][]string{}
	lastCourseString := ""
	doc.Find(".tbllist tr:not(:first-of-type)").Each(func(_ int, selection *goquery.Selection) {
		columns := selection.Find("td").Map(func(_ int, columnSelection *goquery.Selection) string {
			attr, exists := columnSelection.Attr("onclick")
			teacherIdRegex := regexp.MustCompile(`tid=(\d+)`)
			matchResult := teacherIdRegex.FindStringSubmatch(attr)
			if exists && len(matchResult) == 2 {
				return matchResult[1]
			}
			return strings.Trim(columnSelection.Text(), " \n")
		})
		if len(columns) == 13 {
			// 是包含课程编号和名称的标题行
			lastCourseString = columns[0] + "&" + columns[1] + "&" + columns[2]
			allColumns[lastCourseString] = [][]string{columns[3:]}
		} else {
			// 只有教师信息的行
			allColumns[lastCourseString] = append(allColumns[lastCourseString], columns)
		}
	})
	for key, value := range allColumns {
		keyStrings := strings.Split(key, "&")
		analyzeCourseColumns(semesterId, keyStrings[0], keyStrings[1], keyStrings[2], value)
	}
}
