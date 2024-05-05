import {createApp} from 'vue'

import App from './App.vue'
import {setupStore} from "./store";
import {setupNaive, setupNaiveDiscreteApi} from "./plugins";
import {setupI18n} from "./locales";

async function bootstrap() {
    const app = createApp(App);

    // 挂载状态管理
    setupStore(app);

    // Multilingual configuration
    // 多语言配置
    // Asynchronous case: language files may be obtained from the server side
    // 异步案例：语言文件可能从服务器端获取
    await setupI18n(app);

    // 注册全局常用的 naive-ui 组件
    setupNaive(app);

    // 挂载 naive-ui 脱离上下文的 Api
    setupNaiveDiscreteApi();

    app.mount('#app');
}

void bootstrap();

declare global {
    interface Window {
        $message: any;
    }
}
