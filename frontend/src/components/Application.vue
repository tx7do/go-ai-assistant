<template>
  <n-layout
    ref="rootRef"
    class="h-full w-full lg:w-screen-lg lg:mx-auto"
    style="min-height: 100vh;"
  >
    <n-layout-header bordered>
      <n-space justify="end">
        <n-button type="primary" @click="handleClick">创建AI助理</n-button>
      </n-space>
    </n-layout-header>

    <n-layout-content
      embedded
      class="flex flex-col flex-1 overflow-hidden w-full min-h-0"
      style="min-height: calc(100vh - 56px);"
    >
      <ChatView :loading="loadingBar" />
    </n-layout-content>
  </n-layout>
  <SystemSettingModal :show="showSystemSetting" @close-modal="closeSystemSetting" />
</template>

<script lang="ts" setup>
import {ref} from "vue";
import {NButton} from "naive-ui";

import {screenWidthGreaterThan} from '../utils/media';
import ChatView from "./chat/ChatView.vue";
import SystemSettingModal from "./SystemSettingModal.vue";

const loadingBar = ref(false);
const rootRef = ref();
const gtmd = screenWidthGreaterThan('md');
const showSystemSetting = ref<boolean>(false)

function openSystemSetting() {
  showSystemSetting.value = true;
}

function closeSystemSetting() {
  showSystemSetting.value = false;
}

function handleClick() {
  openSystemSetting();
}

</script>

<style scoped>
:deep(.n-layout) {
  min-height: 100vh;
  width: 100vw;
}

.n-layout-content {
  flex: 1 1 auto;
  min-height: 0;
  width: 100%;
}

html, body, #app {
  height: 100%;
  margin: 0;
  padding: 0;
}

.content {
  display: block;
  width: 98%;
  height: 100%;
  margin: auto;
  padding: 2% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}

@media print {
  body * {
    visibility: hidden;
  }
  #print-content * {
    visibility: visible;
  }
  @page {
    margin-left: 0;
    margin-right: 0;
  }
}
</style>
