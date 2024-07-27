package email

type Message interface {
	GetFrom() *string
	GetTo() []string
	GetSubject() *string
	GetText() *string
	GetHtml() *string
	GetTemplate() *Template
	GetFiles() []*File
	GetAttributes() map[string]any
}

type Template struct {
	Name       string
	Attributes map[string]any
}

func NewTemplate(name string, attributes map[string]any) *Template {
	return &Template{
		Name:       name,
		Attributes: attributes,
	}
}

func NewTemplateMessage(
	to []string,
	template *Template,
	attributes map[string]any,
) *GenericMessage {
	return &GenericMessage{
		To:         to,
		Template:   template,
		Attributes: attributes,
	}
}

func NewTemplateMessageWithAttachments(
	to []string,
	template *Template,
	files []*File,
	attributes map[string]any,
) *GenericMessage {
	return &GenericMessage{
		To:         to,
		Template:   template,
		Files:      files,
		Attributes: attributes,
	}
}

func NewGenericMessage(
	from *string,
	to []string,
	subject *string,
	text *string,
	html *string,
	template *Template,
	files []*File,
	attributes map[string]any,
) *GenericMessage {
	return &GenericMessage{
		From:       from,
		To:         to,
		Subject:    subject,
		Text:       text,
		Html:       html,
		Template:   template,
		Files:      files,
		Attributes: attributes,
	}
}

type GenericMessage struct {
	From    *string
	To      []string
	Subject *string
	Text    *string
	Html    *string

	// email templates
	Template *Template

	// adds attachments
	Files []*File

	// Attributes
	Attributes map[string]any
}

func (m *GenericMessage) GetFrom() *string {
	return m.From
}

func (m *GenericMessage) GetTo() []string {
	return m.To
}

func (m *GenericMessage) GetSubject() *string {
	return m.Subject
}

func (m *GenericMessage) GetText() *string {
	return m.Text
}

func (m *GenericMessage) GetHtml() *string {
	return m.Html
}

func (m *GenericMessage) GetTemplate() *Template {
	return m.Template
}

func (m *GenericMessage) GetFiles() []*File {
	return m.Files
}

func (m *GenericMessage) GetAttributes() map[string]any {
	return m.Attributes
}

type File struct {
	Name string
	Data []byte
}
