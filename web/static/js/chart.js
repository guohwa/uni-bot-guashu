$(document).ready(function() {
    const up_arrow = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512"><path fill="#0ecb81" d="M374.6 246.6C368.4 252.9 360.2 256 352 256s-16.38-3.125-22.62-9.375L224 141.3V448c0 17.69-14.33 31.1-31.1 31.1S160 465.7 160 448V141.3L54.63 246.6c-12.5 12.5-32.75 12.5-45.25 0s-12.5-32.75 0-45.25l160-160c12.5-12.5 32.75-12.5 45.25 0l160 160C387.1 213.9 387.1 234.1 374.6 246.6z"/></svg>';
    const dn_arrow = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512"><path fill="#f6465d" d="M374.6 310.6l-160 160C208.4 476.9 200.2 480 192 480s-16.38-3.125-22.62-9.375l-160-160c-12.5-12.5-12.5-32.75 0-45.25s32.75-12.5 45.25 0L160 370.8V64c0-17.69 14.33-31.1 31.1-31.1S224 46.31 224 64v306.8l105.4-105.4c12.5-12.5 32.75-12.5 45.25 0S387.1 298.1 374.6 310.6z"/></svg>';

    var symbol = 'ETHUSDT';
    var interval = '15m';
    var limit = '1500';

    var price = 0;
    var nDeepth = 10;
    var ktime = 0;
    var moving = 0;

    var Storage = {
        save: function() {
            let key = `${symbol}_${interval}`;
            let s = JSON.stringify(levels);
            localStorage.setItem(key, s);
        },
        load: function() {
            let key = `${symbol}_${interval}`;
            let s = localStorage.getItem(key);
            if (s && s != 'undefined') {
                return JSON.parse(s);
            }
            return {};
        }
    }
    var levels = Storage.load();

    function xier(arr, callback) {
        let step = parseInt(arr.length / 2);
        while (step > 0) {
            for (let i = 0; i < arr.length; i++) {
                var n = i;
                while (typeof arr[n - step] != 'undefined' && callback(arr[n], arr[n - step])  && n > 0) {
                    var temp = arr[n];
                    arr[n] = arr[n - step];
                    arr[n - step] = temp;
                    n = n - step;
                }
            }
            step = parseInt(step / 2);
        }
        return arr;
    }

    var chartDom = document.getElementById('candles');
    var chart = LightweightCharts.createChart(chartDom, {
        crosshair: {
            mode: LightweightCharts.CrosshairMode.Normal
        },
        layout: {
            background: {
                color: '#161a25'
            },
            textColor: '#fff',
        },
        grid: {
            vertLines: {
                color: '#2b2b43'
            },
            horzLines: {
                color: '#2b2b43'
            }
        },
        timeScale: {
            timeVisible: true,
            secondsVisible: false
        },
        watermark: {
            visible: false
        },
    });

    var candleSeries = chart.addCandlestickSeries({
        scaleMargins: {
            top: 0,
            bottom: 0
        }
    });
    var candleData = [];

    var volumeSeries = chart.addHistogramSeries({
        priceScaleId: '',
        priceLineVisible: false,
        lastValueVisible: false,
        scaleMargins: {
            top: 0.8,
            bottom: 0,
        }
    });
    var volumeData = [];

    var TradeBook = {
        size: 60, // Cache size
        cache: [], // Cache
        update: function(t) {
            if (this.cache.length >= this.size) { // remove last if cache is full
                this.cache.pop();
                $('.trades table>tbody>tr:last-child').remove();
            }
            this.cache.unshift(t); // insert at first
            $('.trades table>tbody').prepend(`<tr><td${t.maker?' class="text-red"':' class="text-green"'}>${t.price}</td><td class="text-right">${t.quantity}</td><td class="text-right">${moment(t.time).format('HH:mm:ss')}</td></tr>`);
            $('.price span').removeClass('text-green').removeClass('text-red').html(price);
            $('.price svg').remove();
            if (t.price > price) { // bull
                $('.price span').addClass('text-green');
                $('.price').append(up_arrow);
            }
            if (t.price < price) { // bear
                $('.price span').addClass('text-red');
                $('.price').append(dn_arrow);
            }
            if (t.price != price) { // update current price
                price = t.price;
            }
        }
    }

    var OrderBook = {
        E: 0, //Event time
        T: 0, //Trade time
        u: 0, //Last update ID
        a: [], //Asks
        b: [], //Bids
        c: {
            E: 0,
            T: 0,
            u: 0,
            a: [],
            b: []
        }, //Cache for update
        d: 0, //Space per level
        t: 0, //Top value
        save: function(o) {
            this.E = o.E
            this.T = o.T;
            this.u = o.lastUpdateId;
            this.a = o.asks;
            this.b = o.bids;
            if (this.c.u > o.lastUpdateId) { // discard if expired
                this.update(this.c);
            }
        },
        cache: function(o) {
            this.c.E = o.E;
            this.c.T = o.T;
            this.c.u = o.u;

            let a = [];
            this.c.a.forEach(function(elem) {
                if (o.a.filter(x => x[0] == elem[0]).length == 0) {
                    a.push(elem);
                }
            });
            o.a.forEach(function(elem) {
                if (elem[1] != 0) {
                    a.push(elem);
                }
            });
            this.c.a = a;

            let b = [];
            this.c.b.forEach(function(elem) {
                if (o.b.filter(x => x[0] == elem[0]).length == 0) {
                    b.push(elem);
                }
            });
            o.b.forEach(function(elem) {
                if (elem[1] != 0) {
                    b.push(elem);
                }
            });
            this.c.b = b;
        },
        update: function(o) {
            if (this.u == 0) {
                this.cache(o); //Cache data when orderbook is not download by ajax
                return;
            }
            //????
            //While listening to the stream, each new event's pu should be equal to the previous event's u,
            //otherwise initialize the process from step 3.
            if (this.u > 0 && o.pu != this.u) {
                console.log(o.pu, this.u);
            }
            // Incremental update
            this.E = o.E;
            this.T = o.T;
            this.u = o.u;

            let a = [];
            this.a.forEach(function(elem) {
                if (o.a.filter(x => x[0] == elem[0]).length == 0) {
                    a.push(elem);
                }
            });
            o.a.forEach(function(elem) {
                if (elem[1] != 0) {
                    a.push(elem);
                }
            });
            this.a = a;

            let b = [];
            this.b.forEach(function(elem) {
                if (o.b.filter(x => x[0] == elem[0]).length == 0) {
                    b.push(elem);
                }
            });
            o.b.forEach(function(elem) {
                if (elem[1] != 0) {
                    b.push(elem);
                }
            });
            this.b = b;

            this.paint();
        },
        paint: async function() {
            if (price == 0 || this.a.length == 0) {
                return;
            }

            // Scale orderbook level by marks space
            let marks = candleSeries.priceScale().marks();
            if (marks.length < 2) {
                return;
            }

            let step = marks[0].label - marks[1].label;

            let p = this.a[0][0].split('.')[1].length; // Precision for price
            let q = this.a[0][1].split('.')[1].length; // Precision for size

            let base = Math.floor(price / step) * step; // orderbook(right list) middle
            let asks = [];
            for (let i = 0; i < nDeepth; i++) {
                let low = (i == 0 ? price : base + i * step);
                let high = base + (i + 1) * step;
                let size = 0.0;
                let sum = i == 0 ? 0.0 : Number(asks[i - 1].sum);

                this.a.forEach(function(elem) {
                    if (elem[0] >= low && elem[0] < high) {
                        size += (+elem[1]);
                        sum += (+elem[1]);
                    }
                });
                asks.push({
                    price: high.toFixed(p),
                    size: size.toFixed(q),
                    sum: sum.toFixed(q)
                });
            }
            $('.asks table>tbody').empty();
            asks.forEach(function(t) {
                $('.asks table>tbody').prepend(`<tr><td class="text-red">${t.price}</td><td class="text-right">${t.size}</td><td class="text-right">${t.sum}</td></tr>`);
            });

            let bids = [];
            for (let i = 0; i < nDeepth; i++) {
                let high = (i == 0 ? price : base - (i - 1) * step);
                let low = base - i * step;

                let size = 0.0;
                let sum = i == 0 ? 0.0 : Number(bids[i - 1].sum);

                this.b.forEach(function(elem) {
                    if (elem[0] > low && elem[0] <= high) {
                        size += (+elem[1]);
                        sum += (+elem[1]);
                    }
                });
                bids.push({
                    price: low.toFixed(p),
                    size: size.toFixed(q),
                    sum: sum.toFixed(q)
                });
            }

            $('.bids table>tbody').empty();
            bids.forEach(function(t) {
                $('.bids table>tbody').append(`<tr><td class="text-green">${t.price}</td><td class="text-right">${t.size}</td><td class="text-right">${t.sum}</td></tr>`);
            });
        }
    }

    // receive update of klines after download
    function getKlines() {
        $.ajax({
            url: `https://fapi.binance.com/fapi/v1/klines?symbol=${symbol}&interval=${interval}&limit=${limit}`,
            success: function(klines, textStatus) {
                klines.forEach(function(k, i) {
                    candleData.push({
                        time: k[0] / 1000,
                        open: +k[1],
                        high: +k[2],
                        low: +k[3],
                        close: +k[4],
                    });
                    volumeData.push({
                        time: k[0] / 1000,
                        value: +k[5],
                        color: k[4] > k[1] ? 'rgba(0, 150, 136, 0.6)' : 'rgba(255,82,82, 0.6)'
                    });
                });

                candleSeries.setData(candleData);
                volumeSeries.setData(volumeData);

                let candleConn = new WebSocket(`wss://fstream.binance.com/ws/${symbol.toLowerCase()}@kline_${interval}`);
                candleConn.onmessage = function(evt) {
					let j = JSON.parse(evt.data);
					if (!j || !j.k) {
                        console.log(evt);
                        return;
					}
                    let k = j.k;
                    let t = k.t / 1000;
                    if (t > ktime) {
                        ktime = t;
                    }
                    candleSeries.update({
                        time: t,
                        open: +k.o,
                        high: +k.h,
                        low: +k.l,
                        close: +k.c
                    });
                    volumeSeries.update({
                        time: t,
                        value: +k.v,
                        color: k.c > k.o ? 'rgba(0, 150, 136, 0.6)' : 'rgba(255,82,82, 0.6)'
                    });
                    if (moving == 0) {
                        let s = Number(k.close) > Number(k.open) ? 'text-green' : 'text-red';
                        $('.legend>.open').html(`open=<small class="${s}">${k.o}</small>`);
                        $('.legend>.high').html(`high=<small class="${s}">${k.h}</small>`);
                        $('.legend>.low').html(`low=<small class="${s}">${k.l}</small>`);
                        $('.legend>.close').html(`close=<small class="${s}">${k.c}</small>`);
                    }
                };
                candleConn.onclose = function(evt) {
                    console.log("Candle ws connection closed.");
                };
            }
        });
        return true;
    }

    // receive orderbook update by ws
    function getDeepth() {
        let deepthConn = new WebSocket(`wss://fstream.binance.com/ws/${symbol.toLowerCase()}@depth`);
        deepthConn.onmessage = function(evt) {
            let o = JSON.parse(evt.data);
            OrderBook.update(o);
            if (OrderBook.u == 0) {
                $.ajax({
                    url: `https://fapi.binance.com/fapi/v1/depth?symbol=${symbol.toLowerCase()}&limit=1000`,
                    success: function(data, textStatus) {
                        OrderBook.save(data);
                    }
                });
            }
        };
        deepthConn.onclose = function(evt) {
            console.log("Deepth ws connection closed.");
        };
        return true;
    }

    // receive update for aggTrade by ws
    function getTrades() {
        let tradeConn = new WebSocket(`wss://fstream.binance.com/ws/${symbol.toLowerCase()}@aggTrade`);
        tradeConn.onmessage = function(evt) {
            let t = JSON.parse(evt.data);
            TradeBook.update({
                time: t.T,
                price: t.p,
                quantity: t.q,
                maker: t.m
            });
        }
        tradeConn.onclose = function(evt) {
            console.log("Trade ws connection closed.");
        };
        return true;
    }

    function getPrices() {
        $.ajax({
            url: `https://fapi.binance.com/fapi/v1/ticker/price`,
            success: function(symbols, textStatus) {
                symbols = xier(symbols, function(a, b) {
                    return a.symbol < b.symbol;
                });
                symbols.forEach(function(s) {
                    $('.symbols table>tbody').append(`<tr data-symbol="${s.symbol}"><td>${s.symbol}</td><td class="text-right">${s.price}</td></tr>`);
                });
                let tickerConn = new WebSocket(`wss://fstream.binance.com/ws/!ticker@arr`);
                tickerConn.onmessage = function(evt) {
                    let t = JSON.parse(evt.data);
                    t.forEach(function(d) {
                        let s = Number(d.p) > 0 ? 'text-green' : 'text-red';
                        $(`.symbols table>tbody>tr[data-symbol=${d.s}]>td`).eq(1).removeClass('text-green').removeClass('text-red').addClass(s).html(d.c);
                    });
                }
                tickerConn.onclose = function(evt) {
                    console.log("Ticker ws connection closed.");
                };
                return true;
            }
        });
        return true;
    }

    // full repaint when browser's size changed
    window.onresize = function() {
        chart.resize(chartDom.clientWidth, chartDom.clientHeight, true);
    }

    // event of mouse click, can test code here
    chart.subscribeClick(function(param) {
        // if (!param.point) {
        //     return;
        // }

        // console.log(`Click at ${param.point.x}, ${param.point.y}. The time is ${param.time}.`);
    });

    // event of crosshair move
    chart.subscribeCrosshairMove(function(param) {
        if (param.time && param.time != ktime) {
            moving = 1;
            if (price > 0) {
                let p = price.split('.')[1].length;
                let v = param.seriesPrices.get(candleSeries);
                let s = v.close > v.open ? 'text-green' : 'text-red';
                $('.legend>.o').html(`open=<small class="${s}">${v.open.toFixed(p)}</small>`);
                $('.legend>.h').html(`high=<small class="${s}">${v.high.toFixed(p)}</small>`);
                $('.legend>.l').html(`low=<small class="${s}">${v.low.toFixed(p)}</small>`);
                $('.legend>.c').html(`close=<small class="${s}">${v.close.toFixed(p)}</small>`);
            }
        } else {
            moving = 0;
        }
    });

    function init(symbol, interval) {
        $('.legend>.symbol').html(symbol);
        $('.legend>.interval').html(interval);
        $('.sidebar .icon').click(function(evt) {
            let $that = $(this);
            let page = $that.data('page');
            let $active = $('.sidebar>.item.active>.icon');
            let active = $active.data('page');
            if (page != active) {
                $('.activity>.page').removeClass('active');
                $(`.activity>.page-${page}`).addClass('active');
                $('.sidebar>.item').removeClass('active');
                $that.parent('.item').addClass('active');
            }
        });
        return true;
    }

    init(symbol, interval) && getPrices() && getKlines() && getDeepth() && getTrades();
});
