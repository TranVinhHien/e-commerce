package services

import (
	"context"
	services_assets_sendMessage "new-project/services/assets/sendMessage"
)

const DISSCOUT = "discoumt_info"

func (s *service) NotiNewDiscount(ctx context.Context) error {

	listDiscounts, err := s.repository.GetDiscountForNoti(ctx)
	if err != nil {
		return nil
	}
	mes := services_assets_sendMessage.MaGiamGiaMoi(listDiscounts)
	s.firebase.SendToTopic(ctx, DISSCOUT, mes)
	return nil
}
