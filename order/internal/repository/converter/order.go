package converter

import (
	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	servModel "github.com/Artyom099/factory/order/internal/service/model"
)

func ToRepoOrder(dto servModel.Order) repoModel.RepoOrder {
	return repoModel.RepoOrder{
		OrderUUID:       dto.OrderUUID,
		UserUUID:        dto.UserUUID,
		PartUuids:       dto.PartUuids,
		TotalPrice:      dto.TotalPrice,
		TransactionUUID: &dto.TransactionUUID,
		PaymentMethod:   repoModel.OrderPaymentMethod(dto.PaymentMethod),
		Status:          repoModel.OrderStatus(dto.Status),
	}
}

func ToModelOrder(dto repoModel.RepoOrder) servModel.Order {
	var transactionUUID string
	if dto.TransactionUUID != nil {
		transactionUUID = *dto.TransactionUUID
	}

	return servModel.Order{
		OrderUUID:       dto.OrderUUID,
		UserUUID:        dto.UserUUID,
		PartUuids:       dto.PartUuids,
		TotalPrice:      dto.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   servModel.OrderPaymentMethod(dto.PaymentMethod),
		Status:          servModel.OrderStatus(dto.Status),
	}
}
