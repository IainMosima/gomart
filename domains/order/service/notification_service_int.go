package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/order/schema"
)

type NotificationService interface {
	SendEmail(ctx context.Context, order *schema.OrderResponse) error
	SendSMS(ctx context.Context, order *schema.OrderResponse) error
}
