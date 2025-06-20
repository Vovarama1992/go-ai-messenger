package dto

type AiFeedPayload struct {
	SenderEmail string `json:"senderEmail"`
	Text        string `json:"text"`
	ThreadID    string `json:"threadId"`
	BindingType string `json:"bindingType"`
}
