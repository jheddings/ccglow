#!/usr/bin/env node

import { readFileSync, writeFileSync } from 'node:fs';
import { parseArgs } from './cli-parser.js';
import { run } from './runner.js';

const HELP = `Usage: ccnow [options]

Composable statusline for Claude Code.
Reads session JSON from stdin, outputs styled statusline to stdout.

Options:
  --preset <name>     Use a named preset (default, minimal, full)
  --config <path>     Load JSON config file
  --format <type>     Output format: ansi (default), plain
  --tee <path>        Write raw stdin JSON to file before processing
  --help              Show help
  --version           Show version

Examples:
  npx -y ccnow
  npx -y ccnow --preset=minimal
  npx -y ccnow --config ~/.claude/ccnow.json
`;

async function main(): Promise<void> {
  const args = parseArgs(process.argv.slice(2));

  if (args.help) {
    process.stdout.write(HELP);
    return;
  }

  if (args.version) {
    try {
      const pkg = JSON.parse(readFileSync(new URL('../package.json', import.meta.url), 'utf-8'));
      process.stdout.write(`${pkg.version}\n`);
    } catch {
      process.stdout.write('unknown\n');
    }
    return;
  }

  // Read stdin
  let stdin: string;
  try {
    stdin = readFileSync(0, 'utf-8');
  } catch {
    stdin = '';
  }

  // Tee: write raw stdin to file before processing
  if (args.tee) {
    try {
      writeFileSync(args.tee, stdin, 'utf-8');
    } catch (err) {
      process.stderr.write(`ccnow: failed to write tee file: ${err}\n`);
    }
  }

  const output = await run(args, stdin);
  if (output) process.stdout.write(output);
}

main().catch((err) => {
  process.stderr.write(`ccnow: ${err}\n`);
  process.exit(1);
});
