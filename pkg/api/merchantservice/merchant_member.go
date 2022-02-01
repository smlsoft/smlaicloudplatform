package merchantservice

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (svc *MerchantService) SearchMember(ctx microservice.IServiceContext) error {

	authUsername := ctx.UserInfo().Username
	merchantID := ctx.Param("merchant_id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

	merchantList := []models.MemberInfo{}

	searchText := bson.M{"$regex": primitive.Regex{
		Pattern: ".*" + q + ".*",
		Options: "",
	}}

	pagination, err := pst.FindPage(&models.MemberInfo{}, limit, page, bson.M{"createdby": authUsername, "merchantId": merchantID, "$or": []interface{}{
		bson.M{"email": searchText},
		bson.M{"username": searchText},
		bson.M{"name": searchText},
	}}, &merchantList)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "pagination": pagination, "data": merchantList})
	return nil
}

func (svc *MerchantService) CreateMember(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantID := ctx.Param("merchant_id")

	input := ctx.ReadInput()

	memberReq := &models.MemberRequest{}
	err := json.Unmarshal([]byte(input), &memberReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(memberReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	findMember := &models.Member{}
	err = pst.FindOne(&models.Member{}, bson.M{"merchant_id": merchantID, "$or": []interface{}{
		bson.M{"email": memberReq.Email},
		bson.M{"username": memberReq.Username},
	}}, findMember)

	if err != nil && err.Error() != "mongo: no documents in result" {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if findMember.Email == memberReq.Email {
		ctx.ResponseError(400, "email is exists.")
		return err
	}

	if findMember.Username == memberReq.Username {
		ctx.ResponseError(400, "username is exists.")
		return err
	}

	hashPassword, err := utils.HashPassword(memberReq.Password)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	member := &models.Member{
		MerchantID: merchantID,
		Email:      memberReq.Email,
		Username:   memberReq.Username,
		Password:   hashPassword,
		Name:       memberReq.Name,
		CreatedBy:  authUsername,
		CreatedAt:  time.Now(),
	}

	idx, err := pst.Create(&models.Member{}, member)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "id": idx})
	return nil
}

func (svc *MerchantService) EditMember(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	input := ctx.ReadInput()

	memberReq := &models.MemberRequestEdit{}
	err := json.Unmarshal([]byte(input), &memberReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(memberReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	findIdx, _ := primitive.ObjectIDFromHex(id)
	findMember := &models.Member{}
	err = pst.FindByID(&models.Member{}, "guidFixed", findIdx, findMember)
	if err != nil && err.Error() != "mongo: no documents in result" {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if authUsername != findMember.CreatedBy {
		ctx.ResponseError(400, "merchant invalid")
		return err
	}

	if len(findMember.Email) < 1 || len(findMember.Username) < 1 {
		ctx.ResponseError(400, "member is not exists")
		return err
	}

	findMember.Email = memberReq.Email
	findMember.Username = memberReq.Username
	findMember.Name = memberReq.Name

	idx, _ := primitive.ObjectIDFromHex(id)
	err = pst.Update(&models.Member{}, findMember, "guidFixed", idx)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (svc *MerchantService) DeleteMember(ctx microservice.IServiceContext) error {
	username := ctx.UserInfo().Username
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	findMember := &models.Member{}

	idx, _ := primitive.ObjectIDFromHex(id)
	err := pst.FindOne(&models.Member{}, bson.M{"_id": idx, "createdby": username}, findMember)

	if err != nil && err.Error() != "mongo: no documents in result" {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if len(findMember.CreatedBy) < 1 {
		ctx.ResponseError(400, "id invalid")
		return err
	}

	err = pst.DeleteByID(&models.Member{}, id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}

func (svc *MerchantService) GetMemberInfo(ctx microservice.IServiceContext) error {
	username := ctx.UserInfo().Username
	merchantID := ctx.Param("merchant_id")
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	memberInfo := &models.MemberInfo{}

	err := pst.FindOne(&models.MemberInfo{}, bson.M{"guidFixed": id, "createdby": username, "merchant_id": merchantID}, memberInfo)

	if err != nil {
		ctx.ResponseError(400, "not found")
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "data": memberInfo})
	return nil
}

func (svc *MerchantService) ChangePasswordMember(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	id := ctx.Param("id")

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	if _, err := utils.HasPermissionMerchant(pst, ctx); err != nil {
		ctx.ResponseError(400, err.Error())
		return nil
	}

	input := ctx.ReadInput()

	memberReq := &models.MemberRequestPassword{}
	err := json.Unmarshal([]byte(input), &memberReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(memberReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	findMember := &models.Member{}
	err = pst.FindByID(&models.Member{}, "guidFixed", id, findMember)
	if err != nil && err.Error() != "mongo: no documents in result" {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if authUsername != findMember.CreatedBy {
		ctx.ResponseError(400, "merchant invalid")
		return err
	}

	if len(memberReq.Password) > 0 {
		hashPassword, err := utils.HashPassword(memberReq.Password)

		if err != nil {
			ctx.ResponseError(400, err.Error())
			return err
		}
		findMember.Password = hashPassword
	}

	err = pst.Update(&models.Member{}, findMember, "guidFixed", id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
	return nil
}
