package vk

// MessagesGetHistoryResponse messages.getHistory response
type MessagesGetHistoryResponse struct {
	Response struct {
		Count int `json:"count"`
		Items []struct {
			Date                  int           `json:"date"`
			FromID                int           `json:"from_id"`
			ID                    int           `json:"id"`
			Out                   int           `json:"out"`
			PeerID                int           `json:"peer_id"`
			Text                  string        `json:"text"`
			ConversationMessageID int           `json:"conversation_message_id"`
			FwdMessages           []interface{} `json:"fwd_messages"`
			Important             bool          `json:"important"`
			RandomID              int           `json:"random_id"`
			Attachments           []interface{} `json:"attachments"`
			IsHidden              bool          `json:"is_hidden"`
		} `json:"items"`
	} `json:"response"`
}
