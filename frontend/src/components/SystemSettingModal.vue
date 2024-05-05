<template>
  <n-modal
      v-model:show="showModal"
      preset="dialog"
      :title="props.title"
      :mask-closable="true"
      positive-text="创建"
      negative-text="取消"
      :on-mask-click="handleModalClosed"
      :on-positive-click="handleSubmit"
      @negative-click="handleCancel"
  >
    <n-form
        :model="formValue"
        size="small"
        label-placement="left"
    >
      <n-form-item label="助手名称" path="register">
        <n-input v-model:value="formValue.assistantName" placeholder="请输入助手名称"/>
      </n-form-item>
      <n-form-item label="助手描述" path="description">
        <n-input v-model:value="formValue.description" placeholder="请输入助手描述"/>
      </n-form-item>
      <n-form-item label="助手指引" path="instructions">
        <n-input v-model:value="formValue.instructions" type="textarea" placeholder="指引AI助手如何聊天"/>
      </n-form-item>
      <n-form-item label="知识库文档" path="filepath">
        <n-input v-model:value="formValue.uploadFileName" placeholder="请输入文件地址"></n-input>
        <n-button @click="handleOpenFile">
          <template #icon>
            <n-icon>
              <train-icon/>
            </n-icon>
          </template>
        </n-button>
        <input type="file" style="display: none" id="fileInput" ref="fileInputRef" @change="handleChangeFile"/>
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import {computed, onMounted, reactive, ref} from "vue";
import {FileTray as TrainIcon} from '@vicons/ionicons5'
import {NButton, useMessage} from "naive-ui";

import {EventsOn} from "../../wailsjs/runtime";
import {DestroyCurrentAssistant, GetAssistantConfig, InitAssistant} from "../../wailsjs/go/main/App";
import {FileReaderEx} from "../utils/file";
import {getWindow} from "../utils/window";

interface Props {
  show: boolean;
  title?: string;
}

const props = withDefaults(defineProps<Props>(), {
  show: false,
  title: '创建AI助手',
})

const emits = defineEmits<{
  (e: "update:value", val: boolean): void;
  (e: "close-modal", val: boolean): void;
}>();

const fileInputRef = ref();
const message = useMessage();

const formValue = reactive({
  assistantName: '',
  description: '',
  instructions: '',
  uploadFileName: '',
  uploadFileData: new ArrayBuffer(0),
})

const showModal = computed({
  get() {
    return props.show;
  },
  set(val: boolean) {
    if (!val) {
      handleCancel();
    }
    emits('update:value', val);
  }
})

function handleOpenFile() {
  console.log('handleOpenFile')
  fileInputRef.value.click();
}

async function handleChangeFile(event: any) {

  const file = event.target.files[0];
  formValue.uploadFileName = file.name;
  formValue.uploadFileData = await new FileReaderEx().readAsArrayBuffer(file) as ArrayBuffer;

  console.log('handleChangeFile', formValue.uploadFileName, file, formValue.uploadFileData)
}

function closeModal() {
  emits('close-modal', false);
}

async function handleSubmit() {
  if (formValue.assistantName == '') {
    message.error('请输入助理名称');
    return false;
  }
  if (formValue.description == '') {
    message.error('请输入助理描述信息');
    return false;
  }
  if (formValue.instructions == '') {
    message.error('请输入助理导引信息');
    return false;
  }
  if (formValue.uploadFileName == '') {
    message.error('请选择文件');
    return false;
  }

  getWindow('$loading').start();

  await DestroyCurrentAssistant()

  const arrFile = new Uint8Array(formValue.uploadFileData);
  await InitAssistant(formValue.assistantName, formValue.description, formValue.instructions, formValue.uploadFileName, Array.from(arrFile));

  getWindow('$loading').finish();

  message.success('AI助理创建成功！');

  closeModal()

  return true
}

function handleCancel() {
  closeModal()
}

function handleModalClosed() {
  handleCancel();
}

onMounted(async () => {
  const cfg = await GetAssistantConfig();
  formValue.assistantName = cfg.name || '';
  formValue.description = cfg.description || '';
  formValue.instructions = cfg.instructions || '';

  // formValue.uploadFileName = '';
  // formValue.uploadFileData = new ArrayBuffer(0);
})
</script>

<style scoped>
</style>
