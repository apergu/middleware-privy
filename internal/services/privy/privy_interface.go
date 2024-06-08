package middleware

type PrivyToNetsuitService interface {
	ToNetsuit(req, responseStruct interface{}, url, script, serviceName string) (interface{}, error)
}
