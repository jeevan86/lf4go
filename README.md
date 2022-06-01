# lf4go
### usage
#### config.yml
```yaml
logging:
  factory: logrus # zap | logrus
  root-name: learngolang
  root-level: INFO
  package-levels:
    "protocol/ip/tcp": WARN
  formatter: json # normal | json
  log-file-dir: ./logs
  log-file-name: application.log
  log-to-stdout: true
  log-writer-options:
    max-file-size: 52428800 # 字节
    max-file-backups: 20
    max-file-age: 86400s
    local-time: true
    compress: true
```
#### config.go
```go
type config struct {
Logging logging `yaml:"logging"`
}

type logging struct {
Factory          string            `yaml:"factory"`
RootName         string            `yaml:"root-name"`
RootLevel        string            `yaml:"root-level"`
PackageLevels    map[string]string `yaml:"package-levels"`
Formatter        string            `yaml:"formatter"`
LogFileDir       string            `yaml:"log-file-dir"`
LogFileName      string            `yaml:"log-file-name"`
LogToStdout      bool              `yaml:"log-to-stdout"`
LogWriterOptions logWriter         `yaml:"log-writer-options"`
}

type logWriter struct {
MaxFileSize    int    `yaml:"max-file-size"`
MaxFileBackups int    `yaml:"max-file-backups"`
MaxFileAge     string `yaml:"max-file-age"` // 秒
LocalTime      bool   `yaml:"local-time"`
Compress       bool   `yaml:"compress"`
}

var configYml = "./config/config.yml"
var Config = loadConfigYml(configYml)

func loadConfigYml(path string) *config {
if len(path) == 0 {
path = configYml
}
yml, err := ioutil.ReadFile(path)
if err != nil {
fmt.Println(fmt.Sprintf("解析配置错误：%s", err.Error()))
}
var c = new(config)
err = yaml.Unmarshal(yml, c)
if err != nil {
fmt.Println(fmt.Sprintf("解析配置错误：%s", err.Error()))
return nil
}
return c
}
```
#### logging.go
```go
const EMPTY = ""
const SLASH = "/"

var logging = config.Config.Logging
var writer = logging.LogWriterOptions
var outPaths []string
var loggerFactory *factory.LoggerFactory
var mutex = sync.Mutex{}

func initLogging() {
if loggerFactory != nil {
return
}
mutex.Lock()
defer mutex.Unlock()
if loggerFactory != nil {
return
}
outPaths = make([]string, 0)
logFileDir := strings.TrimSpace(logging.LogFileDir)
if len(logFileDir) <= 0 {
logFileDir = "./logs"
}
logFileName := strings.TrimSpace(logging.LogFileName)
if len(logFileName) <= 0 {
logFileName = "./application.log"
}
logFilePath := logFileDir + SLASH + logFileName
outPaths = append(outPaths, logFilePath)
if logging.LogToStdout {
outPaths = append(outPaths, "stdout")
}

loggerFactory = factory.NewLoggerFactory(
logging.Factory,
func(caller string) string {
projectName := ""
rootName := logging.RootName
if len(rootName) > 0 {
projectName = rootName
} else {
// 会使用：go.mod文件中的module值 + "/" + 包名
myPackage := reflect.TypeOf(EMPTY).PkgPath()
projectName = myPackage[:strings.LastIndex(myPackage, "/")]
}
// 当go.mod文件中的module值，与当前项目的目录名称不一样时，这里会有问题
projectNameIdx := strings.Index(caller, projectName)
if projectNameIdx < 0 {
fmt.Println("FATAL!: 项目名称与go.mod中的module不一致，需手动配置config.yml:logging.root-name为源码项目目录的名称")
os.Exit(-1)
}
callerPackage := caller[projectNameIdx:]
firstSlash := strings.Index(callerPackage, SLASH)
lastSlash := strings.LastIndex(callerPackage, SLASH)
if firstSlash < lastSlash {
callerPackage = callerPackage[firstSlash+1 : lastSlash]
} else {
callerPackage = callerPackage[firstSlash+1:]
}
return callerPackage
},
)

}

var NewLogger = func() *factory.Logger {
_, callFilePath, _, _ := runtime.Caller(1)
initLogging()
MaxFileAge, _ := time.ParseDuration(writer.MaxFileAge)
return loggerFactory.NewLogger(callFilePath, logging.Formatter, outPaths,
writer.MaxFileSize, writer.MaxFileBackups, MaxFileAge, writer.LocalTime, writer.Compress)
}
```
#### actuator.go
```go
var logger = logging.NewLogger()

var context = "/actuator"
var loggers = context + "/loggers"

type LogLevel struct {
Prefix string `json:"prefix"`
Level  string `json:"level"`
}

var body404 = "404"
var fun404 = func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(404)
_, _ = w.Write(*(*[]byte)(unsafe.Pointer(&body404)))
}
var body200 = "200"

type HttpMethod string

const (
GET    HttpMethod = "GET"
POST   HttpMethod = "POST"
PUT    HttpMethod = "PUT"
DELETE HttpMethod = "DELETE"
PATCH  HttpMethod = "PATCH"
HEAD   HttpMethod = "HEAD"
)

type handler struct {
path   string
method HttpMethod
get    func(w http.ResponseWriter, r *http.Request)
post   func(w http.ResponseWriter, r *http.Request)
put    func(w http.ResponseWriter, r *http.Request)
delete func(w http.ResponseWriter, r *http.Request)
patch  func(w http.ResponseWriter, r *http.Request)
head   func(w http.ResponseWriter, r *http.Request)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
switch h.method {
case GET:
h.get(w, r)
break
case POST:
h.post(w, r)
break
case PUT:
h.put(w, r)
break
case DELETE:
h.delete(w, r)
break
case PATCH:
h.patch(w, r)
break
case HEAD:
h.head(w, r)
break
}
}

type HttpFunc func(w http.ResponseWriter, r *http.Request)

func trueOrDefault(b bool, f HttpFunc, def HttpFunc) HttpFunc {
if b {
return f
}
return def
}

func newHandler(path string, method HttpMethod, f HttpFunc) *handler {
return &handler{
path:   path,
method: method,
get:    trueOrDefault(method == GET, f, fun404),
post:   trueOrDefault(method == POST, f, fun404),
put:    trueOrDefault(method == PUT, f, fun404),
delete: trueOrDefault(method == DELETE, f, fun404),
patch:  trueOrDefault(method == PATCH, f, fun404),
head:   trueOrDefault(method == HEAD, f, fun404),
}
}

var actuatorLoggerUpdate = newHandler(
loggers+"/update",
POST,
func(w http.ResponseWriter, r *http.Request) {
body, _ := ioutil.ReadAll(r.Body)
logLevel := new(LogLevel)
_ = json.Unmarshal(body, logLevel)
logger.SetLevels(logLevel.Prefix, logLevel.Level)
w.WriteHeader(200)
_, _ = w.Write(*(*[]byte)(unsafe.Pointer(&body200)))
},
)

var actuatorLoggerSelect = newHandler(
loggers+"/select",
POST,
func(w http.ResponseWriter, r *http.Request) {
body, _ := ioutil.ReadAll(r.Body)
logLevel := new(LogLevel)
_ = json.Unmarshal(body, logLevel)
outBytes, _ := json.Marshal(logger.GetLevels(logLevel.Prefix))
w.WriteHeader(200)
_, _ = w.Write(outBytes)
},
)

func StartActuator() {
go func() {
http.Handle(actuatorLoggerUpdate.path, actuatorLoggerUpdate)
http.Handle(actuatorLoggerSelect.path, actuatorLoggerSelect)
err := http.ListenAndServe("127.0.0.1:8630", nil)
if err != nil {
fmt.Println(err.Error())
}
}()
}
```