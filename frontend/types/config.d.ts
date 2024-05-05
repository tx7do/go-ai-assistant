type LocaleType = 'zh-CN' | 'en-US' | 'ru' | 'ja' | 'ko';

interface LocaleSetting {
    showPicker: boolean;
    // Current language
    locale: LocaleType;
    // default language
    fallback: LocaleType;
    // available Locales
    availableLocales: LocaleType[];
}
