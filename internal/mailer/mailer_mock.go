package mailer

type mockSmtpMailer struct{}

func (m *mockSmtpMailer) SendActivationEmail(e, u, h string) (err error) {
	// do nothing just log
	logEmailSent(e, "📧 Activation Email Sent !!!", err)
	return
}

func (m *mockSmtpMailer) SendConfirmationEmail(e string) (err error) {
	// do nothing just log
	logEmailSent(e, "📧 Confirmation Email Sent !!!", err)
	return
}

var MockSmtpMailer = mockSmtpMailer{}
