type EventCallback = (data?: any) => any

interface SubscriptionMap {
    [name: string]: EventCallback[]
}

interface Subscription {
    unsubscribe: () => any
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

    private subscribers: SubscriptionMap = {}

    constructor() {
        const ipcRenderer = getIpcRenderer();
        ipcRenderer.on(IpcEventBus.PONG, (data) => this.publish(IpcEventBus.PONG, data));
        ipcRenderer.on(IpcEventBus.KILLED, (data) => this.publish(IpcEventBus.KILLED, data));
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