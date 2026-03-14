import type { DataProvider, SessionData } from '../types.js';

export interface CostData {
  usd: string | null;
  duration: string | null;
}

function formatUsd(amount: number): string {
  return `$${amount.toFixed(2)}`;
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

export const costProvider: DataProvider = {
  name: 'cost',
  async resolve(session: SessionData): Promise<CostData> {
    const cost = session.cost;
    if (!cost) {
      return { usd: null, duration: null };
    }

    return {
      usd: formatUsd(cost.total_cost_usd),
      duration: formatDuration(cost.total_duration_ms),
    };
  },
};
