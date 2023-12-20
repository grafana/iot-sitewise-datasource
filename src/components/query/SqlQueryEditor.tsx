import {
    Field,
    TextArea
} from '@grafana/ui';

import React, { useCallback, FC, FormEvent } from 'react';
import { SitewiseQueryEditorProps } from './types';
import { QueryType, SqlQuery } from 'types';

const QUERY_REGEX = /(\s*([\0\b\'\"\n\r\t\%\_\\]*\s*(((select\s+\S.*\s+from\s+\S+)|(insert\s+into\s+\S+)|(update\s+\S+\s+set\s+\S+)|(delete\s+from\s+\S+)|(((drop)|(create)|(alter)|(backup))\s+((table)|(index)|(function)|(PROCEDURE)|(ROUTINE)|(SCHEMA)|(TRIGGER)|(USER)|(VIEW))\s+\S+)|(truncate\s+table\s+\S+)|(exec\s+)|(\/\*)|(--)))(\s*[\;]\s*)*)+)/i;

interface SqlQueryEditorProps extends SitewiseQueryEditorProps<SqlQuery> {}

export const SqlQueryEditor: FC<SqlQueryEditorProps> = ({ query, onChange, onRunQuery }) => {
    const onChangeHandler = useCallback(({ currentTarget }: FormEvent<HTMLTextAreaElement>) => {
        const { value } = currentTarget;
        onChange({
            queryStatement: value,
            queryType: QueryType.SQL,
            refId: ''
        });
    }, [onChange]);

    return (
        <div>
            <Field label="Query Statement">
                <TextArea invalid={!QUERY_REGEX.test(query.queryStatement)} onChange={onChangeHandler} value={query.queryStatement} onBlur={onRunQuery} />
            </Field>
        </div>
    );
};
