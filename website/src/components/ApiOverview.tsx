import { ArrowRight, Braces, Database, HeartPulse } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { apiEndpoints } from "../lib/api";
import { ApiPlayground } from "./ApiPlayground";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

const icons = [Database, Braces, HeartPulse];
export function ApiOverview() {
  const { locale } = useUiStore();
  const t = useCopy(locale);
  return (
    <section className="api-section">
      <div className="container">
        <div className="section-heading api-heading">
          <div>
            <p className="eyebrow">{t.apiKicker}</p>
            <h2>{t.apiTitle}</h2>
            <p>{t.apiText}</p>
          </div>
          <Link className="text-link" to="/api">
            {t.apiLink}
            <ArrowRight size={15} />
          </Link>
        </div>
        <div className="api-grid">
          {apiEndpoints.map((endpoint, index) => {
            const Icon = icons[index] ?? Database;
            return (
              <article className="api-card" key={endpoint.path}>
                <div className="api-card-icon">
                  <Icon size={17} />
                </div>
                <div>
                  <div className="api-route">
                    <span>{endpoint.method}</span>
                    {endpoint.path}
                  </div>
                  <h3>{endpoint.label}</h3>
                  <p>{endpoint.description}</p>
                </div>
              </article>
            );
          })}
        </div>
        <ApiPlayground />
      </div>
    </section>
  );
}
