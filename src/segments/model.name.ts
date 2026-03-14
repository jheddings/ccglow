import type { Segment, SegmentContext } from '../types.js';
import type { ModelData } from '../providers/model.js';

export const modelNameSegment: Segment = {
  name: 'model.name',
  provider: 'model',
  render(context: SegmentContext): string | null {
    const data = context.provider as ModelData | undefined;
    return data?.name ?? null;
  },
};
