<template>
  <n-config-provider
      :locale="zhCN"
      :theme="getDarkTheme"
      :date-locale="dateZhCN"
  >
    <n-dialog-provider>
      <n-notification-provider>
        <n-message-provider>
          <n-loading-bar-provider>
            <n-scrollbar trigger="none">
              <Application/>
            </n-scrollbar>
          </n-loading-bar-provider>
        </n-message-provider>
      </n-notification-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<script lang="ts" setup>
import {computed} from "vue";
import {
  zhCN,
  dateZhCN,
  darkTheme,
  NMessageProvider,
  NDialogProvider,
  NNotificationProvider,
  NLoadingBarProvider
} from 'naive-ui';
import Application from './components/Application.vue'

const getDarkTheme = computed(() => (darkTheme));
</script>
