export function setWindow(name: string, data: any) {
    (window as { [key: string]: any })[name] = data;
}

export function getWindow(name: string): any {
    return (window as { [key: string]: any })[name];
}
