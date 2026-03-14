import type { DataProvider, SessionData } from '../types.js';

export interface ModelData {
  name: string | null;
}

export const modelProvider: DataProvider = {
  name: 'model',
  async resolve(session: SessionData): Promise<ModelData> {
    return {
      name: session.model?.display_name ?? null,
    };
  },
};
