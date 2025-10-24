package converter

import (
	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	servModel "github.com/Artyom099/factory/order/internal/service/model"
)

func ToRepoOrder(dto servModel.Order) repoModel.RepoOrder {
	return repoModel.RepoOrder{
		UserUUID:   dto.UserUUID,
		PartUuids:  dto.PartUuids,
		TotalPrice: dto.TotalPrice,
	}
}

func ToModelOrder(dto repoModel.RepoOrder) servModel.Order {
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
