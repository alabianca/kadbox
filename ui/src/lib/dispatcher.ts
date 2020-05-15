export interface Dispatcher {
    dispatch(name: string, arg?: any): any
}

export class IpcRenderer implements Dispatcher{
    constructor(private ipcRenderer: any) {
    }

    public dispatch(name: string, arg?: any): any {
        this.ipcRenderer.send(name, arg)
    }
}