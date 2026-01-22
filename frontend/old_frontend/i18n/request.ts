import { getRequestConfig } from "next-intl/server";
import { locales, defaultLocale, type Locale } from "./config";

export default getRequestConfig(async ({ locale }) => {
  const current: Locale = (locales as readonly string[]).includes(locale as string) ? (locale as Locale) : defaultLocale;

  return {
    locale: current,
    messages: (await import(`./${current}.json`)).default,
  };
});
