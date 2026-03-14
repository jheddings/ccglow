import { Chalk } from 'chalk';
import type { StyleAttrs } from './types.js';

// Use 256-color level for broad terminal compatibility
let chalk = new Chalk({ level: 2 });

export function setColorLevel(level: 0 | 1 | 2 | 3): void {
  chalk = new Chalk({ level });
}

export function applyStyle(value: string, style: StyleAttrs | undefined): string {
  if (!style) return value;

  // Build the decorated string: prefix + value + suffix
  // Icon is prepended AFTER styling so it renders in default color.
  let result = value;
  if (style.prefix) result = style.prefix + result;
  if (style.suffix) result = result + style.suffix;

  // Apply chalk styling to the full string
  let painter: typeof chalk = chalk;

  // Apply modifiers before color — some terminals require bold
  // before color for bright white to render correctly.
  if (style.bold) painter = painter.bold;
  if (style.dim) painter = painter.dim;
  if (style.italic) painter = painter.italic;

  if (style.color) {
    // Support named colors and hex
    if (style.color.startsWith('#')) {
      painter = painter.hex(style.color);
    } else {
      painter = (painter as any)[style.color] ?? painter;
    }
  }

  // Only apply chalk if we actually set any style
  if (painter !== chalk) {
    // Apply chalk, then replace granular close sequences with a full
    // reset (\e[0m) for compatibility with statusline renderers that
    // don't handle partial ANSI resets correctly.
    result = painter(result).replace(/(\x1b\[\d+m)+$/, '\x1b[0m');
  }

  // Prepend icon outside styled region so it renders in default color
  if (style.icon) result = style.icon + result;

  return result;
}
