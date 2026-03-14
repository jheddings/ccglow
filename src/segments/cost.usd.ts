import type { Segment, SegmentContext } from '../types.js';
import type { CostData } from '../providers/cost.js';

export const costUsdSegment: Segment = {
  name: 'cost.usd',
  provider: 'cost',
  render(context: SegmentContext): string | null {
    const data = context.provider as CostData | undefined;
    return data?.usd ?? null;
  },
};
