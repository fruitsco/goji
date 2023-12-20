package notification

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Action struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Url   string `json:"url"`
	Style string `json:"style"`
}

type Message struct {
	Text string  `json:"text"`
	URL  *string `json:"url"`
	// Color      *string `json:"color"`
	// PreText    *string `json:"pretext"`
	// AuthorName *string `json:"author_name"`
	// AuthorLink *string `json:"author_link"`
	// AuthorIcon *string `json:"author_icon"`
	// Title      *string `json:"title"`
	// TitleLink  *string `json:"title_link"`

	// ImageUrl   *string   `json:"image_url"`
	Fields []*Field `json:"fields"`
	// Footer     *string   `json:"footer"`
	// FooterIcon *string   `json:"footer_icon"`
	// Timestamp  *int64    `json:"ts"`
	// MarkdownIn *[]string `json:"mrkdwn_in"`
	// Actions    []*Action `json:"actions"`
}

func (attachment *Message) AddField(field Field) *Message {
	attachment.Fields = append(attachment.Fields, &field)
	return attachment
}

// func (attachment *Message) AddAction(action Action) *Message {
// 	attachment.Actions = append(attachment.Actions, &action)
// 	return attachment
// }
