package microservice

func MicroServiceHealthyHandler(ctx IContext) error {
	ctx.ResponseS(200, "OK")
	return nil
}
