import type { SegmentRegistry } from '../registry.js';
import { literalSegment } from './literal.js';
import { sepSegment } from './sep.js';
import { pwdNameSegment } from './pwd.name.js';
import { pwdPathSegment } from './pwd.path.js';
import { pwdSmartSegment } from './pwd.smart.js';
import { gitBranchSegment } from './git.branch.js';
import { gitInsertionsSegment } from './git.insertions.js';
import { gitDeletionsSegment } from './git.deletions.js';
import { contextTokensSegment } from './context.tokens.js';
import { contextSizeSegment } from './context.size.js';
import { contextPercentSegment } from './context.percent.js';
import { modelNameSegment } from './model.name.js';
import { costUsdSegment } from './cost.usd.js';
import { costDurationSegment } from './cost.duration.js';

export function registerBuiltinSegments(registry: SegmentRegistry): void {
  registry.register(literalSegment);
  registry.register(sepSegment);
  registry.register(pwdNameSegment);
  registry.register(pwdPathSegment);
  registry.register(pwdSmartSegment);
  registry.register(gitBranchSegment);
  registry.register(gitInsertionsSegment);
  registry.register(gitDeletionsSegment);
  registry.register(contextTokensSegment);
  registry.register(contextSizeSegment);
  registry.register(contextPercentSegment);
  registry.register(modelNameSegment);
  registry.register(costUsdSegment);
  registry.register(costDurationSegment);
}
