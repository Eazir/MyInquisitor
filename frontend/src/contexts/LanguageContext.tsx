import type { ReactNode } from 'react';
import { createContext, useContext, useState, useCallback } from 'react';
import { translations, type Language } from '../i18n';
export type { Language };

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (key: string, params?: Record<string, string | number>) => string;
}

const LanguageContext = createContext<LanguageContextType>({} as LanguageContextType);

function getInitialLanguage(): Language {
  const stored = localStorage.getItem('language');
  if (stored === 'en' || stored === 'es') return stored;
  return 'es';
}

export function LanguageProvider({ children }: { children: ReactNode }) {
  const [language, setLanguageState] = useState<Language>(getInitialLanguage);

  const setLanguage = useCallback((lang: Language) => {
    setLanguageState(lang);
    localStorage.setItem('language', lang);
  }, []);

  const t = useCallback((key: string, params?: Record<string, string | number>): string => {
    let value = translations[language][key];
    if (value === undefined) {
      value = translations.en[key];
    }
    if (value === undefined) {
      return key;
    }
    if (params) {
      for (const [k, v] of Object.entries(params)) {
        value = value.replace(`{{${k}}}`, String(v));
      }
    }
    return value;
  }, [language]);

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  );
}

export const useLanguage = () => useContext(LanguageContext);
