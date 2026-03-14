import type { SegmentNode } from '../types.js';
import { type BaseProps, type CompositeProps, extractStyle } from './primitives.js';

export function ContextGroup(props: CompositeProps = {}): (children: () => SegmentNode[]) => SegmentNode {
  const { enabled, ...styleProps } = props;
  return (children) => ({
    type: 'context',
    enabled,
    style: extractStyle(styleProps),
    children: children(),
  });
}

export function ContextTokens(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'context.tokens', provider: 'context', enabled, style: extractStyle(styleProps) };
}

export function ContextSize(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'context.size', provider: 'context', enabled, style: extractStyle(styleProps) };
}

export function ContextPercent(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'context.percent', provider: 'context', enabled, style: extractStyle(styleProps) };
}
