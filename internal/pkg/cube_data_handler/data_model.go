package cube_data_handler

import (
	"fmt"
	"strings"
)

type RealTimeMeasurements struct {
	AccX              float32
	AccY              float32
	AccZ              float32
	AccMax            float32
	VelX              float32
	VelY              float32
	VelZ              float32
	VelMax            float32
	T                 float32
	ShaftSpeed        float32
	VelDAUnbalanced   float32
	VelDAMisalignment float32
	VelDALooseness    float32
	VelDAOther        float32
}

type AnalogDigitalOutputs struct {
	U1Analog  float32
	U2Analog  float32
	U2Digital uint16
}

type DeviceStatuses struct {
	GoodOperatingMode        bool
	AnyHardwareFault         bool
	MEMSHardwareFault        bool
	EEPROMHardwareFault      bool
	ConfigurationFault       bool
	BootInProgress           bool
	DigitalOutputDriverFault bool
	InvalidSetup             bool
}

type ChannelStatuses struct {
	XOn         bool
	XInvalid    bool
	XFault      bool
	XSaturation bool
	YOn         bool
	YInvalid    bool
	YFault      bool
	YSaturation bool
	ZOn         bool
	ZInvalid    bool
	ZFault      bool
	ZSaturation bool
	TOn         bool
	TInvalid    bool
	TFault      bool
	TSaturation bool
}

type MeasurementsStatuses struct {
	AnyAlert          bool
	AnyVibrationAlert bool
	AnyAccAlert       bool
	AnyVelAlert       bool
	AnyXAlert         bool
	AnyYAlert         bool
	AnyZAlert         bool
	AccXAlert         bool
	AccYAlert         bool
	AccZAlert         bool
	VelXAlert         bool
	VelYAlert         bool
	VelZAlert         bool
	AccMaxAlert       bool
	VelMaxAlert       bool
	TAlert            bool
}

type dt struct {
	Unit  string
	Value float64
}

type VectorialMeasures struct {
	Dt      dt
	Unit    string
	ValuesX JSONasbleSlice
	ValuesY JSONasbleSlice
	ValuesZ JSONasbleSlice
}

type JSONasbleSlice []uint8

func (u JSONasbleSlice) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ", ")
	}
	return []byte(result), nil
}
