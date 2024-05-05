import {saveAs} from 'file-saver';

import {RoleEnum} from '../enums/role';

export const saveAsMarkdown = (chatHistory: ChatMessage[]) => {
    let content = `# 聊天\n\n`;
    // content += `Date: ${conv.createTime}\nModel: ${conv.model || null}\n`;
    content += '---\n\n';

    let index = 0;
    for (const message of chatHistory) {
        // 选取第一行作为标题，最多50个字符，如果有省略则加上...
        // TODO 适配不同对话
        let title = '';
        if (title.length >= 50) {
            title = title.slice(0, 47) + '...';
        }
        // console.log('------------------', message);

        if (message.role === RoleEnum.USER) {
            content += `## ${++index}. ${title}\n\n`;
            content += `### User\n\n${message.content}\n\n`;
        } else {
            content += `### ChatGPT\n\n${message.content}\n\n`;
            content += '---\n\n';
        }
    }
    const blob = new Blob([content], {type: 'text/plain;charset=utf-8'});
    saveAs(blob, `chatgpt_history.md`);
};
