package memcache

const NormalReaderSize = 4096

const (
	FlushAllLen = 1
	DelDelayLen = 3
	DelNowLen = 2
	StoreLen = 5
	GetLen = 2
)

const (
	CliErr = "CLIENT_ERROR"
	SrvErr = "SERVER_ERROR"
	CommonErr = "ERROR"
)

const (
	NotStored = "NOT_STORED"
	NotFound = "NOT_FOUND"
	Invaild = "invaild"
	Deleted = "DELETED"
	Stored = "STORED"
	Value = "VALUE"
	End = "END"
	OK = "OK"
	CRLF = "\r\n"
)

const (
//	DEL_EXPIR = "del_expir"
	FLUSH_ALL = "flush_all"
	REPLACE = "replace"
	DELETE = "delete"
	SET = "set"
	ADD = "add"
	GET = "get"
)
