import type { SegmentNode } from '../types.js';
import { type BaseProps, extractStyle } from './primitives.js';

export function ModelName(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'model.name', provider: 'model', enabled, style: extractStyle(styleProps) };
}
