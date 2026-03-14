import type { Segment, SegmentContext } from '../types.js';
import type { SessionMetrics } from '../providers/session.js';

export const sessionLinesRemovedSegment: Segment = {
  name: 'session.lines-removed',
  provider: 'session',
  render(context: SegmentContext): string | null {
    const data = context.provider as SessionMetrics | undefined;
    if (!data?.linesRemoved) return null;
    return `${data.linesRemoved}`;
  },
};
