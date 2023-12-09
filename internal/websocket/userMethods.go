package websocket

import (
	"github.com/goccy/go-json"
	"strconv"
)

func (handler *Handler) getOnlineUsers(request ParsedRequest) {
	var array []interface{}
	err := json.Unmarshal(request.Data, &array)
	if err != nil {
		return
	}
	var onlineUsersSlice = []int{}
	var needToGetFromRedis = []string{}
	if request.Client.userId != -1 {
		var needToSubscribe = []string{}
		clientIdStr := strconv.Itoa(request.Client.userId)
		for _, userId := range array {
			userIdI := int(userId.(float64))
			strId := strconv.Itoa(userIdI)
			needToSubscribe = append(needToSubscribe, strId)
			if handler.hub.AuthClients[userIdI] != nil {
				onlineUsersSlice = append(onlineUsersSlice, userIdI)
			} else {
				needToGetFromRedis = append(needToGetFromRedis, strId)
			}
		}
		if len(needToSubscribe) > 0 {
			err = handler.services.SubscribeOnUsers(request.Client.ctx, needToSubscribe, clientIdStr)
			if err != nil {
				onlineUsersJSON, _ := json.Marshal(map[string]interface{}{
					"data":   onlineUsersSlice,
					"method": "getOnlineUsers",
				})
				request.Client.send <- onlineUsersJSON
				return
			}
			request.Client.subscriptions = needToSubscribe
		}
	} else {
		for _, userId := range array {
			userIdI := int(userId.(float64))
			strId := strconv.Itoa(userIdI)
			if handler.hub.AuthClients[userIdI] != nil {
				onlineUsersSlice = append(onlineUsersSlice, userIdI)
			} else {
				needToGetFromRedis = append(needToGetFromRedis, strId)
			}
		}
	}
	var users []int
	if len(needToGetFromRedis) > 0 {
		users, err = handler.services.GetOnlineUsers(request.Client.ctx, needToGetFromRedis)
		if err != nil {
			onlineUsersJSON, _ := json.Marshal(map[string]interface{}{
				"data":   onlineUsersSlice,
				"method": "getOnlineUsers",
			})
			request.Client.send <- onlineUsersJSON
			return
		}
		onlineUsersSlice = append(onlineUsersSlice, users...)
	}

	onlineUsersJSON, _ := json.Marshal(map[string]interface{}{
		"data":   onlineUsersSlice,
		"method": "getOnlineUsers",
	})
	request.Client.send <- onlineUsersJSON
}
