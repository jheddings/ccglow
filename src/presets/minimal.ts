import { StatusLine, Sep, Branch } from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const minimalPreset: SegmentNode[] = StatusLine(() => [
  { type: 'pwd.name', provider: 'pwd', style: { color: 'cyanBright', bold: true } },
  Sep({ char: '|', color: '240' }),
  Branch({ color: 'whiteBright', bold: true, icon: '\ue0a0 ' }),
]);
