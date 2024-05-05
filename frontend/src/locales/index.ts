import {App, computed} from 'vue';
import {createI18n, type I18n} from 'vue-i18n';

import {useLocaleStoreWithOut} from '../store/modules';
import {LOCALE} from '../settings/localeSetting';

import EN from './lang/en-US.json';
import ZH from './lang/zh-CN.json';


let i18n: I18n;

export function setHtmlPageLang(locale: LocaleType) {
    document.querySelector('html')?.setAttribute('lang', locale);
}

function setI18nLanguage(locale: LocaleType) {
    const localeStore = useLocaleStoreWithOut();

    if (i18n.mode === 'legacy') {
        i18n.global.locale = locale;
    } else {
        (i18n.global.locale as any).value = locale;
    }
    localeStore.setLocaleInfo({locale: locale});
    setHtmlPageLang(locale);

    return locale;
}

export async function changeLocale(locale: LocaleType) {
    const localeStore = useLocaleStoreWithOut();

    const globalI18n = i18n.global;

    localeStore.setLocaleInfo({locale: locale});

    let message: any;
    if (locale === LOCALE.ZH_CN) {
        message = ZH;
    } else if (locale === LOCALE.EN_US) {
        message = EN;
    }

    globalI18n.setLocaleMessage(locale, message);

    setI18nLanguage(locale);

    return locale;
}

const init = () => {
    const localeStore = useLocaleStoreWithOut();
    const getLocale = computed(() => localeStore.getLocale);

    i18n = createI18n({
        legacy: false,
        locale: getLocale.value,
        messages: {
            'en-US': {
                ...EN,
            },
            'zh-CN': {
                ...ZH,
            },
        },
    });
};

export function useI18n() {
    return i18n.global.t as any;
}

init();

// setup i18n instance with glob
export async function setupI18n(app: App) {
    app.use(i18n);
}
