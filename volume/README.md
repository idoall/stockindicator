# Volume 交易量相关交易指标


- [Accumulation Distribution Indicator](#accumulation-distribution-indicator)
- [Chaikin Money Flow](#chaikin-money-flow)
- [Ease of Movement](#ease-of-movement)
- [On Balance Volume(OBV)](#on-balance-volume)
- [Volume Price Trend(VPT)](#volume-price-trend)
- [Volume Weighted Moving Average(VWMA)](#volume-weighted-moving-average)


### Accumulation Distribution Indicator

Accumulation/Distribution Indicator (A/D)，吸筹/派发是通过价格和交易量的变化来判断的。在价格变化时交易量实际上作为权重系数 ― 系数 (交易量) 越高，价格变化 (在这段时间周期内) 对指标值的贡献越大。

事实上, 该指标是较常用的 能量潮 指标的一个版本。两者都通过衡量各自的业务交易量均价来确认价格的变化。

当吸筹/派发指标增长，意味着特定证券正在积累 (吸筹中)，因为业务量的压倒性份额与价格的上涨趋势相关。当指标下降时，意味着证券的派发 (抛售中)，因为大部分的业务是在价格下跌的过程中发生的。

吸筹/派发指标与证券价格之间的背离表明即将到来的价格变化。作为一条规则，通常情况下，如果出现这种背离，价格趋势会朝着指标运行的方向发展。 因此，如果指标增长，证券价格下降，应期待价格出现转折。

价格波动在上限和下限的区间之内，价格涨跌幅度加大时，带状区会变宽，涨跌幅度狭小盘整时，带状区会变窄。价格突破上限时，代表超买，价格跌穿下限时，代表超卖。

```golang
stock := NewDefaultAccumulationDistribution(list)

var dataList = stock.GetData()
```

### Chaikin Money Flow

Chaikin Money Flow (CMF) 蔡金资金流量是用于在一段时间内衡量资金流量的技术分析指标。

资金流量(Marc Chaikin创立的一个概念)是用于衡量单一期间证券的买卖压力的指标。

CMF在用户指定的回溯期内对资金流量进行加总。 任何回溯期都可使用，但最受欢迎的设定是20或21天。


 - 1、一般而言，CMF大于零，市场处于牛市，CMF小于零，市场处于熊市。
 - 2、CMF大于零（或小于零）的时间长短也值得注意。停留时间越长，趋势越可靠。
 - 3、CMF可以结合趋势线及支撑线、阻力线突破进行分析。
 - 4、CMF与价格之间的背离具有重要意义，通常预示着行情即将转变。

```golang
stock := NewDefaultChaikinMoneyFlow(list)

var dataList = stock.GetData()
```

### Ease of Movement
The Ease of Movement (EMV) 简易波动指标（Ease of Movement Value）又称EMV指标.

它是由RichardW．ArmJr．根据等量图和压缩图的原理设计而成,目的是将价格与成交量的变化结合成一个波动指标来反映股价或指数的变动状况。

由于股价的变化和成交量的变化都可以引发该指标数值的变动,因此,EMV实际上也是一个量价合成指标。


```golang
stock := NewDefaultKeltnerChannel(list)

var dataList = stock.GetData()

var side = stock.AnalysisSide()
```

### On Balance Volume
On Balance Volume (OBV) 是一种动量技术分析工具，OBV指标就是将成交量数据化，编制成趋势线，配合股价趋势，从量能角度判断股价走向。

 - 能量潮指标通过统计成交量变动的趋势来预测股价未来走向。
 - 当价格上涨时，成交量会被加入到运行总量中；当价格下跌时，成交量将从运行总量中减去。
 - 能量潮指标给交易者提供了检测成交量和股价是否背离的方法。

```golang
stock := NewObv(list)

var dataList = stock.GetData()
```

### Volume Price Trend
Volume Price Trend (VPT) 量价曲线是一种可以一定程度反映股价运动趋势的曲线，可以供投资者参考。

量价是不可分割的，量价曲线（VPT）将量能的增减和股价的涨跌结合起来分析，

从而确定量能的主要运动方向，进而得出主力资金进出的真实意图。投资者可以根据

这种量价的变化最终掌握股价的运动趋势。

```golang
stock := NewDefaultVolumePriceTrend(list)

var dataList = stock.GetData()
```

### Volume Weighted Moving Average
Volume Weighted Moving Average (VWMA) 基于交易量的一种加权移动平均指数（WMA）技术指标

```golang
stock := NewDefaultVwma(list)

var dataList = stock.GetData()
```
