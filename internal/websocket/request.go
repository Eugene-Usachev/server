package websocket

var (
	getOnlineUsers = uint8(0)

	createChat = uint8(1)
	updateChat = uint8(2)
	deleteChat = uint8(3)

	sendMessage   = uint8(4)
	updateMessage = uint8(5)
	deleteMessage = uint8(6)
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
