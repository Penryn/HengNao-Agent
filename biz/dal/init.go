package dal

import (
	"meeting_agent/biz/dal/mysql"
	"meeting_agent/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
