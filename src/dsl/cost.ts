import type { SegmentNode } from '../types.js';
import { type BaseProps, extractStyle } from './primitives.js';

export function CostUSD(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'cost.usd', provider: 'cost', enabled, style: extractStyle(styleProps) };
}
