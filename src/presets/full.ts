import {
  StatusLine, Git, Group, Branch, Insertions, Deletions,
  Context, Tokens, Size, Percent, Literal,
  ModelName, CostUsd,
} from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const fullPreset: SegmentNode[] = StatusLine(() => [
  { type: 'pwd.smart', provider: 'pwd', style: { color: '31' } },
  { type: 'pwd.name', provider: 'pwd', style: { color: '39', bold: true } },
  Git({ prefix: ' | ', color: '240' })(() => [
    Branch({ color: 'whiteBright', bold: true, prefix: '\ue0a0 ' }),
    Group({ prefix: ' [', suffix: ']' })(() => [
      Insertions({ color: 'green', prefix: '+' }),
      Deletions({ color: 'red', prefix: ' -' }),
    ]),
  ]),
  { type: 'literal', props: { text: ' | ' }, style: { color: '240' } },
  ModelName({ color: '240' }),
  { type: 'literal', props: { text: ' · ' } },
  Context()(() => [
    Tokens({ color: 'white', bold: true }),
    Literal({ text: '/' }),
    Size({ color: 'white' }),
    Literal({ text: ' (' }),
    Percent({ color: 'white' }),
    Literal({ text: ')' }),
  ]),
  { type: 'literal', props: { text: ' · ' } },
  CostUsd({ color: 'yellow' }),
]);
