import { Languages } from "lucide-react";

import { Button } from "@/components/ui/button";
import type { Locale } from "@/lib/i18n";

export function LanguageSwitcher({
  locale,
  inverted = false,
}: {
  locale: Locale;
  inverted?: boolean;
}) {
  const target = locale === "en" ? "zh-CN" : "en";

  return (
    <Button
      asChild
      variant="ghost"
      size="sm"
      className={inverted ? "text-card hover:bg-card/10 hover:text-card" : undefined}
    >
      <a
        href={`?lang=${target}`}
        hrefLang={target}
        lang={target}
        aria-label={locale === "en" ? "Switch to Chinese" : "切换为英文"}
      >
        <Languages aria-hidden="true" />
        {target === "en" ? "EN" : "中文"}
      </a>
    </Button>
  );
}
