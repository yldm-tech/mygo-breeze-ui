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
    <div className="mt-3 rounded-lg border border-line bg-panel-2 p-5">
      <div className="mb-3.5 flex items-center gap-2 font-mono text-[10px] text-muted">
        <span className="size-1.5 rounded-full bg-mint-strong" />
        <span>{t.apiSearch}</span>
      </div>
      <form
        className="flex gap-2.5 max-[680px]:flex-wrap"
        onSubmit={(event) => {
          event.preventDefault();
          event.stopPropagation();
          void form.handleSubmit();
        }}
      >
        <form.Field name="query">
          {(field) => (
            <input
              className="min-w-0 flex-1 rounded-md border border-line bg-bg px-3 py-2 text-xs outline-none placeholder:text-dim focus:border-mint max-[680px]:min-w-full"
              value={field.state.value}
              onChange={(event) => field.handleChange(event.target.value)}
              placeholder={t.apiSearchPlaceholder}
            />
          )}
        </form.Field>
        <form.Field name="kind">
          {(field) => (
            <select
              className="rounded-md border border-line bg-bg px-3 py-2 text-xs outline-none focus:border-mint max-[680px]:flex-1"
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
        <button
          className="inline-flex items-center justify-center gap-2 rounded-[7px] bg-mint px-[17px] py-[11px] text-xs font-extrabold text-bg hover:bg-[#c0f9df]"
          type="submit"
        >
          {t.apiRun}
        </button>
      </form>
      {result && (
        <pre className="mt-3.5 whitespace-pre-wrap rounded-md border border-line bg-bg p-3 font-mono text-[10px] text-muted">
          <span className="mb-2 block text-mint">{t.apiResult}</span>
          {result}
        </pre>
      )}
    </div>
  );
}
