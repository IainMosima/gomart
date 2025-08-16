package notification

import (
	"context"

	"github.com/IainMosima/gomart/configs"
	"github.com/IainMosima/gomart/domains/order/schema"
)

type EmailNotificationServiceImpl struct {
	config *configs.Config
}

func NewEmailNotificationServiceImpl(config *configs.Config) *EmailNotificationServiceImpl {
	return &EmailNotificationServiceImpl{
		config: config,
	}
}

func (e EmailNotificationServiceImpl) SendOrderNotificationEmail(ctx context.Context, order *schema.OrderResponse) error {
	//TODO implement me
	panic("implement me")
}
