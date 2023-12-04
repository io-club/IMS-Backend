package service

import (
	"context"
	"ims-server/internal/user/bll/pack"
	"ims-server/internal/user/dal/model"
	"ims-server/internal/user/param"
	egoerror "ims-server/pkg/error"
	ioginx "ims-server/pkg/ginx"
	"ims-server/pkg/util"
)

type sensorService struct {
}

type sensorRepo struct {
	ioginx.IRepo[model.Sensor]
}

func NewSensorRepo() *sensorRepo {
	return &sensorRepo{}
}

func (s *sensorService) CreateMqttDate(ctx context.Context, msg string, req *param.MqDateRequest) (*param.CreateMqDateRes, error) {
	sensorDate := &model.Sensor{
		Type:       req.Type,
		SensorDate: req.SensorDate,
		TerminalID: req.TerminalID,
	}
	err := NewSensorRepo().Create(ctx, sensorDate)
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
	sensorDate, err := NewSensorRepo().Get(ctx, id)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}

	resp := pack.ToMqDateResponse(sensorDate)
	return &param.GetMqDateByIDRes{
		MqDateResponse: resp,
	}, nil
}

func (s *sensorService) MGetSensorByIDs(ctx context.Context, req *param.MGetMqDateByIDsReq) (*param.MGetMqDateByIDsRes, error) {
	res, err := NewSensorRepo().MGet(ctx, req.IDs)
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
	_, err1 := NewSensorRepo().Get(ctx, req.ID)
	if err1 != nil {
		return nil, egoerror.ErrNotFound
	}

	mqDateMap := util.RequestToSnakeMapWithIgnoreZeroValueAndIDKey(req)

	update, err2 := NewSensorRepo().Update(ctx, req.ID, mqDateMap)
	if err2 != nil {
		return nil, egoerror.ErrInvalidParam
	}

	resp := pack.ToMqDateResponse(update)
	return &param.UpdateMqDateByIDRes{
		MqDateResponse: resp,
	}, nil
}

func (s *sensorService) DeleteMqDaterByID(ctx context.Context, req *param.DeleteMqDateByIDReq) (*param.DeleteMqDateByIDRes, error) {
	_, err := NewSensorRepo().Get(ctx, req.ID)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	err = NewSensorRepo().Delete(ctx, req.ID)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	return &param.DeleteMqDateByIDRes{}, nil
}
