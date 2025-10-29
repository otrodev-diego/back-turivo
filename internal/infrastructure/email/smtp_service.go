package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"time"

	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type SMTPService struct {
	host     string
	port     int
	username string
	password string
	from     string
	logger   *zap.Logger
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewSMTPService(config SMTPConfig, logger *zap.Logger) *SMTPService {
	// Use provided config values
	host := config.Host
	port := config.Port
	username := config.Username
	password := config.Password
	from := config.From

	// Use username as default from address if not provided
	if from == "" {
		from = username
	}

	logger.Info("📧 SMTP Service initialized",
		zap.String("host", host),
		zap.Int("port", port),
		zap.String("username", username),
		zap.String("from", from),
	)

	return &SMTPService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		logger:   logger,
	}
}

func (s *SMTPService) SendWelcomeEmail(email, name string, registrationToken string) error {
	s.logger.Info("📧 === SendWelcomeEmail Started ===",
		zap.String("email", email),
		zap.String("name", name),
		zap.String("smtp_host", s.host),
		zap.Int("smtp_port", s.port),
		zap.String("smtp_username", s.username),
	)

	// Frontend URL from environment or default
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:8080"
	}

	data := domain.WelcomeEmailData{
		Name:              name,
		Email:             email,
		RegistrationToken: registrationToken,
		RegistrationURL:   fmt.Sprintf("%s/auth/register?token=%s", frontendURL, registrationToken),
	}

	subject := "¡Bienvenido a Turivo! Completa tu registro"
	body, err := s.generateWelcomeEmailHTML(data)
	if err != nil {
		s.logger.Error("Failed to generate email HTML", zap.Error(err))
		return fmt.Errorf("failed to generate email HTML: %w", err)
	}

	return s.sendEmail(email, subject, body)
}

func (s *SMTPService) SendReservationCreated(to string, reservation *domain.Reservation, user *domain.User) error {
	s.logger.Info("📧 === SendReservationCreated Started ===",
		zap.String("email", to),
		zap.String("reservation_id", reservation.ID),
		zap.String("user_name", user.Name),
	)

	// Handle optional notes field
	notes := ""
	if reservation.Notes != nil {
		notes = *reservation.Notes
	}

	data := domain.ReservationEmailData{
		UserName:       user.Name,
		UserEmail:      user.Email,
		ReservationID:  reservation.ID,
		Pickup:         reservation.Pickup,
		Destination:    reservation.Destination,
		DateTime:       reservation.DateTime.Format("02/01/2006 15:04"),
		Passengers:     reservation.Passengers,
		VehicleType:    "", // TODO: Add vehicle type to reservation
		Amount:         reservation.Amount,
		Status:         string(reservation.Status),
		Notes:          notes,
		Stops:          0,     // TODO: Add stops count to reservation
		HasSpecialLang: false, // TODO: Add special language flag to reservation
	}

	subject := fmt.Sprintf("Confirmación de Reserva - %s", reservation.ID)
	body, err := s.generateReservationCreatedHTML(data)
	if err != nil {
		s.logger.Error("Failed to generate reservation email HTML", zap.Error(err))
		return fmt.Errorf("failed to generate reservation email HTML: %w", err)
	}

	return s.sendEmail(to, subject, body)
}

func (s *SMTPService) SendReservationNotification(to string, reservation *domain.Reservation, user *domain.User) error {
	s.logger.Info("📧 === SendReservationNotification Started ===",
		zap.String("email", to),
		zap.String("reservation_id", reservation.ID),
		zap.String("user_name", user.Name),
	)

	// Handle optional notes field
	notes := ""
	if reservation.Notes != nil {
		notes = *reservation.Notes
	}

	data := domain.ReservationEmailData{
		UserName:       user.Name,
		UserEmail:      user.Email,
		ReservationID:  reservation.ID,
		Pickup:         reservation.Pickup,
		Destination:    reservation.Destination,
		DateTime:       reservation.DateTime.Format("02/01/2006 15:04"),
		Passengers:     reservation.Passengers,
		VehicleType:    "", // TODO: Add vehicle type to reservation
		Amount:         reservation.Amount,
		Status:         string(reservation.Status),
		Notes:          notes,
		Stops:          0,     // TODO: Add stops count to reservation
		HasSpecialLang: false, // TODO: Add special language flag to reservation
	}

	subject := fmt.Sprintf("[Nueva Reserva] %s - %s", reservation.ID, user.Name)
	body, err := s.generateReservationNotificationHTML(data)
	if err != nil {
		s.logger.Error("Failed to generate reservation notification HTML", zap.Error(err))
		return fmt.Errorf("failed to generate reservation notification HTML: %w", err)
	}

	return s.sendEmail(to, subject, body)
}

func (s *SMTPService) SendSupportRequest(to string, request *domain.SupportRequest, user *domain.User) error {
	s.logger.Info("📧 === SendSupportRequest Started ===",
		zap.String("email", to),
		zap.String("user_name", user.Name),
		zap.String("user_email", user.Email),
	)

	data := domain.SupportEmailData{
		UserID:      user.ID.String(),
		UserName:    user.Name,
		UserEmail:   user.Email,
		Descripcion: request.Descripcion,
		Detalle:     request.Detalle,
		Timestamp:   time.Now().Format("02/01/2006 15:04:05"),
	}

	subject := fmt.Sprintf("[Soporte Usuario] %s", user.Email)
	body, err := s.generateSupportRequestHTML(data)
	if err != nil {
		s.logger.Error("Failed to generate support request HTML", zap.Error(err))
		return fmt.Errorf("failed to generate support request HTML: %w", err)
	}

	return s.sendEmail(to, subject, body)
}

func (s *SMTPService) sendEmail(to, subject, body string) error {
	s.logger.Info("📤 === Starting email send process ===",
		zap.String("to", to),
		zap.String("subject", subject),
		zap.String("smtp_host", s.host),
		zap.Int("smtp_port", s.port),
		zap.String("smtp_username", s.username),
		zap.String("smtp_from", s.from),
	)

	// Set up authentication information
	s.logger.Info("🔐 Setting up SMTP authentication")
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Compose message
	s.logger.Info("📝 Composing email message")
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", s.from, to, subject, body)

	s.logger.Info("📧 Message composed",
		zap.String("message_size", fmt.Sprintf("%d bytes", len(msg))),
	)

	// Send email
	smtpAddr := fmt.Sprintf("%s:%d", s.host, s.port)
	s.logger.Info("📡 Attempting to send email",
		zap.String("smtp_address", smtpAddr),
		zap.String("auth_username", s.username),
	)

	// Para puerto 465, necesitamos usar TLS directo (como en tu PHP)
	var err error
	if s.port == 465 {
		s.logger.Info("🔒 Using SSL/TLS for port 465")
		err = s.sendEmailWithTLS(smtpAddr, auth, s.from, []string{to}, []byte(msg))
	} else {
		s.logger.Info("📡 Using standard SMTP")
		err = smtp.SendMail(smtpAddr, auth, s.from, []string{to}, []byte(msg))
	}

	if err != nil {
		s.logger.Error("❌ FAILED to send email",
			zap.Error(err),
			zap.String("to", to),
			zap.String("subject", subject),
			zap.String("smtp_host", s.host),
			zap.Int("smtp_port", s.port),
			zap.String("smtp_username", s.username),
			zap.String("error_type", fmt.Sprintf("%T", err)),
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("✅ Email sent successfully",
		zap.String("to", to),
		zap.String("smtp_host", s.host),
	)
	return nil
}

// sendEmailWithTLS sends email using TLS (for port 465) like your PHP code
func (s *SMTPService) sendEmailWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	s.logger.Info("🔒 === Starting TLS email send ===")

	// Create TLS config (equivalent to your PHP SMTPOptions)
	tlsConfig := &tls.Config{
		ServerName:         s.host,
		InsecureSkipVerify: true, // Like verify_peer: false in PHP
	}

	s.logger.Info("📞 Dialing TLS connection",
		zap.String("address", addr),
		zap.String("server_name", tlsConfig.ServerName),
	)

	// Connect to the server using TLS
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		s.logger.Error("❌ Failed to dial TLS", zap.Error(err))
		return err
	}
	defer conn.Close()

	s.logger.Info("✅ TLS connection established")

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		s.logger.Error("❌ Failed to create SMTP client", zap.Error(err))
		return err
	}
	defer client.Quit()

	s.logger.Info("✅ SMTP client created")

	// Authenticate
	if auth != nil {
		s.logger.Info("🔐 Authenticating with SMTP server")
		if err = client.Auth(auth); err != nil {
			s.logger.Error("❌ SMTP authentication failed", zap.Error(err))
			return err
		}
		s.logger.Info("✅ SMTP authentication successful")
	}

	// Set sender
	s.logger.Info("📤 Setting mail sender", zap.String("from", from))
	if err = client.Mail(from); err != nil {
		s.logger.Error("❌ Failed to set sender", zap.Error(err))
		return err
	}

	// Set recipients
	for _, recipient := range to {
		s.logger.Info("📮 Adding recipient", zap.String("to", recipient))
		if err = client.Rcpt(recipient); err != nil {
			s.logger.Error("❌ Failed to add recipient", zap.Error(err), zap.String("to", recipient))
			return err
		}
	}

	// Send message
	s.logger.Info("📧 Sending message data")
	writer, err := client.Data()
	if err != nil {
		s.logger.Error("❌ Failed to get data writer", zap.Error(err))
		return err
	}

	_, err = writer.Write(msg)
	if err != nil {
		s.logger.Error("❌ Failed to write message", zap.Error(err))
		return err
	}

	err = writer.Close()
	if err != nil {
		s.logger.Error("❌ Failed to close writer", zap.Error(err))
		return err
	}

	s.logger.Info("🎉 Email sent successfully via TLS")
	return nil
}

func (s *SMTPService) generateWelcomeEmailHTML(data domain.WelcomeEmailData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Bienvenido a Turivo</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px 20px;
            text-align: center;
            border-radius: 8px 8px 0 0;
        }
        .content {
            background: #f9f9f9;
            padding: 30px 20px;
            border-radius: 0 0 8px 8px;
        }
        .button {
            display: inline-block;
            background: #667eea;
            color: white;
            padding: 12px 30px;
            text-decoration: none;
            border-radius: 5px;
            margin: 20px 0;
            font-weight: bold;
        }
        .button:hover {
            background: #5a6fd8;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            font-size: 12px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>¡Bienvenido a Turivo!</h1>
        <p>Tu plataforma de gestión de transporte</p>
    </div>
    
    <div class="content">
        <h2>Hola {{.Name}},</h2>
        
        <p>Te damos la bienvenida a <strong>Turivo</strong>, tu nueva plataforma de gestión de transporte.</p>
        
        <p>Se ha creado una cuenta para ti con el correo electrónico: <strong>{{.Email}}</strong></p>
        
        <p>Para completar tu registro y establecer tu contraseña, haz clic en el siguiente botón:</p>
        
        <div style="text-align: center;">
            <a href="{{.RegistrationURL}}" class="button">Completar Registro</a>
        </div>
        
        <p><strong>Importante:</strong> Este enlace expirará en 24 horas por seguridad.</p>
        
        <p>Si tienes alguna pregunta o necesitas ayuda, no dudes en contactarnos.</p>
        
        <p>¡Esperamos verte pronto en Turivo!</p>
        
        <p>Saludos cordiales,<br>
        <strong>El equipo de Turivo</strong></p>
    </div>
    
    <div class="footer">
        <p>Este es un mensaje automático, por favor no respondas a este correo.</p>
        <p>© 2024 Turivo. Todos los derechos reservados.</p>
    </div>
</body>
</html>
`

	t, err := template.New("welcome").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (s *SMTPService) generateReservationCreatedHTML(data domain.ReservationEmailData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Confirmación de Reserva - Turivo</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px 20px;
            text-align: center;
            border-radius: 8px 8px 0 0;
        }
        .content {
            background: #f9f9f9;
            padding: 30px 20px;
            border-radius: 0 0 8px 8px;
        }
        .reservation-details {
            background: white;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
            border-left: 4px solid #667eea;
        }
        .detail-row {
            display: flex;
            justify-content: space-between;
            margin: 10px 0;
            padding: 8px 0;
            border-bottom: 1px solid #eee;
        }
        .detail-label {
            font-weight: bold;
            color: #555;
        }
        .detail-value {
            color: #333;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            font-size: 12px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>¡Reserva Confirmada!</h1>
        <p>Tu reserva ha sido creada exitosamente</p>
    </div>
    
    <div class="content">
        <h2>Hola {{.UserName}},</h2>
        
        <p>Tu reserva ha sido creada exitosamente. A continuación encontrarás todos los detalles:</p>
        
        <div class="reservation-details">
            <h3>Detalles de la Reserva</h3>
            
            <div class="detail-row">
                <span class="detail-label">ID de Reserva:</span>
                <span class="detail-value">{{.ReservationID}}</span>
            </div>
            
            <div class="detail-row">
                <span class="detail-label">Origen:</span>
                <span class="detail-value">{{.Pickup}}</span>
            </div>
            
            <div class="detail-row">
                <span class="detail-label">Destino:</span>
                <span class="detail-value">{{.Destination}}</span>
            </div>
            
            <div class="detail-row">
                <span class="detail-label">Fecha y Hora:</span>
                <span class="detail-value">{{.DateTime}}</span>
            </div>
            
            <div class="detail-row">
                <span class="detail-label">Pasajeros:</span>
                <span class="detail-value">{{.Passengers}} personas</span>
            </div>
            
            <div class="detail-row">
                <span class="detail-label">Estado:</span>
                <span class="detail-value">{{.Status}}</span>
            </div>
            
            {{if .Amount}}
            <div class="detail-row">
                <span class="detail-label">Monto:</span>
                <span class="detail-value">${{.Amount}}</span>
            </div>
            {{end}}
            
            {{if .Notes}}
            <div class="detail-row">
                <span class="detail-label">Notas:</span>
                <span class="detail-value">{{.Notes}}</span>
            </div>
            {{end}}
        </div>
        
        <p><strong>Próximos pasos:</strong></p>
        <ul>
            <li>Recibirás actualizaciones sobre el estado de tu reserva</li>
            <li>Un conductor será asignado próximamente</li>
            <li>Te contactaremos si necesitamos información adicional</li>
        </ul>
        
        <p>Si tienes alguna pregunta o necesitas modificar tu reserva, no dudes en contactarnos.</p>
        
        <p>¡Gracias por elegir Turivo!</p>
        
        <p>Saludos cordiales,<br>
        <strong>El equipo de Turivo</strong></p>
    </div>
    
    <div class="footer">
        <p>Este es un mensaje automático, por favor no respondas a este correo.</p>
        <p>© 2024 Turivo. Todos los derechos reservados.</p>
    </div>
</body>
</html>
`

	t, err := template.New("reservation-created").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse reservation created template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute reservation created template: %w", err)
	}

	return buf.String(), nil
}

func (s *SMTPService) generateReservationNotificationHTML(data domain.ReservationEmailData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Nueva Reserva - Turivo</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
            color: white;
            padding: 30px 20px;
            text-align: center;
            border-radius: 8px 8px 0 0;
        }
        .content {
            background: #f9f9f9;
            padding: 30px 20px;
            border-radius: 0 0 8px 8px;
        }
        .reservation-details {
            background: white;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
            border-left: 4px solid #ff6b6b;
        }
        .detail-row {
            display: flex;
            justify-content: space-between;
            margin: 10px 0;
            padding: 8px 0;
            border-bottom: 1px solid #eee;
        }
        .detail-label {
            font-weight: bold;
            color: #555;
        }
        .detail-value {
            color: #333;
        }
        .user-info {
            background: #e8f4f8;
            padding: 15px;
            border-radius: 8px;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            font-size: 12px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Nueva Reserva Creada</h1>
        <p>Notificación para el equipo de operaciones</p>
    </div>
    
    <div class="content">
        <h2>Nueva reserva registrada</h2>
        
        <p>Se ha creado una nueva reserva en el sistema. Revisa los detalles a continuación:</p>
        
        <div class="user-info">
            <h3>Información del Cliente</h3>
            <div class="detail-row">
                <span class="detail-label">Nombre:</span>
                <span class="detail-value">{{.UserName}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Email:</span>
                <span class="detail-value">{{.UserEmail}}</span>
            </div>
        </div>
        
        <div class="reservation-details">
            <h3>Detalles de la Reserva</h3>
            
            <div class="detail-row">
                <span class="detail-label">ID de Reserva:</span>
                <span class="detail-value">{{.ReservationID}}</span>
            </div>
            
            <div class="detail-row">
                <span class="detail-label">Origen:</span>
                <span class="detail-value">{{.Pickup}}</span>
            </div>
            
            <div class="detail-row">
                <span class="detail-label">Destino:</span>
                <span class="detail-value">{{.Destination}}</span>
            </div>
            
            <div class="detail-row">
                <span class="detail-label">Fecha y Hora:</span>
                <span class="detail-value">{{.DateTime}}</span>
            </div>
            
            <div class="detail-row">
                <span class="detail-label">Pasajeros:</span>
                <span class="detail-value">{{.Passengers}} personas</span>
            </div>
            
            <div class="detail-row">
                <span class="detail-label">Estado:</span>
                <span class="detail-value">{{.Status}}</span>
            </div>
            
            {{if .Amount}}
            <div class="detail-row">
                <span class="detail-label">Monto:</span>
                <span class="detail-value">${{.Amount}}</span>
            </div>
            {{end}}
            
            {{if .Notes}}
            <div class="detail-row">
                <span class="detail-label">Notas:</span>
                <span class="detail-value">{{.Notes}}</span>
            </div>
            {{end}}
        </div>
        
        <p><strong>Acciones requeridas:</strong></p>
        <ul>
            <li>Revisar y confirmar la reserva</li>
            <li>Asignar conductor disponible</li>
            <li>Contactar al cliente si es necesario</li>
        </ul>
        
        <p>Accede al panel de administración para gestionar esta reserva.</p>
    </div>
    
    <div class="footer">
        <p>Sistema de notificaciones automáticas - Turivo</p>
        <p>© 2024 Turivo. Todos los derechos reservados.</p>
    </div>
</body>
</html>
`

	t, err := template.New("reservation-notification").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse reservation notification template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute reservation notification template: %w", err)
	}

	return buf.String(), nil
}

func (s *SMTPService) generateSupportRequestHTML(data domain.SupportEmailData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Solicitud de Soporte - Turivo</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: linear-gradient(135deg, #feca57 0%, #ff9ff3 100%);
            color: white;
            padding: 30px 20px;
            text-align: center;
            border-radius: 8px 8px 0 0;
        }
        .content {
            background: #f9f9f9;
            padding: 30px 20px;
            border-radius: 0 0 8px 8px;
        }
        .user-info {
            background: white;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
            border-left: 4px solid #feca57;
        }
        .support-details {
            background: #fff3cd;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
            border: 1px solid #ffeaa7;
        }
        .detail-row {
            display: flex;
            justify-content: space-between;
            margin: 10px 0;
            padding: 8px 0;
            border-bottom: 1px solid #eee;
        }
        .detail-label {
            font-weight: bold;
            color: #555;
        }
        .detail-value {
            color: #333;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            font-size: 12px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Solicitud de Soporte</h1>
        <p>Un usuario necesita asistencia</p>
    </div>
    
    <div class="content">
        <h2>Nueva solicitud de soporte</h2>
        
        <p>Se ha recibido una nueva solicitud de soporte de un usuario. Revisa los detalles a continuación:</p>
        
        <div class="user-info">
            <h3>Información del Usuario</h3>
            <div class="detail-row">
                <span class="detail-label">ID de Usuario:</span>
                <span class="detail-value">{{.UserID}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Nombre:</span>
                <span class="detail-value">{{.UserName}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Email:</span>
                <span class="detail-value">{{.UserEmail}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Fecha/Hora:</span>
                <span class="detail-value">{{.Timestamp}}</span>
            </div>
        </div>
        
        <div class="support-details">
            <h3>Detalles de la Solicitud</h3>
            
            <div style="margin: 15px 0;">
                <strong>Descripción:</strong>
                <p style="background: white; padding: 15px; border-radius: 4px; margin: 10px 0;">{{.Descripcion}}</p>
            </div>
            
            <div style="margin: 15px 0;">
                <strong>Detalle:</strong>
                <p style="background: white; padding: 15px; border-radius: 4px; margin: 10px 0; white-space: pre-wrap;">{{.Detalle}}</p>
            </div>
        </div>
        
        <p><strong>Próximos pasos:</strong></p>
        <ul>
            <li>Revisar la solicitud y priorizar según urgencia</li>
            <li>Contactar al usuario para brindar soporte</li>
            <li>Documentar la resolución en el sistema</li>
        </ul>
        
        <p>Responde directamente a este correo o contacta al usuario en: <strong>{{.UserEmail}}</strong></p>
    </div>
    
    <div class="footer">
        <p>Sistema de soporte automático - Turivo</p>
        <p>© 2024 Turivo. Todos los derechos reservados.</p>
    </div>
</body>
</html>
`

	t, err := template.New("support-request").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse support request template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute support request template: %w", err)
	}

	return buf.String(), nil
}

func (s *SMTPService) SendPasswordResetEmail(email, name, resetLink string) error {
	s.logger.Info("📧 === SendPasswordResetEmail Started ===",
		zap.String("email", email),
		zap.String("name", name),
		zap.String("smtp_host", s.host),
		zap.Int("smtp_port", s.port),
		zap.String("smtp_username", s.username),
	)

	subject := "Restablecer tu contraseña - Turivo"
	body, err := s.generatePasswordResetEmailHTML(name, resetLink)
	if err != nil {
		s.logger.Error("Failed to generate password reset email HTML", zap.Error(err))
		return fmt.Errorf("failed to generate password reset email HTML: %w", err)
	}

	return s.sendEmail(email, subject, body)
}

func (s *SMTPService) generatePasswordResetEmailHTML(name, resetLink string) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Restablecer Contraseña - Turivo</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
            color: white;
            padding: 30px 20px;
            text-align: center;
            border-radius: 8px 8px 0 0;
        }
        .content {
            background: #f9f9f9;
            padding: 30px 20px;
            border-radius: 0 0 8px 8px;
        }
        .button {
            display: inline-block;
            background: #ff6b6b;
            color: white;
            padding: 12px 30px;
            text-decoration: none;
            border-radius: 5px;
            margin: 20px 0;
            font-weight: bold;
        }
        .button:hover {
            background: #ee5a24;
        }
        .warning {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            padding: 15px;
            border-radius: 8px;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            font-size: 12px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Restablecer Contraseña</h1>
        <p>Turivo - Recuperación de acceso</p>
    </div>
    
    <div class="content">
        <h2>Hola {{.Name}},</h2>
        
        <p>Hemos recibido una solicitud para restablecer la contraseña de tu cuenta en <strong>Turivo</strong>.</p>
        
        <p>Si solicitaste este cambio, haz clic en el siguiente botón para crear una nueva contraseña:</p>
        
        <div style="text-align: center;">
            <a href="{{.ResetLink}}" class="button">Restablecer Contraseña</a>
        </div>
        
        <div class="warning">
            <p><strong>⚠️ Importante:</strong></p>
            <ul>
                <li>Este enlace expirará en 24 horas por seguridad</li>
                <li>Solo puede ser usado una vez</li>
                <li>Si no solicitaste este cambio, puedes ignorar este correo</li>
            </ul>
        </div>
        
        <p>Si el botón no funciona, puedes copiar y pegar el siguiente enlace en tu navegador:</p>
        <p style="word-break: break-all; background: #f0f0f0; padding: 10px; border-radius: 4px; font-family: monospace;">{{.ResetLink}}</p>
        
        <p>Si tienes problemas para acceder a tu cuenta o necesitas ayuda adicional, no dudes en contactarnos.</p>
        
        <p>Saludos cordiales,<br>
        <strong>El equipo de Turivo</strong></p>
    </div>
    
    <div class="footer">
        <p>Este es un mensaje automático, por favor no respondas a este correo.</p>
        <p>© 2024 Turivo. Todos los derechos reservados.</p>
    </div>
</body>
</html>
`

	data := struct {
		Name      string
		ResetLink string
	}{
		Name:      name,
		ResetLink: resetLink,
	}

	t, err := template.New("password-reset").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse password reset template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute password reset template: %w", err)
	}

	return buf.String(), nil
}
