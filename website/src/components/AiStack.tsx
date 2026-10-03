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
    <section className="ai-stack" id="ai-stack">
      <div className="container">
        <div className="ai-stack-heading">
          <div>
            <p className="eyebrow">
              <Sparkles size={12} />
              {t.aiKicker}
            </p>
            <h2>{t.aiTitle}</h2>
            <p>{t.aiText}</p>
          </div>
          <a
            className="text-link"
            href="https://github.com/yldm-tech/ai-logo"
            target="_blank"
            rel="noreferrer"
          >
            {t.aiSource}
            <ArrowUpRight size={14} />
          </a>
        </div>
        <div className="ai-stack-grid">
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
        <p className="ai-stack-note">{t.aiDisclaimer}</p>
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
      className="ai-card"
      initial={{ opacity: 0, y: 12 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, amount: 0.25 }}
      transition={{ delay, duration: 0.45 }}
    >
      <div className="ai-card-icon">{icon}</div>
      <div>
        <h3>{name}</h3>
        <p>{text}</p>
      </div>
    </motion.div>
  );
}
