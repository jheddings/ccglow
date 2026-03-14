import type { SegmentNode } from '../types.js';
import { type BaseProps, extractStyle } from './primitives.js';

interface PwdProps extends BaseProps {
  style?: 'name' | 'path' | 'smart';
}

export function Pwd(props: PwdProps = {}): SegmentNode {
  const { style: variant = 'smart', enabled, ...styleProps } = props;
  return {
    type: `pwd.${variant}`,
    provider: 'pwd',
    enabled,
    style: extractStyle(styleProps),
  };
}
