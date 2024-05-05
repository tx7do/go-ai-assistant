import {defineStore} from 'pinia';

import {LOCALE, localeSetting} from '../../settings/localeSetting';
import {store} from '..';

const lsLocaleSetting = localeSetting as LocaleSetting;

interface LocaleState {
    localInfo: LocaleSetting;
}

export const useLocaleStore = defineStore({
    id: 'app-locale',
    state: (): LocaleState => ({
        localInfo: lsLocaleSetting,
    }),
    getters: {
        getShowPicker(): boolean {
            return !!this.localInfo?.showPicker;
        },
        getLocale(): LocaleType {
            return this.localInfo?.locale ?? LOCALE.ZH_CN;
        },
    },
    actions: {
        /**
         * Set up multilingual information and cache
         * @param info multilingual info
         */
        setLocaleInfo(info: Partial<LocaleSetting>) {
            this.localInfo = {...this.localInfo, ...info};
        },
        /**
         * Initialize multilingual information and load the existing configuration from the local cache
         */
        initLocale() {
            this.setLocaleInfo({
                ...localeSetting,
                ...this.localInfo,
            });
        },
    },
});

// Need to be used outside the setup
export function useLocaleStoreWithOut() {
    return useLocaleStore(store);
}
