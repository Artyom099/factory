package order

import (
	"github.com/brianvoe/gofakeit/v6"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *ServiceSuite) TestCreateSuccess() {
	var (
		orderUuid = gofakeit.UUID()
		partUuid1 = gofakeit.UUID()
		partUuid2 = gofakeit.UUID()

		serviceRequestDto = model.OrderCreateServiceRequestDto{
			UserUUID:  gofakeit.UUID(),
			PartUuids: []string{partUuid1, partUuid2},
		}

		serviceResponseDto = model.OrderCreateServiceResponseDto{
			OrderUUID:  orderUuid,
			TotalPrice: 100.0,
		}

		repoRequestDto = repoModel.OrderCreateRepoRequestDto{
			UserUUID:   serviceRequestDto.UserUUID,
			PartUuids:  serviceRequestDto.PartUuids,
			TotalPrice: 100,
		}

		listPartsRequestDto = model.ListPartsFilter{
			Uuids: repoRequestDto.PartUuids,
		}

		listPartsResponseDto = model.ListPartsResponseDto{
			Parts: []*model.Part{
				{Uuid: partUuid1, Name: gofakeit.Name(), Price: 40.0},
				{Uuid: partUuid2, Name: gofakeit.Name(), Price: 60.0},
			},
		}
	)

	s.orderRepository.On("Create", s.ctx, repoRequestDto).Return(orderUuid, nil)
	s.inventoryClient.On("ListParts", s.ctx, listPartsRequestDto).Return(listPartsResponseDto, nil)

	res, err := s.service.Create(s.ctx, serviceRequestDto)
	s.Require().NoError(err)
	s.Require().Equal(res, serviceResponseDto)
}

func (s *ServiceSuite) TestCreateRepoError() {
	var (
		repoErr   = gofakeit.Error()
		partUuid1 = gofakeit.UUID()
		partUuid2 = gofakeit.UUID()

		serviceRequestDto = model.OrderCreateServiceRequestDto{
			UserUUID:  gofakeit.UUID(),
			PartUuids: []string{partUuid1, partUuid2},
		}

		repoRequestDto = repoModel.OrderCreateRepoRequestDto{
			UserUUID:   serviceRequestDto.UserUUID,
			PartUuids:  serviceRequestDto.PartUuids,
			TotalPrice: 100,
		}

		listPartsRequestDto = model.ListPartsFilter{
			Uuids: repoRequestDto.PartUuids,
		}

		listPartsResponseDto = model.ListPartsResponseDto{
			Parts: []*model.Part{
				{Uuid: partUuid1, Name: gofakeit.Name(), Price: 40.0},
				{Uuid: partUuid2, Name: gofakeit.Name(), Price: 60.0},
			},
		}
	)

	s.orderRepository.On("Create", s.ctx, repoRequestDto).Return("", repoErr)
	s.inventoryClient.On("ListParts", s.ctx, listPartsRequestDto).Return(listPartsResponseDto, nil)

	res, err := s.service.Create(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestCreateINotAllPartsExistInInventoryServiceError() {
	var (
		partUuid1 = gofakeit.UUID()
		partUuid2 = gofakeit.UUID()

		serviceRequestDto = model.OrderCreateServiceRequestDto{
			UserUUID:  gofakeit.UUID(),
			PartUuids: []string{partUuid1, partUuid2},
		}

		listPartsRequestDto = model.ListPartsFilter{
			Uuids: serviceRequestDto.PartUuids,
		}

		listPartsResponseDto = model.ListPartsResponseDto{
			Parts: []*model.Part{
				{Uuid: partUuid1, Name: gofakeit.Name(), Price: 40.0},
			},
		}
	)

	s.inventoryClient.On("ListParts", s.ctx, listPartsRequestDto).Return(listPartsResponseDto, nil)

	res, err := s.service.Create(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, model.ErrNotAllPartsExist)
}

func (s *ServiceSuite) TestCreateInventoryServiceInternalError() {
	var (
		inventoryServiceErr = gofakeit.Error()

		partUuid1 = gofakeit.UUID()
		partUuid2 = gofakeit.UUID()

		serviceRequestDto = model.OrderCreateServiceRequestDto{
			UserUUID:  gofakeit.UUID(),
			PartUuids: []string{partUuid1, partUuid2},
		}

		listPartsRequestDto = model.ListPartsFilter{
			Uuids: []string{partUuid1, partUuid2},
		}
	)

	s.inventoryClient.On("ListParts", s.ctx, listPartsRequestDto).Return(model.ListPartsResponseDto{}, inventoryServiceErr)

	res, err := s.service.Create(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, model.ErrListPartsError)
}
