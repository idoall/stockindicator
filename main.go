package main

import (
	"fmt"
	"github.com/idoall/stockindicator/channel"
	"github.com/idoall/stockindicator/oscillator"
	"github.com/idoall/stockindicator/trend"
	"github.com/idoall/stockindicator/utils"
	"strings"
)

func main() {

	list := utils.GetTestKline()

	strategiesSides := utils.RunStrategies(
		trend.NewDefaultMacd(list),
		trend.NewDefaultDma(list),
		trend.NewDefaultTrix(list),
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
		channel.NewDefaultKeltnerChannel(list),
		trend.NewDefaultBbi(list),
		trend.NewDefaultTTMSqueeze(list),
	)

	for i, v := range list {
		var strategiesSidesMessages []string
		for _, sideMessage := range strategiesSides {
			strategiesSidesMessages = append(strategiesSidesMessages, fmt.Sprintf("%s:%s",
				sideMessage.Name,
				sideMessage.Data[i].PrintColorSide(),
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
