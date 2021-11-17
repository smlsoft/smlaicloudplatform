package microservice


func MicroServiceHealthyHandler(ctx IServiceContext) error {
	ctx.ResponseS(200, "OK")
	return nil
}