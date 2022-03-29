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
	CreateWithGuid(shopId string, username string, guid string, doc models.Member) (string, error)
	CreateMember(shopId string, username string, doc models.Member) (string, error)
	UpdateMember(guid string, shopId string, username string, doc models.Member) error
	DeleteMember(guid string, shopId string, username string) error
	InfoMember(guid string, shopId string) (models.MemberInfo, error)
	SearchMember(shopId string, q string, page int, limit int) ([]models.MemberInfo, paginate.PaginationData, error)
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

// Find guid in postgresql index
func (svc MemberService) IsExistsGuid(shopID string, guidFixed string) (bool, error) {

	count, err := svc.memberPgRepo.Count(shopID, guidFixed)
	if err != nil {
		return false, err
	}

	if count == 0 {
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

func (svc MemberService) CreateWithGuid(shopId string, username string, guid string, doc models.Member) (string, error) {
	dataDoc := models.MemberDoc{}

	newGuid := guid
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

	err := svc.memberRepo.Delete(guid, shopId, username)
	if err != nil {
		return err
	}

	return nil
}

func (svc MemberService) InfoMember(guid string, shopId string) (models.MemberInfo, error) {
	doc, err := svc.memberRepo.FindByGuid(guid, shopId)

	if err != nil {
		return models.MemberInfo{}, err
	}

	return doc.MemberInfo, nil
}

func (svc MemberService) SearchMember(shopId string, q string, page int, limit int) ([]models.MemberInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.memberRepo.FindPage(shopId, q, page, limit)

	if err != nil {
		return docList, pagination, err
	}

	return docList, pagination, nil
}
