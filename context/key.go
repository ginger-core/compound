package composition

type Key string

var (
	ctxDefaultKeys = []Key{
		"IP",
		"Agent",
		"ReqId",
	}
)

func AddDefaultKeys(keys ...Key) {
	ctxDefaultKeys = append(ctxDefaultKeys, keys...)
}
