import { StatusLine, GitGroup, GitBranch, ContextGroup, ContextTokens, Literal, ContextSize } from '../dsl/index.js';
import type { SegmentNode } from '../types.js';

export const minimalPreset: SegmentNode[] = StatusLine(() => [
  { type: 'pwd.name', provider: 'pwd', style: { color: '39' } },
  GitGroup({ prefix: ' | ', color: '240' })(() => [
    GitBranch({ color: 'whiteBright', bold: true }),
  ]),
  ContextGroup({ prefix: ' | ', color: '240' })(() => [
    ContextTokens({ color: 'white' }),
    Literal({ text: '/' }),
    ContextSize({ color: 'white' }),
  ]),
]);
