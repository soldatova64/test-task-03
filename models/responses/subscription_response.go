package responses

import (
	"test-task-03/entity"
	"test-task-03/models"
)

type SubscriptionResponse struct {
	Meta models.Meta           `json:"meta"`
	Data []entity.Subscription `json:"data"`
}
