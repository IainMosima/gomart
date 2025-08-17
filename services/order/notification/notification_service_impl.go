package notification

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/IainMosima/gomart/configs"
	"github.com/IainMosima/gomart/domains/auth/repository"
	"github.com/IainMosima/gomart/domains/order/schema"
	africastalking "github.com/tech-kenya/africastalkingsms"
)

type EmailNotificationServiceImpl struct {
	config   *configs.Config
	authRepo repository.AuthRepository
}

func NewNotificationServiceImpl(config *configs.Config, authRepo repository.AuthRepository) *EmailNotificationServiceImpl {
	return &EmailNotificationServiceImpl{
		config:   config,
		authRepo: authRepo,
	}
}

func (e *EmailNotificationServiceImpl) SendEmail(ctx context.Context, order *schema.OrderResponse) error {
	emailHost := e.config.EmailHost
	emailPort := e.config.EmailPort
	emailUsername := e.config.EmailUsername
	emailPassword := e.config.EmailPassword
	emailFrom := e.config.EmailFrom

	auth := smtp.PlainAuth("", emailUsername, emailPassword, emailHost)

	customer, err := e.authRepo.GetUserByID(ctx, order.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to get customer details: %w", err)
	}

	to := []string{customer.Email}
	subject := fmt.Sprintf("Order Confirmation - %s", order.OrderNumber)
	body := fmt.Sprintf(`
	Dear Customer,

	Your order has been received and is being processed.	
	Order Details:
	- Order Number: %s
	- Total Amount: KES %.2f
	- Status: %s
	- Date: %s

	Thank you for your business!

	Best regards,
	Gomart Team`,
		order.OrderNumber, order.TotalAmount, order.Status, order.CreatedAt.Format("2006-01-02 15:04:05"))

	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to[0], subject, body)

	addr := fmt.Sprintf("%s:%s", emailHost, emailPort)
	err = smtp.SendMail(addr, auth, emailFrom, to, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (e *EmailNotificationServiceImpl) SendSMS(ctx context.Context, order *schema.OrderResponse) error {
	apiKey := e.config.AfricasTalkingAPIKey
	username := e.config.AfricasTalkingUsername
	atShortCode := e.config.AfricasTalkingShortCode
	sandbox := e.config.AfricasTalkingSandbox

	if apiKey == "" || username == "" || atShortCode == "" {
		return fmt.Errorf("AfricasTalking configuration missing in environment variables")
	}

	client, err := africastalking.NewSMSClient(apiKey, username, atShortCode, sandbox)
	if err != nil {
		return fmt.Errorf("failed to create SMS client: %w", err)
	}

	message := fmt.Sprintf("Order %s confirmed! Total: $%.2f. Status: %s. Thank you for shopping with Gomart!",
		order.OrderNumber, order.TotalAmount, order.Status)

	customer, err := e.authRepo.GetUserByID(ctx, order.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to get customer details: %w", err)
	}

	customerPhone := customer.PhoneNumber

	resp, err := client.SendSMS(customerPhone, message)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	fmt.Printf("SMS sent successfully: %v\n", resp)

	return nil
}
