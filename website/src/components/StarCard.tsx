import { ArrowUpRight, Star } from "lucide-react";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

export function StarCard() {
  const { locale } = useUiStore();
  const t = useCopy(locale);
  return (
    <section className="pb-[105px] max-[680px]:pb-[72px]">
      <div className="mx-auto grid w-[calc(100%-48px)] max-w-[1160px] grid-cols-[0.9fr_1.1fr] items-center gap-9 rounded-lg border border-line bg-panel-2 p-[37px] max-[940px]:p-9 max-[680px]:grid-cols-1 max-[680px]:gap-7 max-[680px]:w-[calc(100%-34px)] max-[680px]:p-6">
        <div>
          <p className="mb-5 font-mono text-[10px] uppercase tracking-[0.09em] text-mint">GitHub</p>
          <h2 className="mb-3.5 text-[clamp(29px,3.5vw,44px)] leading-[1.1] tracking-[-0.065em]">
            {t.starTitle}
          </h2>
          <p className="mb-5 text-[13px] text-muted">{t.starText}</p>
          <a
            className="inline-flex items-center justify-center gap-2 rounded-[7px] bg-mint px-[17px] py-[11px] text-xs font-extrabold text-bg hover:bg-[#c0f9df]"
            href="https://github.com/yldm-tech/mygo-breeze-ui"
            target="_blank"
            rel="noreferrer"
          >
            {t.source}
            <ArrowUpRight size={16} />
          </a>
        </div>
        <div className="relative ml-4 h-[170px] overflow-hidden max-[680px]:ml-0">
          <div className="absolute inset-0 bg-[linear-gradient(rgba(157,245,208,0.07)_1px,transparent_1px),linear-gradient(90deg,rgba(157,245,208,0.07)_1px,transparent_1px)] bg-[size:34px_34px] opacity-50" />
          <svg
            className="relative h-full w-full text-mint"
            viewBox="0 0 480 170"
            preserveAspectRatio="none"
            aria-hidden="true"
          >
            <path
              d="M0 158 L42 153 L86 145 L128 140 L174 121 L218 117 L260 98 L305 91 L348 65 L392 54 L437 24 L480 5"
              fill="none"
              stroke="currentColor"
              strokeWidth="3"
            />
          </svg>
          <div className="absolute top-2 right-2 flex items-center gap-1.5 font-mono text-[10px] text-mint">
            <Star size={13} fill="currentColor" /> 0 stars
          </div>
        </div>
      </div>
    </section>
  );
}
