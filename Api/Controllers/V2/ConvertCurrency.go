package V2

import (
	"Currencies/Api/Controllers"
)

func (v *V2Controllers) ConvertCurrency(ctx Controllers.RequestContextImpl) {
	convertedResponse, err := v.Convert(ctx)

	if err != nil {
		err.SendError()
		return
	}

	ctx.Json(convertedResponse)
}
