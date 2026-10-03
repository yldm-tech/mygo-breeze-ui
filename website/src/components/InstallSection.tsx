import { useState } from "react";
import { ArrowRight, Check, Copy } from "lucide-react";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

const commands = ["go get github.com/yldm-tech/mygo-breeze-ui", "go run ./examples/breeze-ui"];
function Command({ value }: { value: string }) {
  const [copied, setCopied] = useState(false);
  return (
    <div className="install-command">
      <code>{value}</code>
      <button
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
    <section className="install-section" id="install">
      <div className="container install-grid">
        <div>
          <p className="eyebrow">{t.installKicker}</p>
          <h2>{t.installTitle}</h2>
          <p>{t.installText}</p>
          <a
            className="text-link"
            href="https://github.com/yldm-tech/mygo-breeze-ui/blob/main/docs/breeze-ui.md"
            target="_blank"
            rel="noreferrer"
          >
            {t.docs}
            <ArrowRight size={15} />
          </a>
        </div>
        <div className="install-card">
          <div className="install-tabs">
            <span className="active">Go</span>
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
