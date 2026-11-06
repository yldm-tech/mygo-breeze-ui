import { useEffect } from "react";

const siteUrl = "https://yldm-tech.github.io/mygo-breeze-ui";
const imageUrl = `${siteUrl}/og-image.png`;

type SeoProps = {
  title: string;
  description: string;
  socialTitle?: string;
  socialDescription?: string;
  path: string;
  locale: "zh" | "en";
};

function setMeta(selector: string, attribute: "name" | "property", key: string, content: string) {
  let element = document.head.querySelector<HTMLMetaElement>(selector);
  if (!element) {
    element = document.createElement("meta");
    element.setAttribute(attribute, key);
    document.head.appendChild(element);
  }
  element.content = content;
}

export function Seo({ title, description, socialTitle = title, socialDescription = description, path, locale }: SeoProps) {
  useEffect(() => {
    const url = `${siteUrl}${path === "/" ? "/" : path}`;
    document.title = title;
    document.documentElement.lang = locale === "zh" ? "zh-CN" : "en";
    setMeta('meta[name="description"]', "name", "description", description);
    setMeta('meta[property="og:title"]', "property", "og:title", socialTitle);
    setMeta('meta[property="og:description"]', "property", "og:description", socialDescription);
    setMeta('meta[property="og:url"]', "property", "og:url", url);
    setMeta('meta[name="twitter:title"]', "name", "twitter:title", socialTitle);
    setMeta('meta[name="twitter:description"]', "name", "twitter:description", socialDescription);
    let canonical = document.head.querySelector<HTMLLinkElement>('link[rel="canonical"]');
    if (!canonical) {
      canonical = document.createElement("link");
      canonical.rel = "canonical";
      document.head.appendChild(canonical);
    }
    canonical.href = url;
    let structured = document.head.querySelector<HTMLScriptElement>('script[data-seo="page"]');
    if (!structured) {
      structured = document.createElement("script");
      structured.type = "application/ld+json";
      structured.dataset.seo = "page";
      document.head.appendChild(structured);
    }
    structured.textContent = JSON.stringify({
      "@context": "https://schema.org",
      "@type": "WebPage",
      name: title,
      description,
      url,
      isPartOf: { "@type": "WebSite", name: "MyGo Breeze", url: `${siteUrl}/` },
      image: imageUrl,
    });
  }, [description, locale, path, socialDescription, socialTitle, title]);
  return null;
}
