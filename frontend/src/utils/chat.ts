export const chatModelColorMap: Record<string, string> = {
    gpt_3_5: 'green',
    gpt_4: 'purple',
    gpt_4_mobile: 'darkblue',
    gpt_4_browsing: 'purple',
    gpt_4_plugins: 'purple',
};

export const getChatModelColor = (model_name: string | null) => {
    if (model_name && chatModelColorMap[model_name]) return chatModelColorMap[model_name];
    else return 'black';
};

export const getChatModelIconStyle = (model_name: string | null) => {
    if (model_name == 'gpt_4_plugins') return 'plugins';
    else if (model_name == 'gpt_4_browsing') return 'browsing';
    else return 'default';
};
