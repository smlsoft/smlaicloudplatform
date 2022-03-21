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
	CreateMember(shopId string, username string, doc models.Member) (string, error)
	UpdateMember(guid string, shopId string, username string, doc models.Member) error
	DeleteMember(guid string, shopId string, username string) error
	InfoMember(guid string, shopId string) (models.MemberDoc, error)
	SearchMember(shopId string, q string, page int, limit int) ([]models.MemberDoc, paginate.PaginationData, error)
}

type MemberService struct {
	memberRepo IMemberRepository
}

func NewMemberService(memberRepo IMemberRepository) *MemberService {
	return &MemberService{
		memberRepo: memberRepo,
	}
}

func (svc MemberService) CreateMember(shopId string, username string, doc models.Member) (string, error) {
	dataDoc := models.MemberDoc{}

	newGuid := utils.NewGUID()
	dataDoc.GuidFixed = newGuid
	dataDoc.ShopID = shopId
	dataDoc.CreatedBy = username
	dataDoc.CreatedAt = time.Now()

	dataDoc.Member = doc

	_, err := svc.memberRepo.Create(dataDoc)

	if err != nil {
		return "", err
	}

	return newGuid, nil
}

func (svc MemberService) UpdateMember(guid string, shopId string, username string, doc models.Member) error {

	findDoc, err := svc.memberRepo.FindByGuid(guid, shopId)

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

func (svc MemberService) DeleteMember(guid string, shopId string, username string) error {

	err := svc.memberRepo.Delete(guid, shopId)
	if err != nil {
		return err
	}

	return nil
}

func (svc MemberService) InfoMember(guid string, shopId string) (models.MemberDoc, error) {
	doc, err := svc.memberRepo.FindByGuid(guid, shopId)

	if err != nil {
		return models.MemberDoc{}, err
	}

	return doc, nil
}

func (svc MemberService) SearchMember(shopId string, q string, page int, limit int) ([]models.MemberDoc, paginate.PaginationData, error) {
	docList, pagination, err := svc.memberRepo.FindPage(shopId, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}
