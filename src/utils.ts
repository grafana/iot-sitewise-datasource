import { DataFrame, MetricFindValue, DataFrameView } from "@grafana/data";

interface SimpleResult {
    id: string;
    name: string;
}

export function frameToMetricFindValues(df: DataFrame): MetricFindValue[] {
    const res:MetricFindValue[] = [];
    const view = new DataFrameView<SimpleResult>(df);
    view.forEach(item => {
        const id = item.id;
        const name = item.name;
        if (id && name) {
            res.push({text: name, value:id});
        }
        const value = id || name;
        if (value) {
            res.push({text: value, value});
        }
    });
    return res;
}