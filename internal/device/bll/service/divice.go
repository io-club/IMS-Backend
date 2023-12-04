package service

import (
	"context"
	"ims-server/internal/device/bll/pack"
	"ims-server/internal/device/dal/model"

	"ims-server/internal/device/param"

	egoerror "ims-server/pkg/error"
	"ims-server/internal/device/dal/repo"
	"ims-server/pkg/util"
)

type sensorService struct {
}



func (s *sensorService) CreateMqttDate(ctx context.Context, req *param.MqDateRequest) (*param.CreateMqDateRes, error) {
	sensorDate := &model.Sensor{
		Type:       req.Type,
		SensorDate: req.SensorDate,
		TerminalID: req.TerminalID,
	}
	err := repo.NewSensorRepo().Create(ctx, sensorDate)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	res := pack.ToMqDateResponse(sensorDate)
	return &param.CreateMqDateRes{
		MqDateResponse: res,
	}, nil
}

func (s *sensorService) GetMqttDateByID(ctx context.Context, req *param.GetMqDateByID) (*param.GetMqDateByIDRes, error) {
	id := req.ID
	sensorDate, err := repo.NewSensorRepo().Get(ctx, id)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}

	resp := pack.ToMqDateResponse(sensorDate)
	return &param.GetMqDateByIDRes{
		MqDateResponse: resp,
	}, nil
}

func (s *sensorService) MGetSensorByIDs(ctx context.Context, req *param.MGetMqDateByIDsReq) (*param.MGetMqDateByIDsRes, error) {
	res, err := repo.NewSensorRepo().MGet(ctx, req.IDs)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}

	resp := []param.GetMqDateByIDRes{}
	for _, sensorDate := range res {
		info := pack.ToMqDateResponse(&sensorDate)
		resp = append(resp, param.GetMqDateByIDRes{
			MqDateResponse: info,
		})
	}

	return &param.MGetMqDateByIDsRes{
		List: resp,
	}, nil
}

func (s *sensorService) UpdateMqDateByID(ctx context.Context, req *param.UpdateMqDateByIDReq) (*param.UpdateMqDateByIDRes, error) {
	_, err1 := repo.NewSensorRepo().Get(ctx, req.ID)
	if err1 != nil {
		return nil, egoerror.ErrNotFound
	}

	mqDateMap := util.RequestToSnakeMapWithIgnoreZeroValueAndIDKey(req)

	update, err2 := repo.NewSensorRepo().Update(ctx, req.ID, mqDateMap)
	if err2 != nil {
		return nil, egoerror.ErrInvalidParam
	}

	resp := pack.ToMqDateResponse(update)
	return &param.UpdateMqDateByIDRes{
		MqDateResponse: resp,
	}, nil
}

func (s *sensorService) DeleteMqDaterByID(ctx context.Context, req *param.DeleteMqDateByIDReq) (*param.DeleteMqDateByIDRes, error) {
	_, err := repo.NewSensorRepo().Get(ctx, req.ID)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	err = repo.NewSensorRepo().Delete(ctx, req.ID)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	return &param.DeleteMqDateByIDRes{}, nil
}
