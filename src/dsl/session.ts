import type { SegmentNode } from '../types.js';
import { type BaseProps, extractStyle } from './primitives.js';

export function SessionDuration(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return {
    type: 'session.duration',
    provider: 'session',
    enabled,
    style: extractStyle(styleProps),
  };
}

export function SessionLinesAdded(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return {
    type: 'session.lines-added',
    provider: 'session',
    enabled,
    style: extractStyle(styleProps),
  };
}

export function SessionLinesRemoved(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return {
    type: 'session.lines-removed',
    provider: 'session',
    enabled,
    style: extractStyle(styleProps),
  };
}
