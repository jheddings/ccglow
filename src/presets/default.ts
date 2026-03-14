import {
  StatusLine, GitGroup, Group, GitBranch, GitInsertions, GitDeletions,
  ContextGroup, ContextTokens, ContextPercent, Literal,
} from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const defaultPreset: SegmentNode[] = StatusLine(() => [
  { type: 'pwd.smart', provider: 'pwd', style: { color: '31' } },
  { type: 'pwd.name', provider: 'pwd', style: { color: '39', bold: true } },
  GitGroup({ prefix: ' | ', color: '240' })(() => [
    GitBranch({ color: 'whiteBright', bold: true, prefix: '\ue0a0 ' }),
    Group({ prefix: ' [', suffix: ']' })(() => [
      GitInsertions({ color: 'green', prefix: '+' }),
      GitDeletions({ color: 'red', prefix: ' -' }),
    ]),
  ]),
  ContextGroup({ prefix: ' | ', color: '240' })(() => [
    Literal({ text: 'ctx: ' }),
    ContextTokens({ color: 'white', bold: true }),
    Literal({ text: ' (' }),
    ContextPercent({ color: 'white' }),
    Literal({ text: ')' }),
  ]),
]);
