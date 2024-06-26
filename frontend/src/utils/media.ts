import {useWindowSize} from '@vueuse/core';

const {width} = useWindowSize();

export const breakpoints = {
    sm: 640,
    md: 768,
    lg: 1024,
    xl: 1280,
} as const;

export function screenWidthGreaterThan(breakpoint: keyof typeof breakpoints | number) {
    return () => {
        return width.value >= (typeof breakpoint === 'number' ? breakpoint : breakpoints[breakpoint]);
    };
}
