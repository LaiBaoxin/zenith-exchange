import { useEffect, useRef } from 'react';
import {
    createChart,
    ColorType,
    AreaSeries,
    CandlestickSeries,
    type Time,
    type IChartApi,
    type ISeriesApi,
} from 'lightweight-charts';

export type Period = '1m' | '5m' | '15m' | '1h' | '1d';

interface KLineChartProps {
    isDark: boolean;
    data?: any[];
    currentPeriod: Period;
    onPeriodChange: (newPeriod: Period) => void;
}

export const KLineChart = ({ isDark, data = [], currentPeriod, onPeriodChange }: KLineChartProps) => {
    const chartContainerRef = useRef<HTMLDivElement>(null);
    const chartRef = useRef<IChartApi | null>(null);
    const seriesRef = useRef<ISeriesApi<any> | null>(null);

    const formatData = (rawItems: any[]) => {
        if (!Array.isArray(rawItems) || rawItems.length === 0) return [];

        // 处理 ClickHouse 重复数据问题
        const uniqueMap = new Map<number, any>();
        rawItems.forEach(item => {
            const t = item.t || Math.floor(new Date(item.ts).getTime() / 1000);
            uniqueMap.set(t, item);
        });

        return Array.from(uniqueMap.values())
            .sort((a, b) => (a.t || 0) - (b.t || 0))
            .map(item => {
                const time = (item.t || Math.floor(new Date(item.ts).getTime() / 1000)) as Time;
                if (currentPeriod === '1m') {
                    return { time, value: Number(item.close) };
                }
                return {
                    time,
                    open: Number(item.open),
                    high: Number(item.high),
                    low: Number(item.low),
                    close: Number(item.close),
                };
            });
    };

    useEffect(() => {
        if (!chartContainerRef.current) return;

        const chart = createChart(chartContainerRef.current, {
            layout: {
                background: { type: ColorType.Solid, color: isDark ? '#141414' : '#ffffff' },
                textColor: isDark ? '#D9D9D9' : '#191919'
            },
            grid: {
                vertLines: { color: isDark ? '#2B2B2B' : '#F0F0F0' },
                horzLines: { color: isDark ? '#2B2B2B' : '#F0F0F0' }
            },
            width: chartContainerRef.current.clientWidth,
            height: 450,
            timeScale: { timeVisible: true, secondsVisible: false },
        });

        const series = currentPeriod === '1m'
            ? chart.addSeries(AreaSeries, { lineColor: '#2962FF', topColor: 'rgba(41, 98, 255, 0.3)', bottomColor: 'rgba(41, 98, 255, 0)', lineWidth: 2 })
            : chart.addSeries(CandlestickSeries, { upColor: '#26a69a', downColor: '#ef5350', wickUpColor: '#26a69a', wickDownColor: '#ef5350' });

        chartRef.current = chart;
        seriesRef.current = series;

        // 初始化数据
        if (data.length > 0) {
            series.setData(formatData(data));
        }

        const handleResize = () => {
            chart.applyOptions({ width: chartContainerRef.current?.clientWidth || 0 });
        };
        window.addEventListener('resize', handleResize);

        return () => {
            window.removeEventListener('resize', handleResize);
            chart.remove();
        };
    }, [isDark, currentPeriod]);

    // 监听 data 变化实现平滑跳动
    useEffect(() => {
        if (!seriesRef.current || data.length === 0) return;

        const formatted = formatData(data);
        if (formatted.length === 0) return;

        // 如果只有一条新数据，使用 update() 而不是 setData()
        // update 会让最后一根柱子实时波动，而不是闪烁重绘
        if (data.length === 1) {
            seriesRef.current.update(formatted[0]);
        } else {
            seriesRef.current.setData(formatted);
        }
    }, [data]);

    return (
        <div style={{ position: 'relative', background: isDark ? '#141414' : '#fff' }}>
            <div style={{ padding: '8px 12px', display: 'flex', gap: '8px', borderBottom: `1px solid ${isDark ? '#2b2b2b' : '#f0f0f0'}` }}>
                {(['1m', '5m', '15m', '1h', '1d'] as Period[]).map(p => (
                    <button key={p} onClick={() => onPeriodChange(p)} style={{
                        padding: '2px 10px', cursor: 'pointer', borderRadius: '4px', border: 'none',
                        backgroundColor: currentPeriod === p ? '#2962FF' : (isDark ? '#333' : '#eee'),
                        color: currentPeriod === p ? '#fff' : (isDark ? '#999' : '#666'), fontSize: '11px'
                    }}>
                        {p === '1m' ? '分时' : p.toUpperCase()}
                    </button>
                ))}
            </div>
            <div ref={chartContainerRef} style={{ height: '450px' }} />
        </div>
    );
};
