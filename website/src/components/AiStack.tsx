import { Claude, Codex, Cursor, DeepSeek, Gemini, OpenAI } from "@yldm-tech/ai-logo";
import type { ReactNode } from "react";
import { motion } from "motion/react";
import { ArrowUpRight, Sparkles } from "lucide-react";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

export function AiStack() {
  const { locale } = useUiStore();
  const t = useCopy(locale);

  return (
    <section
      className="border-t border-line bg-[linear-gradient(180deg,rgba(16,35,55,0.28),rgba(7,17,31,0))] px-0 py-[118px] pb-[124px] max-[680px]:py-[76px] max-[680px]:pb-[82px]"
      id="ai-stack"
    >
      <div className="mx-auto w-[calc(100%-48px)] max-w-[1160px] max-[680px]:w-[calc(100%-34px)]">
        <div className="mb-[38px] flex items-end justify-between gap-[45px] max-[680px]:mb-[26px] max-[680px]:block">
          <div className="max-w-[610px]">
            <p className="mb-5 flex items-center gap-[9px] font-mono text-[10px] uppercase tracking-[0.09em] text-mint">
              <Sparkles size={12} />
              {t.aiKicker}
            </p>
            <h2 className="mb-[13px] text-[clamp(29px,3.5vw,44px)] leading-[1.1] tracking-[-0.065em]">
              {t.aiTitle}
            </h2>
            <p className="m-0 max-w-[560px] text-muted">{t.aiText}</p>
          </div>
          <a
            className="inline-flex items-center gap-1.5 whitespace-nowrap font-mono text-[10px] text-mint hover:text-ink max-[680px]:mt-5"
            href="https://github.com/yldm-tech/ai-logo"
            target="_blank"
            rel="noreferrer"
          >
            {t.aiSource}
            <ArrowUpRight size={14} />
          </a>
        </div>
        <div className="grid grid-cols-3 gap-2.5 max-[940px]:grid-cols-2 max-[680px]:grid-cols-1">
          <AiCard icon={<OpenAI size={30} />} name="OpenAI" text={t.aiOpenAi} delay={0} />
          <AiCard icon={<Claude.Color size={30} />} name="Claude" text={t.aiClaude} delay={0.04} />
          <AiCard icon={<Gemini.Color size={30} />} name="Gemini" text={t.aiGemini} delay={0.08} />
          <AiCard icon={<Cursor size={30} />} name="Cursor" text={t.aiCursor} delay={0.12} />
          <AiCard icon={<Codex.Color size={30} />} name="Codex" text={t.aiCodex} delay={0.16} />
          <AiCard
            icon={<DeepSeek.Color size={30} />}
            name="DeepSeek"
            text={t.aiDeepSeek}
            delay={0.2}
          />
        </div>
        <p className="m-0 mt-[21px] font-mono text-[9px] text-dim">{t.aiDisclaimer}</p>
      </div>
    </section>
  );
}

function AiCard({
  icon,
  name,
  text,
  delay,
}: {
  icon: ReactNode;
  name: string;
  text: string;
  delay: number;
}) {
  return (
    <motion.div
      className="flex min-h-[132px] items-start gap-4 rounded-[11px] border border-line bg-[rgba(12,27,43,0.74)] p-[21px] transition hover:-translate-y-0.5 hover:border-[rgba(157,245,208,0.55)] hover:bg-[rgba(16,35,55,0.9)]"
      initial={{ opacity: 0, y: 12 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, amount: 0.25 }}
      transition={{ delay, duration: 0.45 }}
    >
      <div className="grid size-[46px] shrink-0 place-items-center rounded-[10px] border border-[rgba(157,245,208,0.12)] bg-[rgba(7,17,31,0.75)]">
        {icon}
      </div>
      <div>
        <h3 className="mb-1.5 mt-px text-sm">{name}</h3>
        <p className="m-0 text-[11px] leading-[1.55] text-muted">{text}</p>
      </div>
    </motion.div>
  );
}
