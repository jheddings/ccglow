import type { DataProvider, SessionData } from '../types.js';

export interface CostData {
  usd: string | null;
}

function formatUsd(amount: number): string {
  return `$${amount.toFixed(2)}`;
}

export const costProvider: DataProvider = {
  name: 'cost',
  async resolve(session: SessionData): Promise<CostData> {
    const cost = session.cost;
    if (!cost) {
      return { usd: null };
    }

    return {
      usd: formatUsd(cost.total_cost_usd),
    };
  },
};
