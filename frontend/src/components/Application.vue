<template>
  <n-layout
      ref="rootRef"
      :class="['h-full', 'lg:w-screen-lg lg:mx-auto' ]"
  >

    <n-layout-header bordered>
      <n-space justify="end">
        <n-button type="primary" @click="handleClick">创建AI助理</n-button>
      </n-space>
    </n-layout-header>

    <n-layout-content
        embeded
        :class="['flex flex-col overflow-hidden', gtmd() ? '' : 'min-w-100vw']"
    >
      <ChatView :loading="loadingBar"/>
    </n-layout-content>
  </n-layout>
  <SystemSettingModal :show="showSystemSetting" @close-modal="closeSystemSetting"/>
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

  /* no margin in page */
  @page {
    margin-left: 0;
    margin-right: 0;
  }
}
</style>
