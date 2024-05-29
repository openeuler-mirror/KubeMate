/*
 * Copyright 2024 KylinSoft  Co., Ltd.
 * KubeMate is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
 * PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

package log

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"time"
)

func InitLog() {
	logrus.AddHook(newLogrusHook())
	logWriter := &lumberjack.Logger{
		Filename:   "./logs/entry.log", // 日志文件路径
		MaxSize:    10,                 // 每个日志文件的最大大小（MB）
		MaxBackups: 10,                 // 保留旧日志文件的最大个数
		MaxAge:     28,                 // 保留旧日志文件的最大天数
		Compress:   true,               // 是否压缩/归档旧日志文件
		LocalTime:  true,
	}

	// 设置logrus的输出为我们创建的lumberjack.Logger实例
	logrus.SetOutput(logWriter)

	// 设置logrus的日志格式为文本格式
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.DateTime, // 设置时间格式
		FullTimestamp:   true,          // 显示完整的时间戳
	})

	// 设置logrus的日志级别
	logrus.SetLevel(logrus.InfoLevel)
}

/**
* @Description: LogrusHook
* @param 借助hook可以做参数注入
* @param
* return
*   @resp
*
 */

type LogrusHook struct {
}

func newLogrusHook() *LogrusHook {
	return &LogrusHook{}
}

func (hook *LogrusHook) Fire(entry *logrus.Entry) error {
	//entry.Data["tag"] = "tags"
	return nil
}

func (hook *LogrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
