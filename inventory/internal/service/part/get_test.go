package part

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"

	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
	"github.com/Artyom099/factory/inventory/internal/service/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		partUUID      = gofakeit.UUID()
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
		stringValue   = gofakeit.Word()
		metadata      = map[string]*repoModel.Value{
			"key1": {StringValue: &stringValue},
		}

		serviceRequestDto = model.PartGetServiceRequest{
			Uuid: partUUID,
		}

		repoRequestDto = repoModel.PartGetRepoRequest{
			Uuid: partUUID,
		}

		repoResponseDto = repoModel.PartGetRepoResponse{
			Part: repoModel.Part{
				Uuid:          partUUID,
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
				Tags:      tags,
				Metadata:  metadata,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		}
	)

	s.partRepository.On("Get", s.ctx, repoRequestDto).Return(repoResponseDto, nil)

	res, err := s.service.Get(s.ctx, serviceRequestDto)
	s.Require().NoError(err)
	s.Require().Equal(partUUID, res.Part.Uuid)
	s.Require().Equal(name, res.Part.Name)
	s.Require().Equal(description, res.Part.Description)
	s.Require().Equal(price, res.Part.Price)
	s.Require().Equal(stockQuantity, res.Part.StockQuantity)
	s.Require().Equal(category, int(res.Part.Category))
	s.Require().Equal(width, res.Part.Dimensions.Width)
	s.Require().Equal(height, res.Part.Dimensions.Height)
	s.Require().Equal(length, res.Part.Dimensions.Length)
	s.Require().Equal(weight, res.Part.Dimensions.Weight)
	s.Require().Equal(manufName, res.Part.Manufacturer.Name)
	s.Require().Equal(manufCountry, res.Part.Manufacturer.Country)
	s.Require().Equal(tags, res.Part.Tags)
}
