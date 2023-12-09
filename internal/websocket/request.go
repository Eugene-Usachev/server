package websocket

const (
	getOnlineUsers = uint8(iota)

	createChat
	updateChat
	deleteChat

	sendMessage
	updateMessage
	deleteMessage

	// need to set length to the router
	size
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
