package mailer

type MailBuilder struct {
	templateManager *TemplateManager
	mail            *Mail
}

func NewMailBuilder() *MailBuilder {
	return &MailBuilder{
		templateManager: NewTemplateManager(),
		mail:            &Mail{},
	}
}

func (mb *MailBuilder) To(email string) *MailBuilder {
	mb.mail.To = email
	return mb
}

func (mb *MailBuilder) Subject(subject string) *MailBuilder {
	mb.mail.Subject = subject
	return mb
}

func (mb *MailBuilder) PlainBody(body string) *MailBuilder {
	mb.mail.Body = body
	mb.mail.IsHTML = false
	return mb
}

func (mb *MailBuilder) HTMLBody(htmlBody string) *MailBuilder {
	mb.mail.HTMLBody = htmlBody
	mb.mail.IsHTML = true
	return mb
}

func (mb *MailBuilder) Template(templateName string, data interface{}) *MailBuilder {
	htmlContent, err := mb.templateManager.ParseTemplate(templateName, data)
	if err != nil {
		mb.mail.buildError = err
		return mb
	}

	mb.mail.HTMLBody = htmlContent
	mb.mail.IsHTML = true
	return mb
}

func (mb *MailBuilder) Build() (*Mail, error) {
	if mb.mail.buildError != nil {
		return nil, mb.mail.buildError
	}

	if err := mb.mail.validate(); err != nil {
		return nil, err
	}

	return mb.mail, nil
}
