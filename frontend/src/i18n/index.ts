import { en } from './en';
import { es } from './es';

type NestedRecord = { [key: string]: string | NestedRecord };

function flatten(obj: NestedRecord, prefix = ''): Record<string, string> {
  let result: Record<string, string> = {};
  for (const [key, value] of Object.entries(obj)) {
    const path = prefix ? `${prefix}.${key}` : key;
    if (typeof value === 'string') {
      result[path] = value;
    } else {
      result = { ...result, ...flatten(value, path) };
    }
  }
  return result;
}

export type TranslationMap = Record<string, string>;
export type Language = 'en' | 'es';

export const translations: Record<Language, TranslationMap> = {
  en: flatten(en as unknown as NestedRecord),
  es: flatten(es as unknown as NestedRecord),
};
