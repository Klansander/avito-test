package errors

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// ErrorApp структура для описания ошибкок в приложении
type ErrorApp struct {
	Code    int
	Message string
}

func (e *ErrorApp) Error() string {

	return e.Message

}

func New(code int, message string) error {

	err := &ErrorApp{Code: code, Message: message}

	return errors.Wrap(err, wrapPosition(2))

}

func Wrap(err error) error {

	return errors.Wrap(err, wrapPosition(2))

}

func Cause(err error) error {

	return errors.Cause(err)

}

func wrapPosition(level int) string {

	f := getFrame(level)
	_, filename := path.Split(f.File)
	dir := strings.Split(f.Function, ".")[0]

	return fmt.Sprintf("--> %s/%s/%d", dir, filename, f.Line)

}

func getFrame(skipFrames int) runtime.Frame {

	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame

}

func As(err error, target interface{}) bool { return errors.As(err, target) }
func Is(err error, target error) bool       { return errors.Is(err, target) }
