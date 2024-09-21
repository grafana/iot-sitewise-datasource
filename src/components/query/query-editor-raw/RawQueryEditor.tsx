// import defaults from 'lodash/defaults';
import React from 'react';
import { SQLEditor } from '@grafana/experimental'; // LanguageDefinition
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from 'SitewiseDataSource';
import { SitewiseQuery, SitewiseOptions } from 'types'; // , QueryType, ListAssetsQuery, ListTimeSeriesQuery

type Props = QueryEditorProps<DataSource, SitewiseQuery, SitewiseOptions>;

export const firstLabelWith = 20;

export function RawQueryEditor(props: Props) {
  // const queryRef = useRef<SitewiseRawQuery>(query);
  // useEffect(() => {
  //     queryRef.current = query;
  // }, [query]);

  // const onRawQueryChange = useCallback(
  //     (rawSql: string, processQuery: boolean) => {
  //         const newQuery = {
  //             ...queryRef.current,
  //             rawQuery: true,
  //             rawSql,
  //         };
  //         onChange(newQuery, processQuery);
  //     },
  //     [onChange]
  // );

  // const onQueryTextChange = (event: ChangeEvent<HTMLInputElement>) => {
  //     onChange({ ...query, queryText: event.target.value });
  // };

  return (
    <SQLEditor
      query="select * from sitewise"
      // onChange={onQueryTextChange}
      // language={editorLanguageDefinition}
    />
    //     {children}
    // </SQLEditor>
  );
}
