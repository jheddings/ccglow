import { describe, it, expect } from '@jest/globals';
import { pwdProvider } from '../../src/providers/pwd.js';
import type { SessionData } from '../../src/types.js';

describe('pwd provider', () => {
  it('resolves name, path, and smart from cwd', async () => {
    const session: SessionData = { cwd: '/Users/jheddings/Projects/ccnow' };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.name).toBe('ccnow');
    expect(data.path).toBe('/Users/jheddings/Projects/ccnow');
    expect(data.smart).toBeDefined();
  });

  it('handles root path', async () => {
    const session: SessionData = { cwd: '/' };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.name).toBe('/');
    expect(data.path).toBe('/');
  });

  it('smart keeps short paths as-is', async () => {
    const session: SessionData = { cwd: '/tmp' };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.smart).toBe('/tmp');
  });

  it('smart keeps 3-segment home paths as-is', async () => {
    const home = process.env.HOME ?? '/Users/test';
    const session: SessionData = { cwd: `${home}/Projects/ccnow` };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.smart).toBe('~/Projects/ccnow');
  });

  it('smart abbreviates 2 middle segments then ellipsis for deep paths', async () => {
    const home = process.env.HOME ?? '/Users/test';
    const session: SessionData = { cwd: `${home}/Projects/rise/red/app/routes/virtual-events` };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.smart).toBe('~/P/r/\u2026/virtual-events');
  });

  it('smart abbreviates 4-segment home path without ellipsis', async () => {
    const home = process.env.HOME ?? '/Users/test';
    const session: SessionData = { cwd: `${home}/Projects/rise/red` };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.smart).toBe('~/P/r/red');
  });

  it('smart handles 5-segment home path with ellipsis', async () => {
    const home = process.env.HOME ?? '/Users/test';
    const session: SessionData = { cwd: `${home}/Projects/rise/red/app` };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.smart).toBe('~/P/r/\u2026/app');
  });

  it('smart handles absolute non-home deep paths', async () => {
    const session: SessionData = { cwd: '/usr/local/share/some/deep/path' };
    const data = await pwdProvider.resolve(session) as any;
    expect(data.smart).toBe('/u/l/\u2026/path');
  });

  it('has correct name', () => {
    expect(pwdProvider.name).toBe('pwd');
  });
});
