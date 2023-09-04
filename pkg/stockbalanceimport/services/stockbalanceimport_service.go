package services

import (
	"fmt"
	"math"
	"smlcloudplatform/pkg/stockbalanceimport/models"
	"smlcloudplatform/pkg/stockbalanceimport/repositories"
	stockbalance_models "smlcloudplatform/pkg/transaction/stockbalance/models"
	"smlcloudplatform/pkg/transaction/stockbalance/services"
	"time"
)

type IStockBalanceImportService interface {
	CreateTask(shopID string, req models.StockBalanceImportTaskRequest) (models.StockBalanceImportTask, error)
	GetTaskPart(shopID string, partID string) (models.StockBalanceImportPartCache, error)
	GetTaskMeta(shopID string, taskID string) (models.StockBalanceImportMeta, error)
	SaveTaskPart(shopID string, partID string, details []stockbalance_models.StockBalanceDetail) error
	SaveTaskComplete(shopID string, authUsername string, taskID string) (models.StockBalanceImportMeta, error)
}

type StockBalanceImportService struct {
	chunkSize           int
	sizeID              int
	cacheExpire         time.Duration
	cacheRepo           repositories.IStockBalanceImportCacheRepository
	stockBalanceService services.IStockBalanceHttpService
	GenerateID          func(int) string
}

func NewStockBalanceImportService(
	cacheRepo repositories.IStockBalanceImportCacheRepository,
	stockBalanceService services.IStockBalanceHttpService,
	GenerateID func(int) string,
) *StockBalanceImportService {
	return &StockBalanceImportService{
		chunkSize:           500,
		sizeID:              12,
		cacheExpire:         time.Minute * 60,
		cacheRepo:           cacheRepo,
		stockBalanceService: stockBalanceService,
		GenerateID:          GenerateID,
	}
}

func (svc *StockBalanceImportService) CreateTask(shopID string, req models.StockBalanceImportTaskRequest) (models.StockBalanceImportTask, error) {

	result := models.StockBalanceImportTask{}

	taskID := svc.GenerateID(svc.sizeID)

	result.TaskID = taskID
	result.TotalItem = req.TotalItem
	result.Header = req.Header
	result.Parts = []models.StockBalanceImportPart{}

	totalPart := math.Ceil(float64(req.TotalItem) / float64(svc.chunkSize))
	result.ChunkSize = svc.chunkSize

	for i := 0; i < int(totalPart); i++ {
		partNumber := i + 1
		partID := fmt.Sprintf("%s-%d", taskID, partNumber)
		result.Parts = append(result.Parts, models.StockBalanceImportPart{
			PartID:     partID,
			PartNumber: partNumber,
		})
	}

	svc.createTaskInMemory(shopID, result)

	return result, nil
}

func (svc *StockBalanceImportService) createTaskInMemory(shopID string, task models.StockBalanceImportTask) {

	metaCache := models.StockBalanceImportMeta{
		TaskID:    task.TaskID,
		TotalItem: task.TotalItem,
		Status:    models.TaskStatusPending,
		Header:    task.Header,
	}

	for _, part := range task.Parts {
		partCache := models.StockBalanceImportPartCache{}

		partCache.TaskID = task.TaskID
		partCache.PartID = part.PartID
		partCache.PartNumber = part.PartNumber
		partCache.Status = 0
		partCache.Detail = []stockbalance_models.StockBalanceDetail{}

		svc.cacheRepo.CreatePart(shopID, part.PartID, partCache, svc.cacheExpire)

		metaCache.Parts = append(metaCache.Parts, models.StockBalanceImportPartMeta{
			PartID:     part.PartID,
			PartNumber: part.PartNumber,
			Status:     0,
		})
	}

	svc.cacheRepo.CreateMeta(shopID, task.TaskID, metaCache, svc.cacheExpire)
}

func (svc *StockBalanceImportService) GetTaskPart(shopID string, partID string) (models.StockBalanceImportPartCache, error) {

	result, err := svc.cacheRepo.GetPart(shopID, partID)

	if err != nil {
		return models.StockBalanceImportPartCache{}, err
	}

	if result.PartID == "" {
		result.PartID = partID
		result.Status = models.PartStatusNotFound
	}

	return result, nil
}

func (svc *StockBalanceImportService) GetTaskMeta(shopID string, TaskID string) (models.StockBalanceImportMeta, error) {

	result := models.StockBalanceImportMeta{
		TaskID: TaskID,
		Parts:  []models.StockBalanceImportPartMeta{},
	}
	taskMeta, err := svc.cacheRepo.GetMeta(shopID, TaskID)

	if err != nil {
		return models.StockBalanceImportMeta{}, err
	}

	if taskMeta.TaskID == "" {
		result.Status = models.TaskStatusNotFound
		return result, nil
	}

	if taskMeta.Status == models.TaskStatusSaveSucceded {
		return taskMeta, nil
	}

	for i, part := range taskMeta.Parts {
		partCache, err := svc.cacheRepo.GetPart(shopID, part.PartID)

		if err != nil {
			taskMeta.Parts[i].Status = models.PartStatusError
			continue
		}
		taskMeta.Parts[i].Status = partCache.Status
	}

	result = taskMeta
	result.Status = svc.taskStatus(taskMeta.Parts)

	svc.cacheRepo.UpdateMeta(shopID, TaskID, result)

	return result, nil
}

func (svc *StockBalanceImportService) SaveTaskPart(shopID string, taskID string, details []stockbalance_models.StockBalanceDetail) error {

	doc, err := svc.cacheRepo.GetPart(shopID, taskID)

	if err != nil {
		return err
	}

	if doc.PartID == "" {
		return fmt.Errorf("part not found")
	}

	if doc.Status == models.PartStatusDone {
		return fmt.Errorf("part is done")
	}

	doc.Detail = details
	doc.Status = models.PartStatusDone
	err = svc.cacheRepo.UpdatePart(shopID, taskID, doc)

	if err != nil {
		doc.Status = models.PartStatusError
		doc.Detail = []stockbalance_models.StockBalanceDetail{}
		svc.cacheRepo.UpdatePart(shopID, taskID, doc)
		return err
	}

	return nil
}

func (svc *StockBalanceImportService) SaveTaskComplete(shopID string, authUsername string, taskID string) (models.StockBalanceImportMeta, error) {

	result := models.StockBalanceImportMeta{
		TaskID: taskID,
		Parts:  []models.StockBalanceImportPartMeta{},
	}
	meta, err := svc.GetTaskMeta(shopID, taskID)

	if err != nil {
		result.Status = models.TaskStatusError
		return result, err
	}

	if meta.Status == models.TaskStatusSaveSucceded {
		return meta, nil
	}

	if meta.TaskID == "" {
		result.Status = models.TaskStatusNotFound
		return result, nil
	}

	result = meta
	result.Status = svc.taskStatus(meta.Parts)

	if result.Status != models.TaskStatusDone {
		return result, fmt.Errorf("task is not done")
	}

	tempDetails := []stockbalance_models.StockBalanceDetail{}

	for _, part := range meta.Parts {
		partCache, err := svc.cacheRepo.GetPart(shopID, part.PartID)

		if err != nil {
			result.Status = models.TaskStatusError
		}

		tempDetails = append(tempDetails, partCache.Detail...)
	}

	tempTransaction := stockbalance_models.StockBalance{}
	tempTransaction.StockBalanceHeader = meta.Header
	tempTransaction.Details = &tempDetails

	svc.stockBalanceService.CreateStockBalance(shopID, authUsername, tempTransaction)

	result.Status = models.TaskStatusSaveSucceded
	svc.cacheRepo.UpdateMeta(shopID, taskID, result)

	return result, nil
}

func (svc *StockBalanceImportService) taskStatus(taskParts []models.StockBalanceImportPartMeta) models.TaskStatus {
	var taskStatus models.TaskStatus = models.TaskStatusPending

	for _, part := range taskParts {
		if part.Status == models.PartStatusError {
			taskStatus = models.TaskStatusError
			break
		}

		if part.Status != models.PartStatusDone {
			taskStatus = models.TaskStatusProcessing
			break
		}

		taskStatus = models.TaskStatusDone
	}

	return taskStatus
}
