import { useState } from "react";
import { ArrowRight, Check, Copy } from "lucide-react";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

const commands = ["go get github.com/yldm-tech/mygo-breeze-ui", "go run ./examples/breeze-ui"];
function Command({ value }: { value: string }) {
  const [copied, setCopied] = useState(false);
  return (
    <div className="flex items-center justify-between gap-3 border-b border-line px-4 py-3 last:border-b-0">
      <code className="min-w-0 overflow-hidden text-ellipsis whitespace-nowrap font-mono text-[10px] text-muted">
        {value}
      </code>
      <button
        className="inline-flex shrink-0 items-center gap-1.5 bg-transparent text-[10px] text-dim hover:text-mint"
        onClick={async () => {
          await navigator.clipboard?.writeText(value);
          setCopied(true);
          setTimeout(() => setCopied(false), 1200);
        }}
      >
        {copied ? (
          <>
            <Check size={13} /> Copied
          </>
        ) : (
          <>
            <Copy size={13} /> Copy
          </>
        )}
      </button>
    </div>
  );
}

export function InstallSection() {
  const { locale } = useUiStore();
  const t = useCopy(locale);
  return (
    <section className="px-0 pb-[105px] max-[680px]:pb-20" id="install">
      <div className="mx-auto grid w-[calc(100%-48px)] max-w-[1160px] grid-cols-[0.85fr_1.15fr] items-center gap-[70px] max-[940px]:gap-[45px] max-[680px]:grid-cols-1 max-[680px]:gap-[35px] max-[680px]:w-[calc(100%-34px)]">
        <div>
          <p className="mb-5 font-mono text-[10px] uppercase tracking-[0.09em] text-mint">
            {t.installKicker}
          </p>
          <h2 className="mb-3.5 text-[clamp(29px,3.5vw,44px)] leading-[1.1] tracking-[-0.065em]">
            {t.installTitle}
          </h2>
          <p className="mb-5 text-[13px] text-muted">{t.installText}</p>
          <a
            className="inline-flex items-center gap-1.5 text-xs font-extrabold text-mint hover:text-ink"
            href="https://github.com/yldm-tech/mygo-breeze-ui/blob/main/docs/breeze-ui.md"
            target="_blank"
            rel="noreferrer"
          >
            {t.docs}
            <ArrowRight size={15} />
          </a>
        </div>
        <div className="overflow-hidden rounded-lg border border-line bg-panel">
          <div className="flex gap-5 border-b border-line px-4 py-3 font-mono text-[10px] text-dim">
            <span className="text-mint">Go</span>
            <span>CLI</span>
            <span>MCP</span>
          </div>
          {commands.map((command) => (
            <Command key={command} value={command} />
          ))}
        </div>
      </div>
    </section>
  );
}
