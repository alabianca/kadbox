import {IpcEventBus, getIpcRenderer } from "./eventBus";
import {IpcRenderer} from "./dispatcher";
import {FileManager} from "./fileManager";

export const EVENT_PING = "kad:ping";
export const EVENT_ADD_FILE = "storage:add";
export const EVENT_GET_FILE = "storage:get-files";

class ConnectionManager {
    private evBus: IpcEventBus = null;
    private ipcRenderer: IpcRenderer = null;
    private fileManager: FileManager = null;

    constructor() {
        this.ipcRenderer = new IpcRenderer(getIpcRenderer());
    }

    public ping() {
        console.log("Sending Ping")
        this.ipcRenderer.dispatch(EVENT_PING)
    }

    public EventBus(): IpcEventBus {
        if (!this.evBus) {
            this.evBus = new IpcEventBus();
        }

        return this.evBus;
    }

    public FileManager(): FileManager {
        if (!this.fileManager) {
            this.fileManager = new FileManager(this.ipcRenderer, this.EventBus())
        }

        return this.fileManager;
    }
}

let connMgr: ConnectionManager = null;
export const getConnectionManager = () => {
    if (!connMgr) {
        connMgr = new ConnectionManager();
    }

    return connMgr;
}