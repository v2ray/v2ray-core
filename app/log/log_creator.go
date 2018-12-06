package log

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
)

type HandlerCreatorOptions struct {
	Path string
}

type HandlerCreator func(LogType, HandlerCreatorOptions) (log.Handler, error)

var (
	handlerCreatorMap = make(map[LogType]HandlerCreator)
)

func RegisterHandlerCreator(logType LogType, f HandlerCreator) error {
	if f == nil {
		return newError("nil HandlerCreator")
	}

	handlerCreatorMap[logType] = f
	return nil
}

func createHandler(logType LogType, options HandlerCreatorOptions) (log.Handler, error) {
	creator, found := handlerCreatorMap[logType]
	if !found {
		return nil, newError("unable to create log handler for ", logType)
	}
	return creator(logType, options)
}

func init() {
	common.Must(RegisterHandlerCreator(LogType_Console, func(lt LogType, options HandlerCreatorOptions) (log.Handler, error) {
		return log.NewLogger(log.CreateStdoutLogWriter()), nil
	}))

	common.Must(RegisterHandlerCreator(LogType_File, func(lt LogType, options HandlerCreatorOptions) (log.Handler, error) {
		creator, err := log.CreateFileLogWriter(options.Path)
		if err != nil {
			return nil, err
		}
		return log.NewLogger(creator), nil
	}))

	common.Must(RegisterHandlerCreator(LogType_None, func(lt LogType, options HandlerCreatorOptions) (log.Handler, error) {
		return nil, nil
	}))
}
