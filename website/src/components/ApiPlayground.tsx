import { useState } from "react";
import { useForm } from "@tanstack/react-form";
import { apiRequestSchema } from "../lib/api";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

export function ApiPlayground() {
  const { locale } = useUiStore();
  const t = useCopy(locale);
  const [result, setResult] = useState("");
  const form = useForm({
    defaultValues: { query: "", kind: "all" as "all" | "component" | "utility" | "token" },
    validators: {
      onChange: ({ value }) => {
        const parsed = apiRequestSchema.safeParse(value);
        return parsed.success ? undefined : parsed.error.issues[0]?.message;
      },
    },
    onSubmit: ({ value }) => setResult(JSON.stringify(apiRequestSchema.parse(value), null, 2)),
  });
  return (
    <div className="api-playground">
      <div className="api-playground-head">
        <span className="api-dot" />
        <span>{t.apiSearch}</span>
      </div>
      <form
        onSubmit={(event) => {
          event.preventDefault();
          event.stopPropagation();
          void form.handleSubmit();
        }}
      >
        <form.Field name="query">
          {(field) => (
            <input
              value={field.state.value}
              onChange={(event) => field.handleChange(event.target.value)}
              placeholder={t.apiSearchPlaceholder}
            />
          )}
        </form.Field>
        <form.Field name="kind">
          {(field) => (
            <select
              value={field.state.value}
              onChange={(event) =>
                field.handleChange(event.target.value as typeof field.state.value)
              }
            >
              <option value="all">{t.all}</option>
              <option value="component">{t.component}</option>
              <option value="utility">{t.utility}</option>
              <option value="token">{t.token}</option>
            </select>
          )}
        </form.Field>
        <button className="button button-primary" type="submit">
          {t.apiRun}
        </button>
      </form>
      {result && (
        <pre className="api-result">
          <span>{t.apiResult}</span>
          {result}
        </pre>
      )}
    </div>
  );
}
