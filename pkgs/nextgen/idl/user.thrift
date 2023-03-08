namespace go user

enum ErrCode {
    SuccessCode                = 0,
    ServiceErrCode             = 10001,
    ParamErrCode               = 10002,
    UserAlreadyExistErrCode    = 10003,
    AuthorizationFailedErrCode = 10004,
}

struct BaseResp {
    1: required i64 status_code,
    2: required string status_message,
    3: required i64 service_time,
}

struct User {
    1: required i64 user_id,
    2: required string username,
    3: required string nickname,
    4: required string avater,
    5: required i64 create_time,
    6: required i64 update_time,
}

struct CreateUserRequest {
    1: required string username,
}

struct CreateUserResponse {
    1: required BaseResp base_resp,
    2: required i64 user_id,
}

struct DeleteUserRequest {
    1: required i64 user_id,
}

struct DeleteUserResponse {
    1: required BaseResp base_resp,
}

struct UpdateUserRequest {
    1: required i64 user_id,
    2: optional string nickname,
    3: optional string avater,
}

struct UpdateUserResponse {
    1: required BaseResp base_resp,
}

struct QueryUserRequest {
    1: optional string username,
    2: optional string nickname,
    3: required i32 limit,
    4: required i32 offset,
}

struct QueryUserResponse {
    1: required BaseResp base_resp,
    2: required list<User> users,
    3: required i32 total,
}

struct MGetUserRequest {
    1: required list<i64> user_ids,
}

struct MGetUserResponse {
    1: required BaseResp base_resp,
    2: required list<User> users,
    3: required i32 total,
}

service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req),
    DeleteUserResponse DeleteUser(1: DeleteUserRequest req),
    UpdateUserResponse UpdateUser(1: UpdateUserRequest req),
    QueryUserResponse QueryUser(1: QueryUserRequest req),
    MGetUserResponse MGetUser(1: MGetUserRequest req),
}