package web

import "driver-service/internal/config"

type NatNotificationer struct {
}

func NewNatNotificationer(cfg config.NatNotificationerConfig) *NatNotificationer {
	return &NatNotificationer{}
}
