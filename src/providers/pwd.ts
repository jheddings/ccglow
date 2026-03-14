import path from 'node:path';
import os from 'node:os';
import type { DataProvider, SessionData } from '../types.js';

export interface PwdData {
  name: string;
  path: string;
  smart: string;
}

const MAX_ABBREVIATED = 2;

function smartTruncate(cwd: string): string {
  const home = os.homedir();
  let p = cwd;

  // Replace home dir with ~
  if (p.startsWith(home)) {
    p = '~' + p.slice(home.length);
  }

  const parts = p.split('/');
  // 3 or fewer segments (e.g. ~/Projects/ccnow) — keep as-is
  if (parts.length <= 3) return p;

  const first = parts[0]; // '' for absolute, '~' for home-relative
  const last = parts[parts.length - 1];
  const middle = parts.slice(1, -1);

  // Abbreviate up to MAX_ABBREVIATED middle segments, then ellipsis if more remain
  if (middle.length <= MAX_ABBREVIATED) {
    // All middle segments fit as abbreviations (e.g. ~/P/r/red)
    const abbreviated = middle.map((part) => part[0] ?? '');
    return [first, ...abbreviated, last].join('/');
  }

  // More than MAX_ABBREVIATED middle segments — abbreviate first 2, then …
  const abbreviated = middle.slice(0, MAX_ABBREVIATED).map((part) => part[0] ?? '');
  return [first, ...abbreviated, '\u2026', last].join('/');
}

export const pwdProvider: DataProvider = {
  name: 'pwd',
  async resolve(session: SessionData): Promise<PwdData> {
    const cwd = session.cwd;
    return {
      name: cwd === '/' ? '/' : path.basename(cwd),
      path: cwd,
      smart: smartTruncate(cwd),
    };
  },
};
