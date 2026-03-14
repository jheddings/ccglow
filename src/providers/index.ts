import type { ProviderRegistry } from '../providers.js';
import { pwdProvider } from './pwd.js';
import { gitProvider } from './git.js';
import { contextProvider } from './context.js';
import { modelProvider } from './model.js';
import { costProvider } from './cost.js';
import { sessionProvider } from './session.js';

export function registerBuiltinProviders(registry: ProviderRegistry): void {
  registry.register(pwdProvider);
  registry.register(gitProvider);
  registry.register(contextProvider);
  registry.register(modelProvider);
  registry.register(costProvider);
  registry.register(sessionProvider);
}
