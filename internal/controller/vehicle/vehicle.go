package vehicle

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/thegoodparticle/vehicle-data-layer/internal/model"
	"github.com/thegoodparticle/vehicle-data-layer/internal/store"
)

type Controller struct {
	repository store.Interface
}

type Interface interface {
	ListOne(ID string) (entity model.VehicleViolations, err error)
	ListAll() (entities []model.VehicleViolations, err error)
	Create(entity *model.VehicleViolations) (string, error)
	Update(ID string, entity *model.VehicleViolations) error
	Remove(ID string) error
}

func NewController(repository store.Interface) Interface {
	return &Controller{repository: repository}
}

func (c *Controller) ListOne(id string) (entity model.VehicleViolations, err error) {
	entity.VehicleRegID = id
	response, err := c.repository.FindOne(entity.GetFilterId(), entity.TableName())
	if err != nil {
		return entity, err
	}
	return model.ParseDynamoAtributeToStruct(response.Item)
}

func (c *Controller) ListAll() (entities []model.VehicleViolations, err error) {
	entities = []model.VehicleViolations{}
	var entity model.VehicleViolations

	filter := expression.Name("owner_name").NotEqual(expression.Value(""))
	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return entities, err
	}

	response, err := c.repository.FindAll(condition, entity.TableName())
	if err != nil {
		return entities, err
	}

	if response != nil {
		for _, value := range response.Items {
			entity, err := model.ParseDynamoAtributeToStruct(value)
			if err != nil {
				return entities, err
			}
			entities = append(entities, entity)
		}
	}

	return entities, nil
}

func (c *Controller) Create(entity *model.VehicleViolations) (string, error) {
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()
	_, err := c.repository.CreateOrUpdate(entity.GetMap(), entity.TableName())
	return entity.VehicleRegID, err
}

func (c *Controller) Update(id string, entity *model.VehicleViolations) error {
	found, err := c.ListOne(id)
	if err != nil {
		return err
	}

	found.VehicleRegID = id

	if entity.OwnerName != "" {
		found.OwnerName = entity.OwnerName
	}

	found.TotalFine = entity.TotalFine

	if entity.Violations != nil {
		found.Violations = entity.Violations
	}

	found.UpdatedAt = time.Now()
	_, err = c.repository.CreateOrUpdate(found.GetMap(), entity.TableName())
	return err
}

func (c *Controller) Remove(id string) error {
	entity, err := c.ListOne(id)
	if err != nil {
		return err
	}
	_, err = c.repository.Delete(entity.GetFilterId(), entity.TableName())
	return err
}
