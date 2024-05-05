import * as NaiveUI from 'naive-ui';
import {computed} from 'vue';

import {setWindow} from "../utils/window";
import {lighten} from '../utils';
import {APP_THEME} from '../settings/designSetting';

/**
 * 挂载 Naive-ui 脱离上下文的 API
 * 如果你想在 setup 外使用 useDialog、useMessage、useNotification、useLoadingBar，可以通过 createDiscreteApi 来构建对应的 API。
 * https://www.naiveui.com/zh-CN/dark/components/discrete
 */

export function setupNaiveDiscreteApi() {
    const configProviderPropsRef = computed(() => ({
        theme: NaiveUI.darkTheme,
        themeOverrides: {
            common: {
                primaryColor: APP_THEME,
                primaryColorHover: lighten(APP_THEME, 6),
                primaryColorPressed: lighten(APP_THEME, 6),
            },
            LoadingBar: {
                colorLoading: APP_THEME,
            },
        },
    }));
    const {message, dialog, notification, loadingBar} = NaiveUI.createDiscreteApi(
        ['message', 'dialog', 'notification', 'loadingBar'],
        {
            configProviderProps: configProviderPropsRef,
        }
    );

    setWindow('$message', message);
    setWindow('$dialog', dialog);
    setWindow('$notification', notification);
    setWindow('$loading', loadingBar);
}
