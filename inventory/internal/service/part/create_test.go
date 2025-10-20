package part

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"

	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
	"github.com/Artyom099/factory/inventory/internal/service/model"
)

func (s *ServiceSuite) TestCreateSuccess() {
	var (
		name          = gofakeit.Name()
		description   = gofakeit.BeerAlcohol()
		expectedUUID  = gofakeit.UUID()
		price         = float64(gofakeit.Number(100, 1000))
		stockQuantity = int64(gofakeit.Number(1, 100))
		category      = gofakeit.Number(1, 5)
		width         = gofakeit.Float64()
		height        = gofakeit.Float64()
		length        = gofakeit.Float64()
		weight        = gofakeit.Float64()
		manufName     = gofakeit.Company()
		manufCountry  = gofakeit.Country()
		tags          = []string{gofakeit.Word(), gofakeit.Word()}
		createdAt     = time.Now()
		updatedAt     = time.Now()

		serviceRequestDto = model.PartCreateServiceRequest{
			Name:          name,
			Description:   description,
			Price:         price,
			StockQuantity: stockQuantity,
			Category:      model.Category(category),
			Dimensions: &model.Dimensions{
				Width:  width,
				Height: height,
				Length: length,
				Weight: weight,
			},
			Manufacturer: &model.Manufacturer{
				Name:    manufName,
				Country: manufCountry,
			},
			Tags: tags,
			// Metadata      map[string]*Value
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		repoRequestDto = repoModel.PartCreateRepoRequest{
			Part: repoModel.Part{
				Name:          name,
				Description:   description,
				Price:         price,
				StockQuantity: stockQuantity,
				Category:      repoModel.Category(category),
				Dimensions: &repoModel.Dimensions{
					Width:  width,
					Height: height,
					Length: length,
					Weight: weight,
				},
				Manufacturer: &repoModel.Manufacturer{
					Name:    manufName,
					Country: manufCountry,
				},
				Tags: tags,
				// Metadata      map[string]*Value
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		}
	)

	s.partRepository.On("Create", s.ctx, repoRequestDto).Return(expectedUUID, nil)

	uuid, err := s.service.Create(s.ctx, serviceRequestDto)
	s.Require().NoError(err)
	s.Require().Equal(expectedUUID, uuid)
}

func (s *ServiceSuite) TestCreateRepoError() {
	var (
		repoErr       = gofakeit.Error()
		name          = gofakeit.Name()
		description   = gofakeit.BeerAlcohol()
		price         = float64(gofakeit.Number(100, 1000))
		stockQuantity = int64(gofakeit.Number(1, 100))
		category      = gofakeit.Number(1, 5)
		width         = gofakeit.Float64()
		height        = gofakeit.Float64()
		length        = gofakeit.Float64()
		weight        = gofakeit.Float64()
		manufName     = gofakeit.Company()
		manufCountry  = gofakeit.Country()
		tags          = []string{gofakeit.Word(), gofakeit.Word()}
		createdAt     = time.Now()
		updatedAt     = time.Now()

		serviceRequestDto = model.PartCreateServiceRequest{
			Name:          name,
			Description:   description,
			Price:         price,
			StockQuantity: stockQuantity,
			Category:      model.Category(category),
			Dimensions: &model.Dimensions{
				Width:  width,
				Height: height,
				Length: length,
				Weight: weight,
			},
			Manufacturer: &model.Manufacturer{
				Name:    manufName,
				Country: manufCountry,
			},
			Tags: tags,
			// Metadata      map[string]*Value
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		repoRequestDto = repoModel.PartCreateRepoRequest{
			Part: repoModel.Part{
				Name:          name,
				Description:   description,
				Price:         price,
				StockQuantity: stockQuantity,
				Category:      repoModel.Category(category),
				Dimensions: &repoModel.Dimensions{
					Width:  width,
					Height: height,
					Length: length,
					Weight: weight,
				},
				Manufacturer: &repoModel.Manufacturer{
					Name:    manufName,
					Country: manufCountry,
				},
				Tags: tags,
				// Metadata      map[string]*Value
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		}
	)

	s.partRepository.On("Create", s.ctx, repoRequestDto).Return("", repoErr)

	uuid, err := s.service.Create(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(uuid)
}
