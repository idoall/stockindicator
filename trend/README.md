# Trend 趋势交易



- [Average True Range(ATR)](#average-true-range)
- [Commodity Channel Index(CCi)](#commodity-channel-index)
- [Double Exponential Moving Average(DEMA)](#double-exponential-moving-average)
- [Exponential Moving Average(EMA)](#exponential-moving-average)
- [Stochastic Indicator(KDJ)](#stochastic-indicator)
- [Moving Average(MA)](#moving-average)
- [Moving Average Convergence and Divergence(MACD)](#moving-average-convergence-and-divergence)
- [Relative Strength Index(RSI)](#relative-strength-index)
- [Simple Moving Average(SMA)](#simple-moving-average)
- [Vortex Indicator(Vortex)](#vortex-indicator)
- [Stochastic Relative Strength Index(Stoch RSI)](#stochastic-relative-strength-index)



### Average True Range

Average True Range (ATR)主要应用于了解价格的震荡幅度和节奏，在窄幅整理行情中用于寻找突破时机，但是不能直接用于预测价格走向。通常情况下价格的波动幅度会保持在一定常态下，但是如果有主力资金进出时，价格波幅往往会加剧。另外，在价格横盘整理，波幅减少到极点时，也往往会产生变盘行情。

```golang
stock := NewDefaultAtr(list)

var dataList = stock.GetData()

var long, short = stock.ChandelierExit(14)
```

### Commodity Channel Index

Commodity Channel Index(CCI)是专门衡量价格是否超出常态分布范围，强调价格平均绝对偏差在市场技术分析中的重要性，属于超买超卖类指标的一种。CCI指标却是波动于正无穷大到负无穷大之间，因此不会出现指标钝化现象，这样就有利于投资者更好地研判行情，特别是那些短期内暴涨暴跌的非常态行情。

```golang
stock := NewDefaultCci(list)

var dataList = stock.GetData()

var side = stock.AnalysisSide()
```

### Double Exponential Moving Average
Double Exponential Moving Average(DEMA)，双重指数移动平均线是由Patrick Mulloy开发的。在一篇名为《使用更快的移动平均线平滑数据》（Smoothing Data with Faster Moving Averages）的文章中，Patrick Mulloy首次提出了双重指数移动平均的概念。本文发表于1994年2月的《股票与商品期货技术分析》杂志中。

旨在减少传统移动平均线产生的结果的滞后性。技术交易者使用它来减少可能扭曲价格图表走势 的“噪音”量。

与任何移动平均线一样，DEMA 用于指示股票或其他资产价格的趋势。通过随时间跟踪其价格，交易者可以发现上升趋势（当价格高于其平均水平时）或下降趋势（当价格低于其平均水平时）。当价格穿越平均线时，可能预示着趋势的持续变化。

```golang
stock := NewDefaultDema(list)

var dataList = stock.GetData()
```

### Exponential Moving Average
Exponential Moving Average(EMA)，指数移动平均值。也叫 EXPMA 指标，是一种趋向类指标。其构造原理是：对收盘价进行加权算术平均，用于判断价格未来走势的变动趋势。与MACD指标、DMA指标相比，EMA指标由于其计算公式中着重考虑了当天价格（当期）行情的权重，决定了其作为一类趋势分析指标，在使用中克服了MACD指标对于价格走势的滞后性缺陷，同时，也在一定程度上消除了DMA指标在某些时候对于价格走势所产生的信号提前性，是一个非常有效的分析指标。

理解了 MA、EMA 的含义后，就可以理解其用途了，简单的说，当要比较数值与均价的关系时，用 MA 就可以了，而要比较均价的趋势快慢时，用 EMA 更稳定；有时，在均价值不重要时，也用 EMA 来平滑和美观曲线。

```golang
stock := NewDefaultEma(list)

var dataList = stock.GetData()
```

### Stochastic Indicator
Stochastic Indicator(KDJ)，KDJ指标的中文名称又叫随机指标，最早起源于期货市场，由乔治·莱恩（George Lane）首创。随机指标KDJ最早是以KD指标的形式出现，而KD指标是在威廉指标的基础上发展起来的。不过KD指标只判断股票的超买超卖的现象，在KDJ指标中则融合了移动平均线速度上的观念，形成比较准确的买卖信号依据。在实践中，K线与D线配合J线组成KDJ指标来使用。KDJ指标在设计过程中主要是研究最高价、最低价和收盘价之间的关系，同时也融合了动量观念、强弱指标和移动平均线的一些优点。因此，能够比较迅速、快捷、直观地研判行情，被广泛用于股市的中短期趋势分析，是期货和股票市场上最常用的技术分析工具。

```golang
stock := NewDefaultKdj(list)

var dataList = stock.GetData()
```



### Moving Average

Moving Average（MA）移动平均线是用统计分析的方法，将一定时期内的证券价格（指数）加以平均，并把不同时间的平均值连接起来，形成一根MA，用以观察证券价格变动趋势的一种技术指标。

移动平均线是由著名的美国投资专家Joseph E.Granville（葛兰碧，又译为格兰威尔）于20世纪中期提出来的。均线理论是当今应用最普遍的技术指标之一，它帮助交易者确认现有趋势、判断将出现的趋势、发现过度延生即将反转的趋势。

```golang
stock := NewDefaultMa(list)

var dataList = stock.GetData()
```

### Moving Average Convergence and Divergence
Moving Average Convergence / Divergence(MACD)，是Geral Appel 于1979年提出的，利用收盘价的短期（常用为12日）指数移动平均线与长期（常用为26日）指数移动平均线之间的聚合与分离状况，对买进、卖出时机作出研判的技术指标。

MACD在应用上应先行计算出快速（一般选12日）移动平均值与慢速（一般选26日）移动平均值。以这两个数值作为测量两者（快速与慢速线）间的“差离值”依据。所谓“差离值”（DIF），即12日EMA数值减去26日EMA数值。因此，在持续的涨势中，12日EMA在26日EMA之上。其间的正差离值（+DIF）会愈来愈大。反之在跌势中，差离值可能变负（-DIF），此时是绝对值愈来愈大。至于行情开始回转，正或负差离值要缩小到一定的程度，才真正是行情反转的信号。MACD的反转信号界定为“差离值”的9日移动平均值（9日DIF）。 

```golang
stock := NewMacd(list)

var dataList = stock.GetData()

var side = stock.AnalysisSide()
```



### Relative Strength Index

Relative Strength Index（RSI）RSI强弱指标(亦叫相对强弱指标)是利用一定时期内平均收盘涨数与平均收盘跌数的比值来反映股市走势的。“一定时期”选择不同，RSI选用天数可为5天，10天，14天。一般讲，天数选择短，易对起伏的股市产生动感，不易平衡长期投资的心理准备，做空做多的短期行为增多。天数选择长，对短期的投资机会不易把握。因此，参考5天、14天的RSI，是比较理智的。当然股民也可以自己选择更适合自己操作的天数。

```golang
stock := NewDefaultRsi(list)

var dataList = stock.GetData()

var side = stock.AnalysisSide()
```

### Simple Moving Average

Simple Moving Average（SMA）简单移动平均是指股票在特定时期内的平均收盘价。均线之所以被称为“移动”，是因为股票价格不断变化，所以移动均线也随之变化。

```golang
stock := NewDefaultSma(list)

var dataList = stock.GetData()
```

### Vortex Indicator

Vortex指标是Etienne Botes和Douglas Siepman发明的技术指标，用于识别金融市场中新趋势的开始或现有趋势的延续。它发表在2010年1月的《股票和商品技术分析》中。

它是一种趋势跟踪指标，用于发现趋势反转。 它具有两个振荡器，可以捕捉现有市场中趋势的正面和负面走势。 两条振荡线都随着市场移动。 当两条振荡线相互交叉时，交易决策基于该指标。 它既可以发出信号 翻转 或图表上的重大价格变动。

Vortex Indicator 线的计算取自指定时间内的特定价格高点和低点。 正振荡线是通过考虑近期低点和当前高点之间的范围来计算的，而负振荡线是通过考虑市场的最后高点和当前低点之间的范围来计算的。 对于图表中高度波动的市场或强劲的价格走势，VI（涡流指标）将显示每个振荡器之间的巨大差距或巨大距离。

```golang
stock := NewDefaultVortex(list)

var dataList = stock.GetData()
```

### Stochastic Relative Strength Index

Stoch RSI(Stochastic Relative Strength Index)结合了两种非常流行的技术分析指标，随机指标KDJ和相对强弱指标RSI。KDJ指标和RSI指标是基于价格计算得出的，而Stoch RSI是基于价格的RSI指标得出的；换句话说，就是把KDJ指标用在RSI指标上计算产生的。

應用法則

1. 超卖线交叉：Stoch RSI从超卖区穿越至超卖线(20)以上，这是一个买入信号。

2. 中线交叉：Stoch RSI从下往上穿过0.50线，可以用来确认其它技术指标产生的买入信号。

3. 底背离：在价格持续下跌的过程中，Stoch RSI却止跌回涨并穿越0.20线，形成底背离，这是一个买入信号。

4. 超买线交叉: Stoch RSI从超买区穿越至超买线(80) 以下，这是一个卖出信号。

5. 中线交叉：Stoch RSI从上往下穿过0.50线，可以用来确认其它技术指标产生的卖出信号。

6. 顶背离：在价格持续上涨的过程中，Stoch RSI却开始下跌并穿越0.80线，形成顶背离，这是一个卖出信号。

7. 在价格强力上涨和下跌的过程中，Stoch RSI和KDJ及RSI指标一样，会造成短期时间内的钝化现象。


```golang
stock := NewDefaultStochRsi(list)

var dataList = stock.GetData()
```