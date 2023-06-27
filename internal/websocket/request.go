package websocket

var (
	getOnlineUsersMethod = uint8(0)

	createChatMethod = uint8(1)
	updateChatMethod = uint8(2)
	deleteChatMethod = uint8(3)

	sendMessageMethod   = uint8(4)
	updateMessageMethod = uint8(5)
	deleteMessageMethod = uint8(6)
)

type ParsedRequest struct {
	Method uint8
	Data   []byte
	Client *Client
}

func parseRequest(request []byte, client *Client) ParsedRequest {
	var parsedRequest ParsedRequest
	methodBytes := request[0]
	parsedRequest.Method = methodBytes
	parsedRequest.Data = request[1:]
	parsedRequest.Client = client

	return parsedRequest
}
