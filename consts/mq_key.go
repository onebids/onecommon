package consts

const (
	MQ_BOT_NOTIFY_KEY    = "BOT_NOTIFY"
	MQ_BOT_NOTIFY_ROUNTE = "BOT_NOTIFY_ROUTE"
	MQ_BOT_NOTIFY_QUEUE  = "BOT_NOTIFY_QUEUE"

	MqPopQueue       = "LuckeyIdCheckQueue"
	MqPopDelayPopKey = "DelayPopKeyV2"
	MqPopRoutingKey  = "MqPopRoutingKey"

	// 订单过期队列
	MqOrderExpiredQueue      = "OrderExpiredQueue"
	MqOrderExpiredKey        = "OrderExpiredKey"
	MqOrderExpiredRoutingKey = "OrderExpiredRoutingKey"
	// 开奖活动队列
	MqActivityQueue      = "OpenActivityQueue"
	MqActivityKey        = "OpenActivityKey"
	MqActivityRoutingKey = "OpenActivityRoutingKey"

	RdmOrderExpiredKey = "OrderExpired:%s"
)
