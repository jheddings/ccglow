import { describe, it, expect } from '@jest/globals';
import { applyStyle } from '../src/style.js';
import type { StyleAttrs } from '../src/types.js';

describe('applyStyle', () => {
  it('returns value unchanged when no style attrs', () => {
    expect(applyStyle('hello', {})).toBe('hello');
    expect(applyStyle('hello', undefined)).toBe('hello');
  });

  it('applies prefix before value', () => {
    const result = applyStyle('42', { prefix: '+' });
    expect(result).toContain('+');
    expect(result).toContain('42');
    expect(result.indexOf('+')).toBeLessThan(result.indexOf('42'));
  });

  it('applies suffix after value', () => {
    const result = applyStyle('42', { suffix: '%' });
    expect(result).toContain('42');
    expect(result).toContain('%');
  });

  it('applies icon before prefix and value', () => {
    const result = applyStyle('main', { icon: '\ue0a0 ', prefix: '' });
    expect(result).toContain('\ue0a0');
    expect(result).toContain('main');
  });

  it('applies color via chalk (ANSI codes present)', () => {
    const result = applyStyle('hello', { color: 'cyan' });
    // chalk wraps with ANSI escape codes
    expect(result).toContain('hello');
    expect(result.length).toBeGreaterThan('hello'.length);
  });

  it('applies bold via chalk', () => {
    const result = applyStyle('hello', { bold: true });
    expect(result).toContain('hello');
    expect(result.length).toBeGreaterThan('hello'.length);
  });

  it('combines multiple style attrs', () => {
    const style: StyleAttrs = { color: 'green', bold: true, prefix: '+' };
    const result = applyStyle('12', style);
    expect(result).toContain('+');
    expect(result).toContain('12');
  });

  it('handles null/undefined style fields gracefully', () => {
    const style: StyleAttrs = { color: undefined, bold: undefined };
    expect(applyStyle('hello', style)).toBe('hello');
  });
});
