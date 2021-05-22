package aurafx

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	adminfx,
	configfx,
	debugfx,
	healthcheckfx,
	metricsfx,
	tracingfx,
)
