// Custom mock for `react-inlinesvg` used in tests.
//
// The default mock scaffolded by `@grafana/create-plugin` only renders the
// `data-testid` derived from the icon file name and drops every other prop.
// Since `@grafana/ui` v13 the Select "Clear value" control is rendered via the
// `Icon` component (which forwards `role`/`aria-label` to `react-inlinesvg`),
// so dropping those props makes the control impossible to query by role/name in
// tests. This mock forwards the remaining props (e.g. `role`, `aria-label`,
// `title`) so accessible queries keep working, matching the real library's
// behaviour.

import React from 'react';

type Callback = (...args: any[]) => void;

export interface StorageItem {
  content: string;
  queue: Callback[];
  status: string;
}

export const cacheStore: { [key: string]: StorageItem } = Object.create(null);

const SVG_FILE_NAME_REGEX = /(.+)\/(.+)\.svg$/;

const InlineSVG = ({ src, ...rest }: { src: string; [key: string]: unknown }) => {
  // testId will be the file name without extension (e.g. `public/img/icons/angle-double-down.svg` -> `angle-double-down`)
  const testId = src.replace(SVG_FILE_NAME_REGEX, '$2');
  // `innerRef` and `onLoad`/`onError` are react-inlinesvg specifics that must not
  // be spread onto the DOM node.
  const { innerRef, onLoad, onError, loader, cacheRequests, preProcessor, uniquifyIDs, ...domProps } = rest as Record<
    string,
    unknown
  >;
  return <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" data-testid={testId} {...domProps} />;
};

export default InlineSVG;
