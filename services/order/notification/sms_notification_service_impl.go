package notification

import (
	"context"

	"github.com/IainMosima/gomart/configs"
)

type SMSEmailNotificationServiceImpl struct {
	config *configs.Config
}

func NewSMSEmailNotificationServiceImpl(config *configs.Config) *SMSEmailNotificationServiceImpl {
	return &SMSEmailNotificationServiceImpl{
		config: config,
	}
}

func (S SMSEmailNotificationServiceImpl) SendOrderConfirmationSMS(ctx context.Context, customerPhone, orderNumber string, totalAmount float64) error {
	//TODO implement me
	panic("implement me")
}
