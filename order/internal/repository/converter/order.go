package converter

import (
	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	servModel "github.com/Artyom099/factory/order/internal/service/model"
)

func ModelToRepoOrder(
	dto servModel.Order,
	totalPrice float32,
) repoModel.RepoOrder {
	return repoModel.RepoOrder{
		UserUUID:   dto.UserUUID,
		PartUuids:  dto.PartUuids,
		TotalPrice: totalPrice,
	}
}

func RepoToModelOrder(dto repoModel.RepoOrder) servModel.Order {
	return servModel.Order{
		OrderUUID:       dto.OrderUUID,
		UserUUID:        dto.UserUUID,
		PartUuids:       dto.PartUuids,
		TotalPrice:      dto.TotalPrice,
		TransactionUUID: dto.TransactionUUID,
		PaymentMethod:   servModel.OrderPaymentMethod(dto.PaymentMethod),
		Status:          servModel.OrderStatus(dto.Status),
	}
}
