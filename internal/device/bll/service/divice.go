package service

import (
	"context"
	"ims-server/internal/device/bll/pack"
	"ims-server/internal/device/dal/model"

	"ims-server/internal/device/param"

	"ims-server/internal/device/dal/repo"
	egoerror "ims-server/pkg/error"
	"ims-server/pkg/util"
)

type sensorService struct {
}

func (s *sensorService) CreateMqttDate(ctx context.Context, req *param.MqDateRequest) (*param.CreateMqttDateRes, error) {
	sensorDate := &model.SensorData{
		Type:       req.Type,
		SensorData: req.SensorDate,
		TerminalID: req.TerminalID,
	}
	err := repo.NewSensorRepo().Create(ctx, sensorDate)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	res := pack.ToMqttDateResponse(sensorDate)
	return &param.CreateMqttDateRes{
		MqttDateResponse: res,
	}, nil
}

func (s *sensorService) GetMqttDateByID(ctx context.Context, req *param.GetMqttDateByID) (*param.GetMqttDateByIDRes, error) {
	id := req.ID
	sensorDate, err := repo.NewSensorRepo().Get(ctx, id)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}

	resp := pack.ToMqttDateResponse(sensorDate)
	return &param.GetMqttDateByIDRes{
		MqttDateResponse: resp,
	}, nil
}

func (s *sensorService) MGetSensorByIDs(ctx context.Context, req *param.MGetMqttDateByIDsReq) (*param.MGetMqttDateByIDsRes, error) {
	res, err := repo.NewSensorRepo().MGet(ctx, req.IDs)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}

	resp := []param.GetMqttDateByIDRes{}
	for _, sensorDate := range res {
		info := pack.ToMqttDateResponse(&sensorDate)
		resp = append(resp, param.GetMqttDateByIDRes{
			MqttDateResponse: info,
		})
	}

	return &param.MGetMqttDateByIDsRes{
		List: resp,
	}, nil
}

func (s *sensorService) UpdateMqDateByID(ctx context.Context, req *param.UpdateMqttDateByIDReq) (*param.UpdateMqttDateByIDRes, error) {
	_, err1 := repo.NewSensorRepo().Get(ctx, req.ID)
	if err1 != nil {
		return nil, egoerror.ErrNotFound
	}

	mqDateMap := util.RequestToSnakeMapWithIgnoreZeroValueAndIDKey(req)

	update, err2 := repo.NewSensorRepo().Update(ctx, req.ID, mqDateMap)
	if err2 != nil {
		return nil, egoerror.ErrInvalidParam
	}

	resp := pack.ToMqttDateResponse(update)
	return &param.UpdateMqttDateByIDRes{
		MqttDateResponse: resp,
	}, nil
}

func (s *sensorService) DeleteMqDaterByID(ctx context.Context, req *param.DeleteMqttDateByIDReq) (*param.DeleteMqttDateByIDRes, error) {
	_, err := repo.NewSensorRepo().Get(ctx, req.ID)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	err = repo.NewSensorRepo().Delete(ctx, req.ID)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	return &param.DeleteMqttDateByIDRes{}, nil
}
