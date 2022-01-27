package Controllers

type ApiServiceControllers interface {
	ListAllConversions(ctx RequestContextImpl)
	SingleConversion(ctx RequestContextImpl)
	ConvertCurrency(ctx RequestContextImpl)
}
