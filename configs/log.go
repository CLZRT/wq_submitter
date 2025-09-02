package configs

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strings"
)

//时间,操作，日志等级

type logFormatter struct {
	logrus.TextFormatter
}

func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	prettyCaller := func(frame *runtime.Frame) string {
		_, fileName := filepath.Split(frame.File)
		return fmt.Sprintf("%s:%d", fileName, frame.Line)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(fmt.Sprintf("[%s] %s", entry.Time.Format(f.TimestampFormat), strings.ToUpper(entry.Level.String()))) //打印日志时间,等级
	if entry.HasCaller() {
		b.WriteString(fmt.Sprintf("[%s]", prettyCaller(entry.Caller))) //打印文件名，文件行数
	}
	b.WriteString(fmt.Sprintf(" [%s]\n", entry.Message)) //日志内容

	return b.Bytes(), nil
}
