package errx

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// StackError 是自定义错误类型
type StackError struct {
	msg   string
	cause error
	stack []uintptr
}

// 捕获调用堆栈（只捕获一次）
func captureStack() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:]) // 跳过 runtime + captureStack + New
	return pcs[:n]
}

// New 创建带堆栈的新错误
func New(msg string) error {
	return &StackError{
		msg:   msg,
		stack: captureStack(),
	}
}

// Wrap 封装已有错误（不重复捕获堆栈）
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	if se, ok := err.(*StackError); ok {
		return &StackError{
			msg:   msg,
			cause: se,
			stack: se.stack,
		}
	}
	return &StackError{
		msg:   msg,
		cause: err,
		stack: captureStack(),
	}
}

func (e *StackError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.cause)
	}
	return e.msg
}

func (e *StackError) Unwrap() error { return e.cause }

// StackTrace 打印堆栈信息
func (e *StackError) StackTrace() string {
	var sb strings.Builder
	for _, pc := range e.stack {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pc)
		sb.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", fn.Name(), filepath.ToSlash(file), line))
	}
	return sb.String()
}

// Print 打印错误和堆栈（如果有）
func Print(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)

	var se *StackError
	if errors.As(err, &se) {
		fmt.Println("\nStack trace:")
		fmt.Println(se.StackTrace())
	}
}

// Format 返回完整堆栈字符串（适合日志）
func Format(err error) string {
	if err == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(err.Error())

	var se *StackError
	if errors.As(err, &se) {
		sb.WriteString("\n\nStack trace:\n")
		sb.WriteString(se.StackTrace())
	}
	return sb.String()
}
