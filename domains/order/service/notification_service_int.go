package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/order/schema"
)

type EmailNotificationService interface {
	SendOrderNotificationEmail(ctx context.Context, order *schema.OrderResponse) error
}

type SMSNotificationService interface {
	SendOrderConfirmationSMS(ctx context.Context, customerPhone, orderNumber string, totalAmount float64) error
}
