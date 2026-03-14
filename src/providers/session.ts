import type { DataProvider, SessionData } from '../types.js';

export interface SessionMetrics {
  duration: string | null;
  linesAdded: number | null;
  linesRemoved: number | null;
}

function formatDuration(ms: number): string {
  const totalSec = Math.floor(ms / 1000);
  const hours = Math.floor(totalSec / 3600);
  const mins = Math.floor((totalSec % 3600) / 60);

  if (hours > 0) {
    return `${hours}h ${mins}m`;
  }
  return `${mins}m`;
}

export const sessionProvider: DataProvider = {
  name: 'session',
  async resolve(session: SessionData): Promise<SessionMetrics> {
    const cost = session.cost;
    if (!cost) {
      return { duration: null, linesAdded: null, linesRemoved: null };
    }

    return {
      duration: formatDuration(cost.total_duration_ms),
      linesAdded: cost.total_lines_added ?? null,
      linesRemoved: cost.total_lines_removed ?? null,
    };
  },
};
