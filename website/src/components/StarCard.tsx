import { ArrowUpRight, Star } from "lucide-react";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

export function StarCard() {
  const { locale } = useUiStore();
  const t = useCopy(locale);
  return (
    <section className="star-section">
      <div className="container star-card">
        <div>
          <p className="eyebrow">GitHub</p>
          <h2>{t.starTitle}</h2>
          <p>{t.starText}</p>
          <a
            className="button button-primary"
            href="https://github.com/yldm-tech/mygo-breeze-ui"
            target="_blank"
            rel="noreferrer"
          >
            {t.source}
            <ArrowUpRight size={16} />
          </a>
        </div>
        <div className="star-visual">
          <div className="star-grid" />
          <svg
            className="star-chart"
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
          <div className="star-label">
            <Star size={13} fill="currentColor" /> 0 stars
          </div>
        </div>
      </div>
    </section>
  );
}
