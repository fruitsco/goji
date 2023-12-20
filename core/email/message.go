package email

type Message interface {
	From() *string
	To() []string
	Subject() *string
	Text() *string
	HTML() *string
	Template() *Template
	Files() []*File
	Attributes() map[string]any
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
	from *string,
	to []string,
	template *Template,
	attributes map[string]any,
) *GenericMessage {
	return &GenericMessage{
		from:       from,
		to:         to,
		template:   template,
		attributes: attributes,
	}
}

func NewTemplateMessageWithAttachments(
	from *string,
	to []string,
	template *Template,
	files []*File,
	attributes map[string]any,
) *GenericMessage {
	return &GenericMessage{
		from:       from,
		to:         to,
		template:   template,
		files:      files,
		attributes: attributes,
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
		from:       from,
		to:         to,
		subject:    subject,
		text:       text,
		html:       html,
		template:   template,
		files:      files,
		attributes: attributes,
	}
}

type GenericMessage struct {
	from    *string
	to      []string
	subject *string
	text    *string
	html    *string

	// email templates
	template *Template

	// adds attachments
	files []*File

	// attributes
	attributes map[string]any
}

func (m *GenericMessage) From() *string {
	return m.from
}

func (m *GenericMessage) To() []string {
	return m.to
}

func (m *GenericMessage) Subject() *string {
	return m.subject
}

func (m *GenericMessage) Text() *string {
	return m.text
}

func (m *GenericMessage) HTML() *string {
	return m.html
}

func (m *GenericMessage) Template() *Template {
	return m.template
}

func (m *GenericMessage) Files() []*File {
	return m.files
}

func (m *GenericMessage) Attributes() map[string]any {
	return m.attributes
}

type File struct {
	Name string
	Data []byte
}
