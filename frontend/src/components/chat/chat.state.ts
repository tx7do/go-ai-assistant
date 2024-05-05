import {defineStore} from 'pinia';

import {store} from '../../store';
import {RoleEnum} from "../../enums/role";

export interface ChatState {
    currentSendMessage: Nullable<ChatMessage>;
    currentRecvMessage: Nullable<ChatMessage>;
    chatMessages: ChatMessage[];
}

export const useChatStore = defineStore({
    id: 'chat',
    state: (): ChatState => ({
        currentSendMessage: null,
        currentRecvMessage: null,
        chatMessages: [],
    }),
    getters: {
        getCurrentSendMessage(): ChatMessage | null {
            return this.currentSendMessage;
        },
        getCurrentRecvMessage(): ChatMessage | null {
            return this.currentRecvMessage;
        },
    },
    actions: {
        getChatMessages(): ChatMessage[] | null {
            return this.chatMessages;
        },

        addMessage(prompt: string, answer: string) {
            this.chatMessages.push({
                id: this.chatMessages.length,
                role: RoleEnum.USER,
                content: prompt
            });
            this.chatMessages.push({
                id: this.chatMessages.length,
                role: RoleEnum.AI,
                content: answer
            });
        },
    },
});

// Need to be used outside the setup
export function useChatStoreWithOut() {
    return useChatStore(store);
}
