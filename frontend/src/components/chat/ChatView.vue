<template>
  <div class="h-full relative flex flex-col">
    <!-- 消息记录内容（用于全屏展示） -->
    <n-scrollbar
        ref="historyRef"
        class="relative"
        :content-style="loadingHistory ? { height: '100%' } : {}"
    >
      <!-- 回到底部按钮 -->
      <div class="right-2 bottom-5 absolute z-20">
        <n-button secondary circle size="small" @click="scrollToBottom(true)">
          <template #icon>
            <n-icon :component="ArrowDown"/>
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
    <!-- 下半部分（回复区域） -->
    <InputRegion
        v-model:input-value="inputValue"
        v-model:auto-scrolling="autoScrolling"
        class="sticky bottom-0 z-10"
        :can-abort="canAbort"
        :send-disabled="sendDisabled"
        @abort-request="abortRequest"
        @export-to-markdown-file="exportToMarkdownFile"
        @export-to-pdf-file="exportToPdfFile"
        @send-msg="sendMsg"
        @show-fullscreen-history="showFullscreenHistory"
    />
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
textarea.n-input__textarea-el {
  resize: none;
}

div.n-menu-item-content-header {
  display: flex;
  justify-content: space-between;
}

span.n-menu-item-content-header__extra {
  display: inline-block;
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

  /* no margin in page */
  @page {
    margin-left: 0;
    margin-right: 0;
  }
}
</style>
