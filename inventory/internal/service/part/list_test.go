package part

import (
	"github.com/brianvoe/gofakeit/v6"

	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
	"github.com/Artyom099/factory/inventory/internal/service/model"
)

func (s *ServiceSuite) TestListSuccess() {
	var (
		uuid1        = gofakeit.UUID()
		uuid2        = gofakeit.UUID()
		name1        = gofakeit.Name()
		name2        = gofakeit.Name()
		description1 = gofakeit.BeerAlcohol()
		description2 = gofakeit.BeerAlcohol()
		price1       = float64(gofakeit.Number(100, 1000))
		price2       = float64(gofakeit.Number(100, 1000))

		serviceRequestDto = model.ModelPartFilter{
			Uuids:                 []string{uuid1, uuid2},
			Names:                 []string{},
			Categories:            []model.Category{},
			ManufacturerCountries: []string{},
			Tags:                  []string{},
		}

		repoRequestDto = repoModel.RepoPartFilter{
			Uuids:                 []string{uuid1, uuid2},
			Names:                 []string{},
			Categories:            []repoModel.Category{},
			ManufacturerCountries: []string{},
			Tags:                  []string{},
		}

		repoResponseDto = []repoModel.RepoPart{
			{Uuid: uuid1, Name: name1, Description: description1, Price: price1},
			{Uuid: uuid2, Name: name2, Description: description2, Price: price2},
		}
	)

	s.partRepository.On("List", s.ctx, repoRequestDto).Return(repoResponseDto, nil)

	res, err := s.service.List(s.ctx, serviceRequestDto)
	s.Require().NoError(err)
	s.Require().Equal(len(res), 2)
}

func (s *ServiceSuite) TestListRepoError() {
	var (
		repoErr = gofakeit.Error()
		uuid1   = gofakeit.UUID()
		uuid2   = gofakeit.UUID()

		serviceRequestDto = model.ModelPartFilter{
			Uuids:                 []string{uuid1, uuid2},
			Names:                 []string{},
			Categories:            []model.Category{},
			ManufacturerCountries: []string{},
			Tags:                  []string{},
		}

		repoRequestDto = repoModel.RepoPartFilter{
			Uuids:                 []string{uuid1, uuid2},
			Names:                 []string{},
			Categories:            []repoModel.Category{},
			ManufacturerCountries: []string{},
			Tags:                  []string{},
		}
	)

	s.partRepository.On("List", s.ctx, repoRequestDto).Return([]repoModel.RepoPart{}, repoErr)

	_, err := s.service.List(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
