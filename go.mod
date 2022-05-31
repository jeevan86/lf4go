module github.com/jeevan86/lf4go

go 1.17

require go.uber.org/zap v1.21.0
require github.com/sirupsen/logrus v1.8.1
require github.com/natefinch/lumberjack/v3 v3.0.0-alpha


require (
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
)

// 现在本地测试
replace github.com/natefinch/lumberjack/v3 v3.0.0-alpha => ../../../github.com/lumberjack
