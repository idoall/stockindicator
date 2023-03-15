# Channel 通道类交易指标


- [Bollinger Bands(Boll)](#bollinger-bands)
- [Donchian Channel](#donchian-channel)
- [Keltner Channel](#keltner-channel)
- [Ulcer Index](#ulcer-index)


### Bollinger Bands

Bollinger Bands(BOLL)通过计算价格的“标准差”，再求价格的“信赖区间”，是一个路径型指标，该指标在图形上画出三条线，上下两条线可以分别看成是价格的压力线和支撑线，中间为价格平均线。

价格波动在上限和下限的区间之内，价格涨跌幅度加大时，带状区会变宽，涨跌幅度狭小盘整时，带状区会变窄。价格突破上限时，代表超买，价格跌穿下限时，代表超卖。

```golang
stock := NewDefaultBoll(list)

var dataList = stock.GetData()

var side = stock.AnalysisSide()
```

### Donchian Channel

Donchian 唐奇安通道指标，由“潮流之父”Richard Donchian 在二十世纪中叶开发。

它们是由移动平均线计算生成的三条线，包括由中间范围附近的上限和下限形成的指标。

但是，它们之间的区域是为该周期选择的通道。 这是为了帮助他识别市场交易的趋势。

因此，它在大多数交易平台上都很常见。

```golang
stock := NewDefaultDonchianChannel(list)

var dataList = stock.GetData()

var side = stock.AnalysisSide()
```

### Keltner Channel
Keltner Channels是一个波动性指标，由一位名叫 Chester Keltner 的交易商在他 1960 年的著作《如何在商品中赚钱》中引入。

Linda 版本的 Keltner Channel 使用更广泛，它与布林带非常相似，因为它也由三条线组成。

由于该通道源自ATR，而ATR本身就是一个波动率指标，因此Keltner 通道也会随着波动率收缩和扩张，但不像布林带那样波动。


```golang
stock := NewDefaultKeltnerChannel(list)

var dataList = stock.GetData()

var side = stock.AnalysisSide()
```

### Ulcer Index
Ulcer Index (UI) 溃疡指数是一种技术指标，可根据价格下跌的深度和持续时间来衡量下行风险。

该指数随着价格远离近期高点而增加，并随着价格升至新高而下跌。

该指标通常在 14 天内计算，溃疡指数显示交易者可以预期从该时期的高点回撤的百分比。

```golang
stock := NewDefaultUlcerIndex(list)

var dataList = stock.GetData()
```
