# 部署
本项目已经打包成[docker镜像](https://hub.docker.com/r/shuosc/shu-course/)。
## 支持服务
### postgresql数据库
migration文件位于本repo的 [migration](https://github.com/shuosc/shuCourse/tree/master/migration) 目录中。

建议使用 [golang-migrate](https://github.com/golang-migrate/migrate) 来进行 migrate。
```shell
migrate -source github://[你的Github用户名]:[你的Github Access Token]@shuosc/shuCourse/migration -database [你的postgrsql数据库url] up
```

## 服务本身
### 环境变量
- `PORT`: 服务端口
- `DB_ADDRESS`: 数据库url
- `JWT_SECRET`: jwt密钥
- `PROXY_ADDRESS`: 访问学校选课网站代理服务地址
- `PROXY_AUTH_ADDRESS`: 登录学校选课网站代理服务地址

### k8s
在k8s下使用如下yaml，假设
- `JWT_SECRET`由k8s secret给出
- 数据库服务器在`shu-course-postgres-svc`
- 代理服务地址在`shu-course-proxy-svc`
- 代理服务登录地址在`shu-course-proxy-svc/login`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: shu-course
spec:
  selector:
    matchLabels:
      run: shu-course
  replicas: 1
  template:
    metadata:
      labels:
        run: shu-course
    spec:
      containers:
      - name: shu-course
        image: shuosc/shu-course
        env:
        - name: PORT
          value: "8000"
        - name: DB_ADDRESS
          value: "postgres://shuosc@shu-course-postgres-svc:5432/shu-course?sslmode=disable"
        - name: PROXY_ADDRESS
          value: "http://shu-course-proxy-svc/"
        - name: PROXY_AUTH_ADDRESS
          value: "http://shu-course-proxy-svc/login"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: shuosc-secret
              key: JWT_SECRET
        ports:
        - containerPort: 8000
---
apiVersion: v1
kind: Service
metadata:
  name: shu-course-svc
spec:
  selector:
     run: shu-course
  ports:
  - protocol: TCP
    port: 8000
    targetPort: 8000
```