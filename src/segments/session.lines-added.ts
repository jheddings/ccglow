import type { Segment, SegmentContext } from '../types.js';
import type { SessionMetrics } from '../providers/session.js';

export const sessionLinesAddedSegment: Segment = {
  name: 'session.lines-added',
  provider: 'session',
  render(context: SegmentContext): string | null {
    const data = context.provider as SessionMetrics | undefined;
    if (!data?.linesAdded) return null;
    return `${data.linesAdded}`;
  },
};
