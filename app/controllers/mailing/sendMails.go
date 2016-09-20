package mailing

import (
	"github.com/mailgun/mailgun-go"
)

var (
	apiKey string = "key-d43fcd11e3fa56d5ccc0d8f4e241aff4"
	domain string = "sandboxc9c4b7ab94f84710adb8e847cff7ac43.mailgun.org"
)

func PassRecoverMail(destination string) (string, error) {
  mg := mailgun.NewMailgun(domain, apiKey, "")
  m := mg.NewMessage(
    "Mr-Tooth <mr-tooth@gmail.com>",	// Correo por definir
    "Link para reestablecer contraseña",
    "No responder este correo",
    destination ,
  )
  m.SetHtml("<html><body><h1>Link para reestablecer contraseña</h1><a href='#'>Click aquí</a></body></html>") // Link por definir
  _, id, err := mg.Send(m)
  return id, err
}
