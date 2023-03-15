package main

import (
	"fmt"
	"strings"

	"github.com/idoall/stockindicator/oscillator"
	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
)

func main() {

	list := utils.GetTestKline()

	strategiesSides := utils.RunStrategies(
		trend.NewMacd(list),
		trend.NewDefaultKdj(list),
		oscillator.NewAbsolutePriceOscillator(list),
		trend.NewDefaultCci(list),
		trend.NewDefaultRsi(list),
		oscillator.NewAwesomeOscillator(list),
		oscillator.NewDefaultChaikinOscillator(list),
		oscillator.NewDefaultIchimokuCloud(list),
		oscillator.NewDefaultPercentagePriceOscillator(list),
		oscillator.NewDefaultStochasticOscillator(list),
		oscillator.NewDefaultWilliamsR(list),
	)

	for i, v := range list {
		var strategiesSidesMessages []string
		for _, sideMessage := range strategiesSides {
			strategiesSidesMessages = append(strategiesSidesMessages, fmt.Sprintf("%s:%s",
				sideMessage.Name,
				getColorSide(sideMessage.Data[i]),
			))
		}
		fmt.Printf("[%d]Time:%s\tClose:%.2f\t %s\n",
			i+1,
			v.Time.Format("2006-01-02 15:04:05"),
			v.Close,
			strings.Join(strategiesSidesMessages, "\t"),
		)
	}

}

var (
	// greenBg      = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	// whiteBg      = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	// yellowBg     = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	// redBg        = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	// blueBg       = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	// magentaBg    = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	// cyanBg       = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	green = string([]byte{27, 91, 51, 50, 109})
	// white        = string([]byte{27, 91, 51, 55, 109})
	// yellow       = string([]byte{27, 91, 51, 51, 109})
	red = string([]byte{27, 91, 51, 49, 109})
	// blue         = string([]byte{27, 91, 51, 52, 109})
	// magenta      = string([]byte{27, 91, 51, 53, 109})
	// cyan         = string([]byte{27, 91, 51, 54, 109})
	reset = string([]byte{27, 91, 48, 109})
	// disableColor = false
)

func getColorSide(s utils.Side) string {
	switch s {
	case utils.Buy:
		return fmt.Sprintf("%s%s%s", green, s.String(), reset) // 绿色
	case utils.Sell:
		return fmt.Sprintf("%s%s%s", red, s.String(), reset) // 红色
	default:
		return s.String()
	}
}
