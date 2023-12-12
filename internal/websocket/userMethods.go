package websocket

import (
	fb "github.com/Eugene-Usachev/fastbytes"
	"github.com/goccy/go-json"
)

func (handler *Handler) getOnlineUsers(request ParsedRequest) {
	// Here I commenter the code for scaling. I didn't test. Uncomment it later if you need.
	var array []interface{}
	err := json.Unmarshal(request.Data, &array)
	if err != nil {
		return
	}
	var onlineUsersSlice = []int{}
	//var needToGetFromRedis = []string{}
	if request.Client.userId != -1 {
		var needToSubscribe = []string{}
		//clientIdStr := fb.B2S(fb.I2B(request.Client.userId))
		for _, userId := range array {
			userIdI := int(userId.(float64))
			strId := fb.B2S(fb.I2B(userIdI))
			needToSubscribe = append(needToSubscribe, strId)
			if handler.hub.AuthClients[userIdI] != nil {
				onlineUsersSlice = append(onlineUsersSlice, userIdI)
			}
			// else {
			//	needToGetFromRedis = append(needToGetFromRedis, strId)
			//}
		}
		if len(needToSubscribe) > 0 {
			//err = handler.services.SubscribeOnUsers(request.Client.ctx, needToSubscribe, clientIdStr)
			//if err != nil {
			//	onlineUsersJSON, _ := createResponse("getOnlineUsers", onlineUsersSlice)
			//	request.Client.send <- onlineUsersJSON
			//	return
			//}
			for _, userId := range needToSubscribe {
				handler.hub.subscribeToClient(request.Client, userId)
			}
		}
	} else {
		for _, userId := range array {
			userIdI := int(userId.(float64))
			//strId := fb.B2S(fb.I2B(userIdI))
			if handler.hub.AuthClients[userIdI] != nil {
				onlineUsersSlice = append(onlineUsersSlice, userIdI)
			}
			//else {
			//	needToGetFromRedis = append(needToGetFromRedis, strId)
			//}
		}
	}

	//var users []int
	//if len(needToGetFromRedis) > 0 {
	//	users, err = handler.services.GetOnlineUsers(request.Client.ctx, needToGetFromRedis)
	//	if err != nil {
	//		onlineUsersJSON, _ := createResponse("getOnlineUsers", onlineUsersSlice)
	//		handler.hub.logger.Error("WS: getOnlineUsers service error:", err)
	//		request.Client.send <- onlineUsersJSON
	//		return
	//	}
	//	onlineUsersSlice = append(onlineUsersSlice, users...)
	//}

	onlineUsersJSON, _ := createResponse("getOnlineUsers", onlineUsersSlice)
	request.Client.send <- onlineUsersJSON
}

func processSubscriptionOnOnline(sub *subscription, clientId int) {
	res, err := createResponse("userOnline", clientId)
	if err != nil {
		return
	}
	sub.m.RLock()
	defer func() {
		sub.m.RUnlock()
	}()
	for _, c := range sub.slice {
		c.send <- res
	}
}

func processSubscriptionOnOffline(sub *subscription, clientId int) {
	res, err := createResponse("userOffline", clientId)
	if err != nil {
		return
	}
	sub.m.RLock()
	defer func() {
		sub.m.RUnlock()
	}()
	for _, c := range sub.slice {
		c.send <- res
	}
}
