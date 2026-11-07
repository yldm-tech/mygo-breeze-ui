import type { ReactNode } from "react";

type TokenKind = "comment" | "string" | "number" | "keyword" | "function";

const goKeywords = new Set([
  "break", "case", "chan", "const", "continue", "default", "defer", "else", "fallthrough", "for", "func", "go", "goto", "if", "import", "interface", "map", "package", "range", "return", "select", "struct", "switch", "type", "var",
]);

const goTokenPattern = /(\/\/[^\n]*|"(?:\\.|[^"\\])*"|`[^`]*`|\b\d+(?:\.\d+)?\b|\b(?:break|case|chan|const|continue|default|defer|else|fallthrough|for|func|go|goto|if|import|interface|map|package|range|return|select|struct|switch|type|var)\b)/g;

function kindOf(token: string): TokenKind | undefined {
  if (token.startsWith("//")) return "comment";
  if (token.startsWith('"') || token.startsWith("`")) return "string";
  if (/^\d/.test(token)) return "number";
  if (goKeywords.has(token)) return "keyword";
  return "function";
}

function highlightGo(code: string): ReactNode[] {
  const nodes: ReactNode[] = [];
  let lastIndex = 0;
  let key = 0;
  for (const match of code.matchAll(goTokenPattern)) {
    const token = match[0];
    const index = match.index ?? 0;
    if (index > lastIndex) nodes.push(code.slice(lastIndex, index));
    const kind = kindOf(token);
    const className = kind === "function" ? "fn" : kind === "string" ? "str" : kind === "number" ? "num" : kind === "keyword" ? "kw" : "comment";
    nodes.push(<span key={`${kind}-${key++}`} className={className}>{token}</span>);
    lastIndex = index + token.length;
  }
  if (lastIndex < code.length) nodes.push(code.slice(lastIndex));
  return nodes;
}

export function GoCodeBlock({ code, className = "" }: { code: string; className?: string }) {
  return <pre className={className}><code>{highlightGo(code)}</code></pre>;
}
