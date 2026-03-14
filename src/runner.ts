import { readFileSync } from 'node:fs';
import type { CliArgs } from './cli-parser.js';
import type { SegmentNode } from './types.js';
import { parseSession } from './session.js';
import { parseConfig } from './config.js';
import { SegmentRegistry } from './registry.js';
import { ProviderRegistry } from './providers.js';
import { registerBuiltinSegments } from './segments/index.js';
import { registerBuiltinProviders } from './providers/index.js';
import { renderTree } from './render.js';
import { getPreset } from './presets/index.js';
import { setColorLevel } from './style.js';

export async function run(args: CliArgs, stdin: string): Promise<string> {
  // Parse session
  const session = parseSession(stdin);
  if (!session) return '';

  // Set color level based on format
  setColorLevel(args.format === 'plain' ? 0 : 1);

  try {
    // Build registries
    const segmentRegistry = new SegmentRegistry();
    registerBuiltinSegments(segmentRegistry);

    const providerRegistry = new ProviderRegistry();
    registerBuiltinProviders(providerRegistry);

    // Resolve render tree: config file > preset
    let tree: SegmentNode[];

    if (args.config) {
      try {
        const configJson = JSON.parse(readFileSync(args.config, 'utf-8'));
        const configTree = parseConfig(configJson);
        tree = configTree.length > 0 ? configTree : getPreset(args.preset) ?? getPreset('default')!;
      } catch (err) {
        process.stderr.write(`ccnow: failed to load config: ${err}\n`);
        tree = getPreset(args.preset) ?? getPreset('default')!;
      }
    } else {
      tree = getPreset(args.preset) ?? getPreset('default')!;
    }

    // Resolve providers
    const providerNames = providerRegistry.collectProviderNames(tree);
    const providerData = await providerRegistry.resolveAll([...providerNames], session);

    // Render
    return renderTree(tree, segmentRegistry, session, providerData);
  } finally {
    setColorLevel(1);
  }
}
