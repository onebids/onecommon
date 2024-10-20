namespace go enums

enum Err {
    Success            = 0,
    NoRoute            = 1,
    NoMethod           = 2,
    BadRequest         = 10000,
    ParamsErr          = 10001,
    AuthorizeFail      = 10002,
    TooManyRequest     = 10003,
    ServiceErr         = 20000,

    RPCUserSrvErr      = 30001,
    RPCUserAdminSrvErr      = 30002,
    RPCOrderSrvErr       = 30003,
    RPCProductSrvErr   = 30004,
    RPCActivitySrvErr      = 30005,

    UserSrvErr         = 40001,
    UserAdminSrvErr      = 40002,
    OrderSrvErr          = 40003,
    ProductSrvErr      = 40004,
    ActivitySrvErr        = 40005,

    RecordNotFound     = 80000,
    RecordAlreadyExist = 80001,
    DirtyData          = 80003,
}