## Oscillator Strategies

震荡器相关的策略.

- [Absolute Price Oscillator](#absolute-price-oscillator)
- [Awesome Oscillator](#awesome-oscillator)
- [Chaikin Oscillator](#chaikin-oscillator)
- [Ichimoku Cloud](#ichimoku-cloud)
- [Percentage Price Oscillator](#percentage-price-oscillator)
- [Projection Oscillator](#projection-oscillator)
- [Stochastic Oscillator](#stochastic-oscillator)
- [Williams R](#williams-r)
- [Volume Oscillator](#volume-oscillator)

### Absolute Price Oscillator

Absolute Price Oscillator(APO) 也叫绝对波动指数. 用于跟踪趋势的技术指标。APO 上穿零表示看涨，而下穿零表示看跌。正值表示上升趋势，负值表示下降趋势。

```golang
stock := NewAbsolutePriceOscillator(klineList)

var dataList = stock.GetData()
var side = stock.AnalysisSide()
```

### Awesome Oscillator

Awesome Oscillator简称AO指标，也叫动量震荡指标。它的发明者是Bill Williams，AO指标是用于显示当前市场的发展趋势，而以其柱形图的形式表现出来。AO 计算34个周期和5个周期简单移动平均的差。使用的简单移动平均不是使用收盘价计算的，而是每个柱的中点价格。AO通常被用来确认趋势或预期可能的逆转。

```golang
stock := NewAwesomeOscillator(klineList)

var dataList = stock.GetData()
var side = stock.AnalysisSide()
```

### Chaikin Oscillator

Chaikin Oscillator 柴金震荡指标是一项应用 MACD 公式来评估另一项被称为累积/派发线（Accumulation-Distribution Line，ADL）的指标动量的指标。这项指标得名于它的创造者马克·柴金（Marc Chaikin）。通过帮您识别动量可能会发生变化的点位，柴金震荡指标能让您知道趋势何时会逆转。

```golang
stock := NewDefaultChaikinOscillator(klineList)

var dataList = stock.GetData()
var side = stock.AnalysisSide()
```

### Ichimoku Cloud

Ichimoku Cloud. 也称为 Ichimoku Kinko Hyo，计算一个多功能指标，定义支撑和阻力，识别趋势方向，衡量动量，并提供交易信号。

经过三十多年的研究和测试，Goichi Hosada认为（9,26,52）的周期设置效果最佳。当时，日本的商业时间表中包括星期六，所以9代表一周半（6 + 3天）的时间。数字26和52分别代表一个月和两个月的时间。

在加密货币市场中，许多交易者通常将Ichimoku的周期范围设为（9,26,52）到（10,30,60），以此适应7*24小时的市场。甚至还可以直接将周期设置为（20,60,120），以减少错误信号的产生。

```golang
stock := NewDefaultIchimokuCloud(klineList)

var dataList = stock.GetData()
var side = stock.AnalysisSide()
```

### Percentage Price Oscillator

Percentage Price Oscillator (PPO). 百分比价格振荡器,是一个动量技术指标，它表示动量方向作为振荡器的迹象

虽然它确实与一些更受欢迎的振荡器（如 MACD）有一些相似之处, 这也是相当奇特的，因为它使用价格的百分比变化来计算动量，而不是绝对价格.

正 PPO 线表示看涨趋势, 而负 PPO 线表示看跌趋势.

还可以根据 PPO 线和信号线的交互方式确定动量. 每当 PPO 线越过信号线上方时，动量是看涨的, 当PPPO线越过信号线下方时，看跌.

```golang
stock := NewDefaultPercentagePriceOscillator(klineList)

var dataList = stock.GetData()
var side = stock.AnalysisSide()
```

### Projection Oscillator

Percentage Price Oscillator (PPO). 由Dr. Mel Widner 研仓。

与其他不同的指标一样，传统的用法也是超买/超卖，背驰，突破等，不少人利用此指标来交易外汇。

传统的参数是12及5，但若应用在港股上，将参数设定为50及10会更好。

分析股票时，初步看，每当由50以下重回至50以上有机会是股价重拾升势的时间，值得留意，不过有关的方法仍有待详细测试。

```golang
stock := NewDefaultProjectionOscillator(klineList)

var dataList = stock.GetData()
var side = stock.AnalysisSide()
```

### Stochastic Oscillator

Stochastic Oscillator. 显示了一定时期内收盘价相对于高低区间的位置。

随机震荡指标不跟随价格，也不跟随交易量或类似的东西。它跟随价格的速度或动量。通常，动量在价格之前改变方向。

因此，随机震荡指标的看涨和看跌背离可用于预示反转。

还可以使用这个振荡器来识别牛市和熊市的设置，以预测未来的逆转。

由于随机震荡指标是区间震荡指标（range-bound），因此它对于识别超买和超卖水平也很有用。

```golang
stock := NewDefaultStochasticOscillator(klineList)

var dataList = stock.GetData()
var side = stock.AnalysisSide()
```

### Williams R

Williams R. 是由美国作家、股市投资家拉里·威廉斯（Larry R. Williams）在1973年出版的《我如何赚得一百万》一书中首先发表，这个指标是一个振荡指标，是依股价的摆动点来度量股票／指数是否处于超买或超卖的现象。

它衡量多空双方创出的峰值（最高价）距每天收市价的距离与一定时间内（如7天、14天、28天等）的股价波动范围的比例，以提供出股市趋势反转的讯号。

根据最低价、最高价和收盘价计算 Williams R。它是一种动量指标，在 0 和 -100 之间移动，衡量超买和超卖水平。

```golang
stock := NewDefaultWilliamsR(klineList)

var dataList = stock.GetData()
var side = stock.AnalysisSide()
```

### Volume Oscillator

Volume Oscillator. 又名移动平均成交量指标，但是，它并非仅仅计算成交量的移动平均线，

而是通过对成交量的长期移动平均线和短期移动平均线之间的比较

分析成交量的运行趋势和及时研判趋势转变方向

```golang
stock := TestVolumeOscillator(klineList)

var dataList = stock.GetData()
```