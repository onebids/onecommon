
namespace go errno

struct ErrNo {
  1: i64 ErrCode,
  2: string ErrMsg
}

struct Response {
  1: i64 Code,
  2: string Message,
  3: binary Data
}

const ErrNo Success = {"ErrCode": 0, "ErrMsg": "success"}
const ErrNo NoRoute = {"ErrCode": 1, "ErrMsg": "no route"}
const ErrNo NoMethod = {"ErrCode": 2, "ErrMsg": "no method"}
const ErrNo BadRequest = {"ErrCode": 3, "ErrMsg": "bad request"}
const ErrNo ParamsErr = {"ErrCode": 4, "ErrMsg": "params error"}
const ErrNo AuthorizeFail = {"ErrCode": 5, "ErrMsg": "authorize failed"}
const ErrNo TooManyRequest = {"ErrCode": 6, "ErrMsg": "too many requests"}
const ErrNo ServiceErr = {"ErrCode": 7, "ErrMsg": "service error"}
const ErrNo RPCUserSrvErr = {"ErrCode": 40100, "ErrMsg": "rpc user service error"}
const ErrNo UserSrvErr = {"ErrCode": 30100, "ErrMsg": "user service error"}
const ErrNo RPCUserAdminSrvErr = {"ErrCode": 40200, "ErrMsg": "rpc system service error"}
const ErrNo UserAdminSrvErr = {"ErrCode": 30200, "ErrMsg": "system service error"}
const ErrNo RPCProductSrvErr = {"ErrCode": 40300, "ErrMsg": "rpc car service error"}
const ErrNo ProductSrvErr = {"ErrCode": 30300, "ErrMsg": "car service error"}
const ErrNo RPCOrderSrvErr = {"ErrCode": 40400, "ErrMsg": "rpc profile service error"}
const ErrNo OrderSrvErr = {"ErrCode": 30400, "ErrMsg": "profile service error"}
const ErrNo RPCActivitySrvErr = {"ErrCode": 40500, "ErrMsg": "rpc trip service error"}
const ErrNo ActivitySrvErr = {"ErrCode": 30500, "ErrMsg": "trip service error"}
const ErrNo RPCCartSrvErr = {"ErrCode": 40600, "ErrMsg": "rpc trip service error"}
const ErrNo CartSrvErr = {"ErrCode": 30600, "ErrMsg": "trip service error"}
const ErrNo RecordNotFound = {"ErrCode": 18, "ErrMsg": "record not found"}
const ErrNo RecordAlreadyExist = {"ErrCode": 19, "ErrMsg": "record already exist"}
const ErrNo DirtyData = {"ErrCode": 20, "ErrMsg": "dirty data"}