package mailer

type MockSmtpMailer struct{}

func (m *MockSmtpMailer) SendActivationEmail(e, u, h string) (err error) {
	// do nothing just log
	logEmailSent(e, "Activation Email Sent !!!", err)
	return
}

func (m *MockSmtpMailer) SendConfirmationEmail(e string) (err error) {
	// do nothing just log
	logEmailSent(e, "Confirmation Email Sent !!!", err)
	return
}
