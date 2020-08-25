package log_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	clog "v2ray.com/core/common/log"
	"v2ray.com/core/testing/mocks"
)

func TestCustomLogHandler(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	var loggedValue []string

	mockHandler := mocks.NewLogHandler(mockCtl)
	mockHandler.EXPECT().Handle(gomock.Any()).AnyTimes().DoAndReturn(func(msg clog.Message) {
		loggedValue = append(loggedValue, msg.String())
	})

	log.RegisterHandlerCreator(log.LogType_Console, func(lt log.LogType, options log.HandlerCreatorOptions) (clog.Handler, error) {
		return mockHandler, nil
	})

	logger, err := log.New(context.Background(), &log.Config{
		ErrorLogLevel: clog.Severity_Debug,
		ErrorLogType:  log.LogType_Console,
		AccessLogType: log.LogType_None,
	})
	common.Must(err)

	common.Must(logger.Start())

	clog.Record(&clog.GeneralMessage{
		Severity: clog.Severity_Debug,
		Content:  "test",
	})

	if len(loggedValue) < 2 {
		t.Fatal("expected 2 log messages, but actually ", loggedValue)
	}

	if loggedValue[1] != "[Debug] test" {
		t.Fatal("expected '[Debug] test', but actually ", loggedValue[1])
	}

	common.Must(logger.Close())
}
