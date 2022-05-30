# lf4go

### config
```golang
type Logging struct {
    RootName      string            `yaml:"root-name"` // RootName is the project-name
    RootLevel     string            `yaml:"root-level"`
    PackageLevels map[string]string `yaml:"package-levels"`
    Encoder       string            `yaml:"encoder"`
    LogFileDir    string            `yaml:"log-file-dir"`
}
```
### usage
#### config.go
```go
type config struct {
	Logging factory.Logging `yaml:"logging"`
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
var LoggerFactory = factory.NewLoggerFactory(
	func(caller string) string {
		projectName := ""
		rootName := config.Config.Logging.RootName
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
		callerPackage = callerPackage[firstSlash:lastSlash]
		return callerPackage
	},
)
```
#### actuator.go
```go
var logger = logging.LoggerFactory.NewLogger([]string{"stdout"}, []string{"stderr"})

type LogLevel struct {
	Prefix string `json:"prefix"`
	Level  string `json:"level"`
}

var body404 = "404"
var body200 = "200"

type actuatorLoggerUpdate string

func (h actuatorLoggerUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		logLevel := new(LogLevel)
		_ = json.Unmarshal(body, logLevel)
		logger.SetLevels(logLevel.Prefix, logLevel.Level)
		w.WriteHeader(200)
		_, _ = w.Write(*(*[]byte)(unsafe.Pointer(&body200)))
	} else {
		w.WriteHeader(404)
		_, _ = w.Write(*(*[]byte)(unsafe.Pointer(&body404)))
	}
}

type actuatorLoggerSelect string

func (h actuatorLoggerSelect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		logLevel := new(LogLevel)
		_ = json.Unmarshal(body, logLevel)
		outBytes, _ := json.Marshal(logger.GetLevels(logLevel.Prefix))
		w.WriteHeader(200)
		_, _ = w.Write(outBytes)
	} else {
		w.WriteHeader(404)
		_, _ = w.Write(*(*[]byte)(unsafe.Pointer(&body404)))
	}
}

func StartActuator() {
	go func() {
		http.Handle("/actuator/loggers/update", actuatorLoggerUpdate("update"))
		http.Handle("/actuator/loggers/select", actuatorLoggerSelect("select"))
		err := http.ListenAndServe("127.0.0.1:8630", nil)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
}
```