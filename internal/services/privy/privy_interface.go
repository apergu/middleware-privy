package middleware

import request "middleware/infrastructure/http/request"

type PrivyToNetsuitService interface {
	ToNetsuit(req request.RequestToNetsuit) (interface{}, error)
}
