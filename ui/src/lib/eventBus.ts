type EventCallback = (data?: any) => any

interface SubscriptionMap {
    [name: string]: EventCallback[]
}

interface Subscription {
    unsubscribe: () => any
}

export interface Subscriber {
    subscribe: (ev: string, callback: EventCallback) => Subscription,
}

export const getIpcRenderer = () => {
    if (!window || !window.process || !window.require) {
        throw new Error("IPC Renderer is not available")
    }

    return window.require('electron').ipcRenderer;
}

export class IpcEventBus {
    public static SERVER_UP = "ipc_server_up";
    public static PONG = "kad:pong";
    public static KILLED = "kad:kill";
    public static STORAGE_SUCCESS = "storage:success";
    public static STORAGE_ERROR = "storage:error";
    public static STORAGE_FILES = "storage:files";

    private subscribers: SubscriptionMap = {}

    constructor() {
        const ipcRenderer = getIpcRenderer();
        ipcRenderer.on(IpcEventBus.PONG, (_, data) => this.publish(IpcEventBus.PONG, data));
        ipcRenderer.on(IpcEventBus.KILLED, (_, data) => this.publish(IpcEventBus.KILLED, data));
        ipcRenderer.on(IpcEventBus.STORAGE_SUCCESS, (_, data) => this.publish(IpcEventBus.STORAGE_SUCCESS, data));
        ipcRenderer.on(IpcEventBus.STORAGE_FILES, (_, data) => this.publish(IpcEventBus.STORAGE_FILES, data))
    }

    public subscribe(ev: string, callback: EventCallback): Subscription {
        if (this.subscribers[ev]) {
            this.subscribers[ev].push(callback)
        } else {
            this.subscribers[ev] = [callback]
        }

        return {
            unsubscribe: () => {
                const index = this.subscribers[ev].indexOf(callback);
                if (index > -1) {
                    this.subscribers[ev].splice(index, 1)
                }
            }
        }
    }

    private publish(ev: string, data: any) {
        if (this.subscribers[ev]) {
            this.subscribers[ev].forEach((cb) => cb(data))
        }
    }
}

let eventBusInstance: IpcEventBus = null;

// returns a singleton IpcEventBus
export const getIPCEventBus = () => {
    if (!eventBusInstance) {
        eventBusInstance = new IpcEventBus();
    }

    return eventBusInstance;
}