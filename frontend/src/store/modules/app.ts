import {defineStore} from 'pinia';
import {useStorage, RemovableRef} from '@vueuse/core';

import {store} from '..';

declare type Preference = {
    sendKey: 'Shift+Enter' | 'Enter' | 'Ctrl+Enter';
    renderUserMessageInMd: boolean;
    codeAutoWrap: boolean;
    widerConversationPage: boolean;
};

interface AppState {
    preference: RemovableRef<Preference>;
}

const useAppStore = defineStore('app', {
    state: (): AppState => ({
        preference: useStorage<Preference>('preference', {
            sendKey: 'Enter',
            renderUserMessageInMd: false,
            codeAutoWrap: false,
            widerConversationPage: true,
        }),
    }),
    getters: {},
    actions: {},
});

// Need to be used outside the setup
export function useAppStoreWithOut() {
    return useAppStore(store);
}
