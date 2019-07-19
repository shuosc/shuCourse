# API Reference

## 模型

### 课程

```json
{
   "id": 	   课程id
   "name":     课程名称
   "credit":   课程学分数
}
```

### 由某一老师授课的课程

```json
{
   "id": 	                课程id
   "course_id":             所属课程id
   "semester_id":           开课所在学期的id,
   "in_course_teacher_id":  选课系统中的教师号
}
```

### 某一节课
```json
{
  "weeks":                哪几周有这节课,
  "weekday":              这节课在周几,
  "begin_sector":         从这一天第几课时开始,
  "end_sector":           上到第几节课时（含）,
  "course_by_teacher_id": 所属“某一老师授课的课程”的id,
  "campus_id":            所在校区id,
  "place":                教室/地点
}
```

### 标准查询输出格式
即在"由某一老师授课的课程"的基础上：
- 将`course_id`展开成`course`的全部信息
- `classes`包含这一课程的所有课

## cli

### addAdmin
这一功能用于新增管理员
```shell
./cli addAdmin 17120238
```
会返回对应JWT。

### setSemesterUrl
用于设置某一个学期对应的选课URL。
```shell
./cli setSemesterUrl 8 http://xk.autoisp.shu.edu.cn
```

### fetchCourses
开始抓取某学期的所有课程数据。
```shell
./cli fetchCourses 8 17120238 [我的学生证密码]
```


## web api

- `GET /ping`

  检查服务是否可用，应该直接返回`pong`。

- `GET /semester?id=[一个id]`

  返回id对应的"标准查询输出格式数据"。

- `GET /semester?semester_id=[一个id]&course_id=[一个id]&in_course_teacher_id=[一个id]`

  用于直接从选课网站上得到的信息查询对应的"标准查询输出格式数据"。