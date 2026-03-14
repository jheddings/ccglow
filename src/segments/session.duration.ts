import type { Segment, SegmentContext } from '../types.js';
import type { SessionMetrics } from '../providers/session.js';

export const sessionDurationSegment: Segment = {
  name: 'session.duration',
  provider: 'session',
  render(context: SegmentContext): string | null {
    const data = context.provider as SessionMetrics | undefined;
    return data?.duration ?? null;
  },
};
