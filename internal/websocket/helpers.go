package websocket

func checkIfClientInroom(room []*Client, c *Client) bool {
	for _, Client := range room {
		if Client == c {
			return true
		}
	}
	return false
}
