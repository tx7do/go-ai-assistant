<template>
  <div class="chat-root">
    <n-scrollbar
      ref="historyRef"
      class="chat-scrollbar"
      :content-style="loadingHistory ? { height: '100%' } : {}"
    >
      <div class="scroll-to-bottom-btn">
        <n-button secondary circle size="small" @click="scrollToBottom(true)">
          <template #icon>
            <n-icon :component="ArrowDown" />
          </template>
        </n-button>
      </div>
      <HistoryContent
        ref="historyContentRef"
        :extra-messages="currentActiveMessages"
        :fullscreen="false"
        :show-tips="showFullscreenTips"
        :loading="loadingHistory"
      />
    </n-scrollbar>
    <div class="chat-input-region-wrapper">
      <InputRegion
        v-model:input-value="inputValue"
        v-model:auto-scrolling="autoScrolling"
        :can-abort="canAbort"
        :send-disabled="sendDisabled"
        @abort-request="abortRequest"
        @export-to-markdown-file="exportToMarkdownFile"
        @export-to-pdf-file="exportToPdfFile"
        @send-msg="sendMsg"
        @show-fullscreen-history="showFullscreenHistory"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import {ArrowDown} from '@vicons/ionicons5';
import {NButton, NIcon} from 'naive-ui';
import {computed, ref} from 'vue';

import {RoleEnum} from '../../enums/role';
import {useI18n} from '../../locales';
import {saveAsMarkdown} from '../../utils/export';
import {getWindow} from "../../utils/window";

import {useChatStoreWithOut} from './chat.state';
import HistoryContent from './HistoryContent.vue';
import InputRegion from './InputRegion.vue';
import {SendChatForAssistant} from "../../../wailsjs/go/main/App";

const props = defineProps<{
  loading: boolean;
}>();

const chatStore = useChatStoreWithOut();

const t = useI18n();

const historyRef = ref();
const historyContentRef = ref();

const loadingBar = ref(false);
const loadingHistory = ref<boolean>(false);

const autoScrolling = ref<boolean>(true);
const showFullscreenTips = ref(false);

const isAborted = ref<boolean>(false);
const canAbort = ref<boolean>(false);
let aborter: (() => void) | null = null;

const currentChatHistory = computed<ChatMessage[] | null>(() => {
  return null;
});

const inputValue = ref('');
const enableAsynSendMessage = false;

// 实际的 currentMessageList，加上当前正在发送的消息
const currentActiveMessages = computed<Array<ChatMessage>>(() => {
  console.log('currentActiveMessages ==== ');
  const result: ChatMessage[] = [];
  if (chatStore.currentSendMessage && result.findIndex((message) => message.id === chatStore.currentSendMessage?.id) === -1) {
    result.push(chatStore.currentSendMessage);
    console.log('currentActiveMessages: send');
  }
  if (chatStore.currentRecvMessage && result.findIndex((message) => message.id === chatStore.currentRecvMessage?.id) === -1) {
    result.push(chatStore.currentRecvMessage);
    console.log('currentActiveMessages: recv');
  }
  return result;
});

const sendDisabled = computed(() => {
  return (
      loadingBar.value ||
      inputValue.value === '' ||
      inputValue.value.trim() == ''
  );
});

const abortRequest = () => {
  if (aborter == null || !canAbort.value) return;
  aborter();
  aborter = null;
};

// 滚屏到底部
const scrollToBottom = (smooth: boolean) => {
  if (historyRef.value === undefined) return;

  const scrollHeight = historyRef.value.$refs.scrollbarInstRef.contentRef.scrollHeight;
  // console.log('scrollToBottom', scrollHeight);
  if (smooth) {
    historyRef.value.scrollTo({
      left: 0,
      top: scrollHeight,
      behavior: 'smooth',
    });
  } else {
    historyRef.value.scrollTo({
      left: 0,
      top: scrollHeight
    });
  }
};


// 发送聊天消息
const sendMsg = async () => {
  if (enableAsynSendMessage) {
    await asynSendMsg();
  } else {
    await syncSendMsg();
  }
};

// 异步发送聊天消息
const asynSendMsg = async () => {
  if (sendDisabled.value) {
    getWindow('$message').error(t('tips.pleaseCheckInput'));
    return;
  }
  if (loadingBar.value) {
    getWindow('$message').error(t('tips.sendingChatMessage'));
    return;
  }

  getWindow('$loading').start();
  loadingBar.value = true;
  const message = inputValue.value;
  inputValue.value = '';

  canAbort.value = false;
  isAborted.value = false;
  let hasGotReply = false;

  chatStore.currentSendMessage = {
    id: 0,
    content: message,
    role: RoleEnum.USER,
  };
  chatStore.currentRecvMessage = {
    id: 0,
    content: '',
    role: RoleEnum.AI,
  };

  setTimeout(() => {
    if (autoScrolling.value) {
      scrollToBottom(true);
    }
  }, 45);

  if (hasGotReply) {
    doGotReplay();
  }
};

// 同步发送聊天消息
const syncSendMsg = async () => {
  if (sendDisabled.value) {
    getWindow('$message').error(t('tips.pleaseCheckInput'));
    return;
  }
  if (loadingBar.value) {
    getWindow('$message').error(t('tips.sendingChatMessage'));
    return;
  }

  getWindow('$loading').start();
  loadingBar.value = true;
  const message = inputValue.value;
  inputValue.value = '';

  canAbort.value = false;
  isAborted.value = false;
  let hasGotReply = false;

  chatStore.currentSendMessage = {
    id: 0,
    content: message,
    role: RoleEnum.USER,
  };
  chatStore.currentRecvMessage = {
    id: 0,
    content: '',
    role: RoleEnum.AI,
  };

  setTimeout(() => {
    if (autoScrolling.value) {
      scrollToBottom(true);
    }
  }, 45);

  // console.log('====================================', currentConversation.value, message);

  try {
    chatStore.currentRecvMessage!.content = t('tips.waiting');

    const resp = await SendChatForAssistant(message);
    console.log('GOT Message:', resp);

    hasGotReply = true;

    if (resp && resp.length > 0) {
      chatStore.currentRecvMessage!.content = resp[0];
    }

    chatStore.addMessage(
        chatStore.currentSendMessage!.content || '',
        chatStore.currentRecvMessage!.content || ''
    );

  } catch (e) {
    chatStore.currentRecvMessage!.content = e as string;
    hasGotReply = false;
    console.error(e);
  }

  chatStore.currentSendMessage = null;
  chatStore.currentRecvMessage = null;

  if (hasGotReply) {
    doGotReplay();
  }
};

function doGotReplay() {
  getWindow('$loading').finish();
  loadingBar.value = false;
  isAborted.value = false;

  if (autoScrolling.value) scrollToBottom(false);
}

// 重载聊天记录
const reloadChatHistory = () => {
  loadingBar.value = true;
  getWindow('$loading').start();
};

// 全屏显示聊天记录
const showFullscreenHistory = () => {
  // focus historyContentRef
  historyContentRef.value.focus();
  historyContentRef.value.toggleFullscreenHistory(true);
};

// 导出PDF文件
const exportToPdfFile = () => {
  if (!currentChatHistory.value) {
    return;
  }

  historyContentRef.value.toggleFullscreenHistory(false);
  window.print();
  historyContentRef.value.toggleFullscreenHistory(false);
};

// 导出Markdown文件
const exportToMarkdownFile = () => {
  if (!currentChatHistory.value) {
    return;
  }

  saveAsMarkdown(currentChatHistory.value!);
};
</script>

<style scoped>
.chat-root {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
  background: var(--n-color);
  position: relative;
}

.chat-input-region-wrapper {
  position: fixed;
  left: 0;
  bottom: 0;
  width: 100vw;
  z-index: 30;
  background: var(--n-color);
  box-shadow: 0 -2px 8px rgba(0,0,0,0.04);
  padding: 0 0 0 0;
  display: flex;
  justify-content: center;
}

.chat-input-region-wrapper > * {
  width: 100vw;
  max-width: 900px;
}

.chat-scrollbar {
  flex: 1 1 0%;
  min-height: 0;
  max-height: calc(100vh - 80px);
  padding-bottom: 120px;
  background: transparent;
  width: 100vw;
  box-sizing: border-box;
}

.scroll-to-bottom-btn {
  position: absolute;
  right: 16px;
  bottom: 80px;
  z-index: 20;
}

@media (min-width: 1024px) {
  .chat-root {
    width: 100vw;
    max-width: 900px;
    margin: 0 auto;
  }
  .chat-input-region-wrapper {
    left: 50%;
    transform: translateX(-50%);
    max-width: 900px;
    width: 100vw;
    padding: 0;
  }
  .chat-input-region-wrapper > * {
    width: 100%;
    max-width: 900px;
  }
  .chat-scrollbar {
    max-width: 900px;
    margin: 0 auto;
  }
}

textarea.n-input__textarea-el {
  resize: none;
}

.left-col .n-card__content {
  @apply flex flex-col overflow-auto !important;
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
