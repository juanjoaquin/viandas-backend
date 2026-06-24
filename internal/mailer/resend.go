package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/juanjoaquin/viandas-backend/settings"
)

const resendEmailsEndpoint = "https://api.resend.com/emails"

type ResendMailer struct {
	apiKey    string
	fromEmail string
	client    *http.Client
}

func New(settings *settings.Settings) service.Mailer {
	return &ResendMailer{
		apiKey:    settings.Resend.APIKey,
		fromEmail: settings.Resend.FromEmail,
		client:    http.DefaultClient,
	}
}

func (m *ResendMailer) SendInvite(ctx context.Context, toEmail, inviteURL string) error {
	if m.apiKey == "" {
		return errors.New("resend api key is not configured")
	}
	if m.fromEmail == "" {
		return errors.New("resend from email is not configured")
	}

	payload := resendEmailRequest{
		From:    m.fromEmail,
		To:      []string{toEmail},
		Subject: "Invitación a Viandas",
		HTML: fmt.Sprintf(
			`<p>Fuiste invitado a la plataforma de Viandas.</p><p><a href="%s">Completar registro</a></p><p>Este enlace expira en 72 horas.</p>`,
			inviteURL,
		),
		Text: fmt.Sprintf("Fuiste invitado a la plataforma de Viandas. Completá tu registro en: %s", inviteURL),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEmailsEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("resend returned status %d", resp.StatusCode)
	}

	return nil
}

func (m *ResendMailer) SendPasswordReset(ctx context.Context, toEmail, resetURL string) error {
	if m.apiKey == "" {
		return errors.New("resend api key is not configured")
	}
	if m.fromEmail == "" {
		return errors.New("resend from email is not configured")
	}

	payload := resendEmailRequest{
		From:    m.fromEmail,
		To:      []string{toEmail},
		Subject: "Restablecer contraseña — Viandas",
		HTML: fmt.Sprintf(
			`<p>Recibimos una solicitud para restablecer tu contraseña en la plataforma de Viandas.</p><p><a href="%s">Restablecer contraseña</a></p><p>Este enlace expira en 1 hora. Si no solicitaste el cambio, podés ignorar este email.</p>`,
			resetURL,
		),
		Text: fmt.Sprintf("Recibimos una solicitud para restablecer tu contraseña. Ingresá al siguiente enlace (válido por 1 hora): %s", resetURL),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEmailsEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("resend returned status %d", resp.StatusCode)
	}

	return nil
}

type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	Text    string   `json:"text"`
}
