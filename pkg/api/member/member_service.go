package member

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMemberService interface {
	IsExistsGuid(shopID string, guidFixed string) (bool, error)
	CreateIndex(doc models.MemberIndex) error
	CreateWithGuid(shopID string, username string, guid string, doc models.Member) (string, error)
	Create(shopID string, username string, doc models.Member) (string, error)
	Update(shopID string, guid string, username string, doc models.Member) error
	Delete(shopID string, guid string, username string) error
	Info(shopID string, guid string) (models.MemberInfo, error)
	Search(shopID string, q string, page int, limit int) ([]models.MemberInfo, paginate.PaginationData, error)
}

type MemberService struct {
	memberRepo   IMemberRepository
	memberPgRepo IMemberPGRepository
}

func NewMemberService(memberRepo IMemberRepository, memberPgRepo IMemberPGRepository) MemberService {
	return MemberService{
		memberRepo:   memberRepo,
		memberPgRepo: memberPgRepo,
	}
}

func (svc MemberService) IsExistsGuid(shopID string, guidFixed string) (bool, error) {

	findDoc, err := svc.memberRepo.FindByGuid(shopID, guidFixed)

	if err != nil {
		return false, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return false, nil
	}

	return true, nil

}

func (svc MemberService) CreateIndex(doc models.MemberIndex) error {

	err := svc.memberPgRepo.Create(doc)
	if err != nil {
		return err
	}

	return nil

}

func (svc MemberService) CreateWithGuid(shopID string, username string, guid string, doc models.Member) (string, error) {
	dataDoc := models.MemberDoc{}

	newGuid := guid
	dataDoc.GuidFixed = newGuid
	dataDoc.ShopID = shopID
	dataDoc.CreatedBy = username
	dataDoc.CreatedAt = time.Now()

	dataDoc.Member = doc

	_, err := svc.memberRepo.Create(dataDoc)

	if err != nil {
		return "", err
	}

	return newGuid, nil
}

func (svc MemberService) Create(shopID string, username string, doc models.Member) (string, error) {
	dataDoc := models.MemberDoc{}

	newGuid := utils.NewGUID()
	dataDoc.GuidFixed = newGuid
	dataDoc.ShopID = shopID
	dataDoc.CreatedBy = username
	dataDoc.CreatedAt = time.Now()

	dataDoc.Member = doc

	_, err := svc.memberRepo.Create(dataDoc)

	if err != nil {
		return "", err
	}

	return newGuid, nil
}

func (svc MemberService) Update(shopID string, guid string, username string, doc models.Member) error {

	findDoc, err := svc.memberRepo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("guid invalid")
	}

	findDoc.UpdatedBy = username
	findDoc.UpdatedAt = time.Now()
	findDoc.Member = doc

	err = svc.memberRepo.Update(guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc MemberService) Delete(shopID string, guid string, username string) error {

	err := svc.memberRepo.Delete(shopID, guid, username)
	if err != nil {
		return err
	}

	return nil
}

func (svc MemberService) Info(shopID string, guid string) (models.MemberInfo, error) {
	doc, err := svc.memberRepo.FindByGuid(shopID, guid)

	if err != nil {
		return models.MemberInfo{}, err
	}

	return doc.MemberInfo, nil
}

func (svc MemberService) Search(shopID string, q string, page int, limit int) ([]models.MemberInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.memberRepo.FindPage(shopID, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}
