package data

import (
	"encoding/json"
	"fmt"
	"time"
	"wehcat-bot-go/internal/ai"
	"wehcat-bot-go/internal/app"
)

func getUserContextKey(uid string) string {
	return fmt.Sprintf("user_context:%s", uid)
}

func GetUserContext(uid string) (data []ai.Message, err error) {
	str, err := app.Rdb.Get(getUserContextKey(uid)).Result()
	if err != nil {
		return
	}
	if str == "" {
		return
	}
	err = json.Unmarshal([]byte(str), &data)
	return data, err
}

func SetUserContext(uid string, context []ai.Message) (err error) {
	key := getUserContextKey(uid)
	bytes, err := json.Marshal(context)
	if err != nil {
		return
	}
	return app.Rdb.Set(key, string(bytes), time.Hour*3).Err()
}
