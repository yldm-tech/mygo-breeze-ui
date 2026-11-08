import type { ReactNode } from "react";

type TokenKind = "comment" | "string" | "number" | "keyword" | "function" | "type" | "flag";

const goKeywords = new Set([
  "break", "case", "chan", "const", "continue", "default", "defer", "else", "fallthrough", "for", "func", "go", "goto", "if", "import", "interface", "map", "package", "range", "return", "select", "struct", "switch", "type", "var",
]);
const goConstants = new Set(["true", "false", "nil", "iota"]);
const shellCommands = new Set(["bun", "go", "npm", "npx", "pnpm", "yarn", "git", "curl", "docker", "make", "sudo"]);

const goTokenPattern = /(\/\/[^\n]*|"(?:\\.|[^"\\])*"|`[^`]*`|\b\d+(?:\.\d+)?\b|\b(?:break|case|chan|const|continue|default|defer|else|fallthrough|for|func|go|goto|if|import|interface|map|package|range|return|select|struct|switch|type|var)\b|\b(?:true|false|nil|iota)\b|\b[A-Z][A-Za-z0-9_]*\b|\b[A-Za-z_][A-Za-z0-9_]*(?=\s*\())/g;
const shellTokenPattern = /(?:#[^\n]*|"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|\B--?[A-Za-z][\w-]*|\b\d+(?:\.\d+)?\b|\b[A-Za-z][A-Za-z0-9_-]*\b)/g;

function kindOfGo(token: string): TokenKind | undefined {
  if (token.startsWith("//")) return "comment";
  if (token.startsWith('"') || token.startsWith("`")) return "string";
  if (/^\d/.test(token)) return "number";
  if (goKeywords.has(token)) return "keyword";
  if (goConstants.has(token)) return "number";
  if (/^[A-Z]/.test(token)) return "type";
  return "function";
}

function kindOfShell(token: string): TokenKind | undefined {
  if (token.startsWith("#")) return "comment";
  if (token.startsWith('"') || token.startsWith("'")) return "string";
  if (token.startsWith("-")) return "flag";
  if (/^\d/.test(token)) return "number";
  if (shellCommands.has(token)) return "function";
  return undefined;
}

function tokenClass(kind: TokenKind | undefined): string {
  return kind ? `syntax-${kind}` : "";
}

function highlight(code: string, pattern: RegExp, getKind: (token: string) => TokenKind | undefined): ReactNode[] {
  const nodes: ReactNode[] = [];
  let lastIndex = 0;
  let key = 0;
  for (const match of code.matchAll(pattern)) {
    const token = match[0];
    const index = match.index ?? 0;
    if (index > lastIndex) nodes.push(code.slice(lastIndex, index));
    const className = tokenClass(getKind(token));
    nodes.push(<span key={`${className}-${key++}`} className={className}>{token}</span>);
    lastIndex = index + token.length;
  }
  if (lastIndex < code.length) nodes.push(code.slice(lastIndex));
  return nodes;
}

export function highlightGo(code: string): ReactNode[] {
  return highlight(code, goTokenPattern, kindOfGo);
}

export function highlightShell(code: string): ReactNode[] {
  return highlight(code, shellTokenPattern, kindOfShell);
}

export function GoCode({ code, className = "" }: { code: string; className?: string }) {
  return <code className={className}>{highlightGo(code)}</code>;
}

export function ShellCode({ code, className = "" }: { code: string; className?: string }) {
  return <code className={className}>{highlightShell(code)}</code>;
}

export function GoCodeBlock({ code, className = "" }: { code: string; className?: string }) {
  return <pre className={className}><GoCode code={code} /></pre>;
}

export function ShellCodeBlock({ code, className = "" }: { code: string; className?: string }) {
  return <pre className={className}><ShellCode code={code} /></pre>;
}
