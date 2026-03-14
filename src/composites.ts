import {
  Pwd, Sep, Git, Branch, Insertions, Deletions,
  Context, Tokens, Percent, Literal,
} from './dsl/index.js';
import type { SegmentNode } from './types.js';

type CompositeBuilder = (sepChar: string) => SegmentNode;

const compositeBuilders: Record<string, CompositeBuilder> = {
  pwd: () => Pwd({ color: 'cyan', bold: true }),
  sep: (sepChar) => Sep({ char: sepChar, dim: true }),
  git: () => Git()(() => [
    Branch({ color: 'white', bold: true, icon: '\ue0a0 ' }),
    Literal({ text: ' [' }),
    Insertions({ color: 'green', prefix: '+' }),
    Literal({ text: ' ' }),
    Deletions({ color: 'red', prefix: '-' }),
    Literal({ text: ']' }),
  ]),
  context: () => Context()(() => [
    Literal({ text: 'ctx: ' }),
    Tokens({ bold: true }),
    Literal({ text: ' (' }),
    Percent(),
    Literal({ text: ')' }),
  ]),
};

export function buildCompositeTree(segments: string[], sepChar: string = '|'): SegmentNode[] {
  return segments
    .map((name) => compositeBuilders[name]?.(sepChar))
    .filter((node): node is SegmentNode => node !== undefined);
}
