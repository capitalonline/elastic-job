package mcode

var (
	OK = add(200, "OK")

	NoLogin                = add(-101, "authentication expired")
	NotSupportClusterType  = add(-102, "not support clsuter type")
	NotEnoughShardAndMongo = add(-103, "the cluster don't have enough shard or mongos")
	ClusterHaveTask        = add(-104, "the cluster has running task")
	Panic                  = add(-105, "panic is commnig on")

	RequestErr         = add(-400, "请求错误")
	Unauthorized       = add(-401, "未认证")
	AccessDenied       = add(-403, "no authentication")
	NotFound           = add(-404, "404")
	MethodNotAllowed   = add(-405, "不支持该方法")
	Conflict           = add(-409, "冲突")
	ServerErr          = add(-500, "服务器错误")
	ServiceUnavailable = add(-503, "过载保护，服务暂时不可用")
	Deadline           = add(-504, "服务调用超时")
	LimitExceed        = add(-509, "超出限制")
	// 集群
	NotFoundCluster = add(-201, "the cluster has running task")

	//参数异常 以3xx开始
	DiskTypeParamExecption = add(-309, "not support disk type")
	NoMoney                = add(-310, "账户余额不足")
	OrderCalError          = add(-311, "订单计费异常")
	SiteError              = add(-312, "站点维护中")
	SoldOutError           = add(-313, "节点售罄")
	DataProductError       = add(-314, "数据异常")
	IPSoldOutError         = add(-315, "IP库存不足")

	//    SUCCESS = {'code': '0000', 'msg': 'success'}
	//    SUCCESS_CODE = SUCCESS['code']
	//    SUCCESS_MSG = SUCCESS['msg']
	//
	//    PARAM_ERROR = {'code': '20001', 'msg': '参数错误'}
	//    PARAM_ERROR_CODE = PARAM_ERROR['code']
	//    SITE_ERROR = {'code': '20002', 'msg': '节点维护中'}
	//    TASK_ERROR = {'code': '20003', 'msg': '任务失败'}
	//    SOLD_OUT = {'code': '20004', 'msg': '节点售罄'}
	//
	//    DATA_ERROR = {'code': '20101', 'msg': '数据异常'}
	//    DATA_ERROR_CODE = DATA_ERROR['code']
	//    DATA_PRO_ERROR = {'code': '20102', 'msg': '产品数据异常'}
	//    DATA_GOODS_ERROR = {'code': '20103', 'msg': '商品数据异常'}
	//    DATA_DUPLICATE_ERROR = {'code': '20104', 'msg': '数据重复'}
	//
	//    ORDER_LOCK_ERROR = {'code': '20201', 'msg': '订单冲突'}
	//    ORDER_LOCK_ERROR_CODE = ORDER_LOCK_ERROR['code']
	//    ORDER_CAL_ERROR = {'code': '20202', 'msg': '订单算价异常'}
	//    ORDER_EXPIRE_ERROR = {'code': '20203', 'msg': '订单过期'}
	//    ORDER_BILL_ERROR = {'code': '20204', 'msg': '订单计费异常'}
	//    ORDER_CANCEL_BILL_ERROR = {'code': '20205', 'msg': '订单取消计费异常'}
	//    REFUND_MONEY_ERROR = {'code': '20206', 'msg': '退费异常'}
	//
	//    ACCOUNT_NO_BALANCE = {'code': '20301', 'msg': '账户余额不足'}
	//    ACCOUNT_NO_BALANCE_CODE = ACCOUNT_NO_BALANCE['code']
	//    ACCOUNT_BALANCE_ERROR = {'code': '20302', 'msg': '查询账户余额异常'}
	//    ACCOUNT_BALANCE_ERROR_CODE = ACCOUNT_BALANCE_ERROR['code']
	//    ACCOUNT_AUTH_ERROR = {'code': '20303', 'msg': '用户权限不足'}
	//    ACCOUNT_REVIEW_UNPASS = {'code': '20304', 'msg': '账户审核未通过'}
	//    ACCOUNT_LOCK_ERROR = {'code': '20305', 'msg': '锁定账户金额异常'}
	//    ACCOUNT_SUBJECT_ERROR = {'code': '20306', 'msg': '测试项目异常'}
	//
	//    OPERATE_ERROR = {'code': '20401', 'msg': u"操作异常"}
	//    OPERATE_ERROR_CODE = OPERATE_ERROR['code']
	//
	//    IP_SOLD_OUT = {'code': '20601', 'msg': 'IP库存不足'}
	//
	//    BILL_ERROR = {'code': '20501', 'msg': 'redis发送计费异常'}
	//    BILL_ERROR_CODE = BILL_ERROR['code']
	//
	//    UNKNOWN_ERROR = {'code': '9999', 'msg': '未知异常'}
	//    UNKNOWN_ERROR_CODE = UNKNOWN_ERROR['code']

	//
	NoFindCluster    = add(-510, "cluster record not found")
	NoVmInfo         = add(-511, "node record not found")
	CreateOrderError = add(-512, "create order error")
	BillinfoError    = add(-513, "billing info error")

	// 备份 6xx
	NoAutoBack         = add(-613, "no auto backup")
	NotFoundBackup     = add(-614, "no backup")
	RemoveBackupFormS3 = add(-615, "remvove backup from s3")
	NofileofBackup     = add(-616, "the backup dont have files")

	// 恢复失败
	LogicalOplogTimeFormatErr        = add(-701, "logical recover oplog time errror")
	LogicalRecoveryTypeFormatErr     = add(-702, "logical recover oplog time errror")
	LogicalRecoveryDatabaseFormatErr = add(-703, "logical recover oplog time errror")
	RecoveryTempClusterFormatErr     = add(-704, "already exists temp cluster")

	//公网 8xx
	ExistPubnetInNode    = add(-801, "pubnet exist in this node")
	IpPoolNotEnough      = add(-802, "ip pool not enough")
	RuleInNodeNotExist   = add(-803, "rule not exist")
	CreateDnatRule       = add(-804, "access dnat service error")
	CluterPubNetInfoSync = add(-805, "update cluster public access sign error")
	DataSync             = add(-806, "data error")
	UnknowAction         = add(-807, "unknow action")
)
