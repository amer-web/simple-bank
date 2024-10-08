package mail

import (
	"github.com/amer-web/simple-bank/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmailSend(t *testing.T) {
	err := config.LoadConfig("./..")
	newSendEmail := NewGmailSender()
	err = newSendEmail.SendEmail("test", "<h1>hello</h1>", []string{"amer.khalil9094@gmail.com"}, nil)
	require.NoError(t, err)
}
