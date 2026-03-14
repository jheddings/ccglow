import type { SegmentNode } from '../types.js';
import { type BaseProps, type CompositeProps, extractStyle } from './primitives.js';

export function GitGroup(props: CompositeProps = {}): (children: () => SegmentNode[]) => SegmentNode {
  const { enabled, ...styleProps } = props;
  return (children) => ({
    type: 'git',
    enabled,
    style: extractStyle(styleProps),
    children: children(),
  });
}

export function GitBranch(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'git.branch', provider: 'git', enabled, style: extractStyle(styleProps) };
}

export function GitInsertions(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'git.insertions', provider: 'git', enabled, style: extractStyle(styleProps) };
}

export function GitDeletions(props: BaseProps = {}): SegmentNode {
  const { enabled, ...styleProps } = props;
  return { type: 'git.deletions', provider: 'git', enabled, style: extractStyle(styleProps) };
}
